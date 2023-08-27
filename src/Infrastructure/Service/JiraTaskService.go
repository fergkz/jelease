package InfrastructureService

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainTool "github.com/fergkz/jelease/src/Domain/Tool"
)

type JiraTaskService struct {
	jiraQueryService     *jiraQueryService
	hostname             string
	CustomFieldObjective string
	CustomFieldTaskType  string
}

func NewJiraTaskService(
	Username string,
	AccessToken string,
	Hostname string,
	CustomFieldObjective string,
	CustomFieldTaskType string,
) *JiraTaskService {
	service := new(JiraTaskService)
	service.jiraQueryService = NewJiraQueryService(Username, AccessToken, Hostname)
	service.hostname = Hostname
	service.CustomFieldObjective = CustomFieldObjective
	service.CustomFieldTaskType = CustomFieldTaskType
	return service
}

func (service JiraTaskService) GetTasksFromSprints(SprintIds []DomainEntity.ProjectSprintId) (Tasks []DomainEntity.ProjectTask, Sprints []DomainEntity.ProjectSprint) {

	Sprints = service.jiraQueryService.GetSprints(SprintIds)

	sprintsIdsStrs := []string{}
	for _, i := range SprintIds {
		sprintsIdsStrs = append(sprintsIdsStrs, strconv.Itoa(int(i)))
	}
	sprintsIdsStr := strings.Join(sprintsIdsStrs, ",")

	jql := `
		sprint IN (` + sprintsIdsStr + `) order by rank ASC
	`

	response := new([]interface{})
	cacheFilename := "cache/sprint-" + sprintsIdsStr
	if !DomainTool.Pretty.GetCache(cacheFilename, response) {
		*response = service.jiraQueryService.Query(jql)
		DomainTool.Pretty.SetCache(cacheFilename, response, 60000)
	}

	Tasks = service.parseToTasks(*response, Sprints)

	return Tasks, Sprints
}

