package InfrastructureController

import (
	"io"
	"net/http"
	"strconv"

	ApplicationUsecase "github.com/fergkz/jelease/src/Application/Usecase"
	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	"github.com/gorilla/mux"
)

type releaseNotesControllerV2 struct {
	TasksRequestService DomainService.TasksRequestService
	JiraQueryService    DomainService.JiraQueryService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
	TemplateFilename    string
	HostnamePrefix      string
}

func NewReleaseNotesControllerV2(
	TasksRequestService DomainService.TasksRequestService,
	JiraQueryService DomainService.JiraQueryService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
	TemplateFilename string,
	HostnamePrefix string,
) *releaseNotesControllerV2 {
	controller := new(releaseNotesControllerV2)
	controller.TasksRequestService = TasksRequestService
	controller.JiraQueryService = JiraQueryService
	controller.ReplaceTeamMembers = ReplaceTeamMembers
	controller.TemplateFilename = TemplateFilename
	controller.HostnamePrefix = HostnamePrefix
	return controller
}

func (controller *releaseNotesControllerV2) Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	boardId := mux.Vars(r)["BoardId"]
	boardIdInt, _ := strconv.Atoi(boardId)
	sprintId := mux.Vars(r)["SprintId"]
	sprintIdInt, _ := strconv.Atoi(sprintId)

	usecase := ApplicationUsecase.NewGenerateReleaseNotesSprintHtmlV2(
		controller.TasksRequestService,
		controller.JiraQueryService,
		controller.ReplaceTeamMembers,
		controller.TemplateFilename,
		controller.HostnamePrefix,
	)
	html := usecase.Run(
		DomainEntity.ProjectBoardId(boardIdInt),
		[]DomainEntity.ProjectSprintId{DomainEntity.ProjectSprintId(sprintIdInt)},
	)

	io.WriteString(w, string(html))
}
