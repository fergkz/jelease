package ApplicationUsecase

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	DomainTool "github.com/fergkz/jelease/src/Domain/Tool"
	"github.com/tyler-sommer/stick"
)

type generateReleaseNotesSprintHtmlV2 struct {
	TasksRequestService DomainService.TasksRequestService
	JiraQueryService    DomainService.JiraQueryService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
	TemplateFilename    string
	HostnamePrefix      string
}

func NewGenerateReleaseNotesSprintHtmlV2(
	TasksRequestService DomainService.TasksRequestService,
	JiraQueryService DomainService.JiraQueryService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
	TemplateFilename string,
	HostnamePrefix string,
) *generateReleaseNotesSprintHtmlV2 {
	usecase := new(generateReleaseNotesSprintHtmlV2)
	usecase.TasksRequestService = TasksRequestService
	usecase.JiraQueryService = JiraQueryService
	usecase.ReplaceTeamMembers = ReplaceTeamMembers
	usecase.TemplateFilename = TemplateFilename
	usecase.HostnamePrefix = HostnamePrefix
	return usecase
}

type RenderHtmlDaoTask struct {
	Task     DomainEntity.ProjectTask
	TaskLink string
	Users    map[string]DomainEntity.ProjectUser
}

type RenderHtmlDaoType struct {
	Tasks map[string]RenderHtmlDaoTask
}

type RenderHtmlDaoEpic struct {
	Epic  DomainEntity.ProjectEpic
	Types map[string]RenderHtmlDaoType
}

type RenderHtmlDao struct {
	Board  DomainEntity.ProjectBoard
	Sprint DomainEntity.ProjectSprint
	Epics  map[string]RenderHtmlDaoEpic
	Users  map[string]DomainEntity.ProjectUser
}

func (usecase *generateReleaseNotesSprintHtmlV2) Run(
	boardId DomainEntity.ProjectBoardId,
	sprintIds []DomainEntity.ProjectSprintId,
) string {

	dao := RenderHtmlDao{
		Board:  usecase.JiraQueryService.GetBoardData(boardId),
		Sprint: usecase.JiraQueryService.GetSprints(sprintIds)[0],
		Epics:  make(map[string]RenderHtmlDaoEpic),
		Users:  make(map[string]DomainEntity.ProjectUser),
	}

	tasks, _ := usecase.TasksRequestService.GetTasksFromSprints(sprintIds)

	for _, task := range tasks {

		if _, ok := dao.Epics[task.Epic.Summary]; !ok {
			dao.Epics[task.Epic.Summary] = RenderHtmlDaoEpic{
				Epic:  task.Epic,
				Types: make(map[string]RenderHtmlDaoType),
			}
		}

		if _, ok := dao.Epics[task.Epic.Summary].Types[task.TaskType]; !ok {
			dao.Epics[task.Epic.Summary].Types[task.TaskType] = RenderHtmlDaoType{
				Tasks: make(map[string]RenderHtmlDaoTask),
			}
		}

		users := make(map[string]DomainEntity.ProjectUser)

		for _, user := range task.Assignees {
			if user.User.Email != "" {
				users[user.User.Email] = user.User
			}
		}

		if task.Reporter.Email != "" {
			users[task.Reporter.Email] = task.Reporter
		}

		for _, user := range task.Comments {
			if user.User.Email != "" {
				if !strings.Contains(strings.ToLower(user.User.Name), "automation") {
					users[user.User.Email] = user.User
				}
			}
		}

		if _, ok := dao.Epics[task.Epic.Summary].Types[task.TaskType].Tasks[task.Key]; !ok {
			dao.Epics[task.Epic.Summary].Types[task.TaskType].Tasks[task.Key] = RenderHtmlDaoTask{
				Task:     task,
				TaskLink: usecase.HostnamePrefix + "/browse/" + task.Key,
				Users:    users,
			}
		}

		for _, user := range users {
			dao.Users[user.Email] = user
		}

		if task.Key == "BACKSHIP-4018" {
			DomainTool.Pretty.Save(task, "debug-test.task.json")
		}

	}

	DomainTool.Pretty.Save(dao, "debug-test.dao.json")

	rendered := usecase.RenderHtml(dao)

	return rendered
}

//go:embed GenerateReleaseNotesSprintHtmlV2.twig
var templateRenderTwig []byte

func (usecase *generateReleaseNotesSprintHtmlV2) RenderHtml(dao RenderHtmlDao) string {
	env := stick.New(nil)
	env.Filters["split"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		return strings.Split(stick.CoerceString(val), stick.CoerceString(args[0]))
	}

	env.Filters["imageBase"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		imgdata := usecase.urlImageToData(stick.CoerceString(val))
		return imgdata
	}

	env.Filters["brDate"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		originDt := stick.CoerceString(val)
		layout := "2006-01-02 15:04:05 -0700 UTC"
		parsedTime, err := time.Parse(layout, originDt)
		if err != nil {
			log.Fatal(err)
		}
		return parsedTime.Format("02/01/2006")
	}

	nParams := map[string]stick.Value{}
	nParams["dao"] = dao

	var b bytes.Buffer
	if err := env.Execute(string(templateRenderTwig), &b, nParams); err != nil {
		log.Fatal(err)
	}

	return b.String()
}

func (usecase *generateReleaseNotesSprintHtmlV2) urlImageToData(url string) string {
	img, err := usecase.JiraQueryService.ApiStr(
		"GET",
		url,
	)

	if err != nil {
		log.Fatal(err)
	}

	contentType := http.DetectContentType([]byte(img))

	if strings.Contains(strings.ToLower(string(img)), "</svg>") {
		contentType = "image/svg+xml"
	}

	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(img)
}