func (service JiraTaskService) parseToTasks(rows []interface{}, sprints []DomainEntity.ProjectSprint) (Tasks []DomainEntity.ProjectTask) {

	type taskDTOStruct struct {
		Id           int `json:",string"`
		Key          string
		FieldsStruct struct {
			Summary   string
			Issuetype struct {
				Name           string
				HierarchyLevel int
				Subtask        bool
			}
			Status struct {
				Name string
			}
			TimeOriginalEstimate int
			TimeSpent            int
			Assignee             struct {
				AccountId    string
				EmailAddress string
				DisplayName  string
				AvatarUrls   struct {
					Image string `json:"48x48"`
				}
			}
			Reporter struct {
				AccountId    string
				EmailAddress string
				DisplayName  string
				AvatarUrls   struct {
					Image string `json:"48x48"`
				}
			}
			Parent struct {
				Key    string
				Fields struct {
					Summary string
					Status  struct {
						Name string
					}
				}
			}
			Subtasks []struct {
				Key    string
				Fields struct {
					Status struct {
						Name string
					}
				}
			}
			Updated           interface{}
			Created           interface{}
			Aggregateprogress struct {
				Percent float64
			}
			Comment struct {
				Comments []struct {
					Author struct {
						AccountId    string
						EmailAddress string
						DisplayName  string
						AvatarUrls   struct {
							Image string `json:"48x48"`
						}
					}
					Body      string
					Created   interface{} //"2022-09-01T15:05:11.156-0300"
					Updated   interface{} //"2022-09-01T15:05:11.156-0300"
					JsdPublic bool
				}
			}
		}
		FieldsMap map[string]interface{} `json:"fields"`
		Changelog struct {
			Histories []struct {
				Author struct {
					AccountId    string
					EmailAddress string
					DisplayName  string
					AvatarUrls   struct {
						Image string `json:"48x48"`
					}
				}
				Created interface{} //"2022-09-01T15:05:11.156-0300"
				Items   []struct {
					From       string
					To         string
					FromString string
					ToString   string
					Field      string
					FieldId    string
				}
			}
		}
	}

	allDTOs := map[string]taskDTOStruct{}
	allOrderedKeys := []string{}

	for _, row := range rows {
		dto := new(taskDTOStruct)
		byteRow, _ := json.Marshal(row)
		json.Unmarshal(byteRow, &dto)

		dbByte, _ := json.Marshal(dto.FieldsMap)
		_ = json.Unmarshal(dbByte, &dto.FieldsStruct)

		allDTOs[dto.Key] = *dto
		DomainTool.Pretty.SetCache("cache/issues/"+dto.Key, row, 0)
		DomainTool.Pretty.SetCache("cache/tasks/"+dto.Key, dto, 0)

		if !dto.FieldsStruct.Issuetype.Subtask {
			allOrderedKeys = append(allOrderedKeys, dto.Key)
		}
	}
	issetUserSet := map[string]bool{}

	for _, key := range allOrderedKeys {
		dto := allDTOs[key]

		Task := DomainEntity.ProjectTask{}

		Task.Key = dto.Key
		Task.Summary = dto.FieldsStruct.Summary
		Task.Type = dto.FieldsStruct.Issuetype.Name
		Task.Status = dto.FieldsStruct.Status.Name
		if service.CustomFieldObjective != "" && dto.FieldsMap[service.CustomFieldObjective] != nil {
			Task.Objective = dto.FieldsMap[service.CustomFieldObjective].(string)
		}
		if service.CustomFieldTaskType != "" && dto.FieldsMap[service.CustomFieldTaskType] != nil {
			type CustomFieldTaskTypeStruct struct {
				Value string
			}
			customFieldTaskType := new(CustomFieldTaskTypeStruct)
			byteCustomFieldTaskType, _ := json.Marshal(dto.FieldsMap[service.CustomFieldTaskType])
			json.Unmarshal(byteCustomFieldTaskType, &customFieldTaskType)
			Task.TaskType = customFieldTaskType.Value
		}

		Task.CreatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", dto.FieldsStruct.Created.(string))
		Task.UpdatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", dto.FieldsStruct.Updated.(string))

		if dto.FieldsStruct.TimeOriginalEstimate > 0 {
			Task.TimeEstimateHours = int(math.Ceil(float64(dto.FieldsStruct.TimeOriginalEstimate) / 60 / 60))
		}
		if dto.FieldsStruct.TimeSpent > 0 {
			Task.TimeSpentHours = int(math.Ceil(float64(dto.FieldsStruct.TimeSpent) / 60 / 60))
		}

		for _, comment := range dto.FieldsStruct.Comment.Comments {
			oComment := DomainEntity.ProjectComment{}
			oComment.Body = comment.Body
			oComment.CreatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", comment.Created.(string))
			oComment.Public = comment.JsdPublic
			oComment.User = DomainEntity.ProjectUser{
				Id:        comment.Author.AccountId,
				Email:     comment.Author.EmailAddress,
				Name:      comment.Author.DisplayName,
				AvatarUrl: comment.Author.AvatarUrls.Image,
			}
			Task.Comments = append(Task.Comments, oComment)
		}

		Task.Assignees = append(Task.Assignees, DomainEntity.ProjectUserExecutor{
			User: DomainEntity.ProjectUser{
				Id:        dto.FieldsStruct.Assignee.AccountId,
				Email:     dto.FieldsStruct.Assignee.EmailAddress,
				Name:      dto.FieldsStruct.Assignee.DisplayName,
				AvatarUrl: dto.FieldsStruct.Assignee.AvatarUrls.Image,
			},
			ContributionPerc: 0,
		})
		issetUserSet[Task.Key+dto.FieldsStruct.Assignee.EmailAddress] = true

		for _, subtask := range dto.FieldsStruct.Subtasks {
			dtoSub := allDTOs[subtask.Key]
			if !issetUserSet[Task.Key+dtoSub.FieldsStruct.Assignee.EmailAddress] {
				Task.Assignees = append(Task.Assignees, DomainEntity.ProjectUserExecutor{
					User: DomainEntity.ProjectUser{
						Id:        dtoSub.FieldsStruct.Assignee.AccountId,
						Email:     dtoSub.FieldsStruct.Assignee.EmailAddress,
						Name:      dtoSub.FieldsStruct.Assignee.DisplayName,
						AvatarUrl: dtoSub.FieldsStruct.Assignee.AvatarUrls.Image,
					},
					ContributionPerc: 0,
				})
			}
			issetUserSet[Task.Key+dtoSub.FieldsStruct.Assignee.EmailAddress] = true
		}

		for _, history := range dto.Changelog.Histories {
			if issetUserSet[Task.Key+history.Author.EmailAddress] {
				Task.Assignees = append(Task.Assignees, DomainEntity.ProjectUserExecutor{
					User: DomainEntity.ProjectUser{
						Id:        history.Author.AccountId,
						Email:     history.Author.EmailAddress,
						Name:      history.Author.DisplayName,
						AvatarUrl: history.Author.AvatarUrls.Image,
					},
					ContributionPerc: 0,
				})
			}
			issetUserSet[Task.Key+history.Author.EmailAddress] = true
		}

		Task.Reporter = DomainEntity.ProjectUser{
			Id:        dto.FieldsStruct.Reporter.AccountId,
			Email:     dto.FieldsStruct.Reporter.EmailAddress,
			Name:      dto.FieldsStruct.Reporter.DisplayName,
			AvatarUrl: dto.FieldsStruct.Reporter.AvatarUrls.Image,
		}

		Task.Epic = DomainEntity.ProjectEpic{
			Key:     dto.FieldsStruct.Parent.Key,
			Summary: dto.FieldsStruct.Parent.Fields.Summary,
			Status:  dto.FieldsStruct.Parent.Fields.Status.Name,
		}

		for _, sprint := range sprints {
			if Task.CreatedAt.Before(sprint.StartDate) {
				Task.Planned = true
			}
		}

		Task.PublicHtmlUrl = service.hostname + "/browse/" + Task.Key

		Tasks = append(Tasks, Task)
	}

	return Tasks
}
