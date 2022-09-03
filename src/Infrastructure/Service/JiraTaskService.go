package InfrastructureService

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
)

type JiraTaskService struct {
	jiraQueryService *jiraQueryService
	hostname         string
}

func NewJiraTaskService(
	Username string,
	AccessToken string,
	Hostname string,
) *JiraTaskService {
	service := new(JiraTaskService)
	service.jiraQueryService = NewJiraQueryService(Username, AccessToken, Hostname)
	service.hostname = Hostname
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

	rows := service.jiraQueryService.Query(jql)

	Tasks = service.parseToTasks(rows, Sprints)

	return Tasks, Sprints
}

func (service JiraTaskService) parseToTasks(rows []interface{}, sprints []DomainEntity.ProjectSprint) (Tasks []DomainEntity.ProjectTask) {

	type taskDTOStruct struct {
		Id     int `json:",string"`
		Key    string
		Fields struct {
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
	}

	allDTOs := map[string]taskDTOStruct{}
	allOrderedKeys := []string{}

	for _, row := range rows {
		dto := new(taskDTOStruct)
		byteRow, _ := json.Marshal(row)
		json.Unmarshal(byteRow, &dto)
		allDTOs[dto.Key] = *dto

		if !dto.Fields.Issuetype.Subtask {
			allOrderedKeys = append(allOrderedKeys, dto.Key)
		}
	}

	for _, key := range allOrderedKeys {
		dto := allDTOs[key]

		Task := DomainEntity.ProjectTask{}

		Task.Key = dto.Key
		Task.Summary = dto.Fields.Summary
		Task.Type = dto.Fields.Issuetype.Name
		Task.Status = dto.Fields.Status.Name

		Task.CreatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", dto.Fields.Created.(string))
		Task.UpdatedAt, _ = time.Parse("2006-01-02T15:04:05.000Z", dto.Fields.Updated.(string))

		if dto.Fields.TimeOriginalEstimate > 0 {
			Task.TimeEstimateHours = int(math.Ceil(float64(dto.Fields.TimeOriginalEstimate) / 60 / 60))
		}
		if dto.Fields.TimeSpent > 0 {
			Task.TimeSpentHours = int(math.Ceil(float64(dto.Fields.TimeSpent) / 60 / 60))
		}

		for _, comment := range dto.Fields.Comment.Comments {
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

		for _, subtask := range dto.Fields.Subtasks {
			dtoSub := allDTOs[subtask.Key]
			Task.Assignees = append(Task.Assignees, DomainEntity.ProjectUserExecutor{
				User: DomainEntity.ProjectUser{
					Id:        dtoSub.Fields.Assignee.AccountId,
					Email:     dtoSub.Fields.Assignee.EmailAddress,
					Name:      dtoSub.Fields.Assignee.DisplayName,
					AvatarUrl: dtoSub.Fields.Assignee.AvatarUrls.Image,
				},
				ContributionPerc: 0,
			})
		}

		Task.Reporter = DomainEntity.ProjectUser{
			Id:        dto.Fields.Reporter.AccountId,
			Email:     dto.Fields.Reporter.EmailAddress,
			Name:      dto.Fields.Reporter.DisplayName,
			AvatarUrl: dto.Fields.Reporter.AvatarUrls.Image,
		}

		Task.Epic = DomainEntity.ProjectEpic{
			Key:     dto.Fields.Parent.Key,
			Summary: dto.Fields.Parent.Fields.Summary,
			Status:  dto.Fields.Parent.Fields.Status.Name,
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
