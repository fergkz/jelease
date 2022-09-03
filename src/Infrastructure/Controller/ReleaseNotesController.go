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

type releaseNotesController struct {
	TasksRequestService DomainService.TasksRequestService
	RenderHtmlService   DomainService.RenderHtmlService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
}

func NewReleaseNotesController(
	TasksRequestService DomainService.TasksRequestService,
	RenderHtmlService DomainService.RenderHtmlService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
) *releaseNotesController {
	controller := new(releaseNotesController)
	controller.TasksRequestService = TasksRequestService
	controller.RenderHtmlService = RenderHtmlService
	controller.ReplaceTeamMembers = ReplaceTeamMembers
	return controller
}

func (controller *releaseNotesController) Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	sprintId := mux.Vars(r)["SprintId"]
	sprintIdInt, _ := strconv.Atoi(sprintId)

	usecase := ApplicationUsecase.NewGenerateReleaseNotesSprintHtml(
		controller.TasksRequestService,
		controller.RenderHtmlService,
		controller.ReplaceTeamMembers,
	)
	html := usecase.Run([]DomainEntity.ProjectSprintId{DomainEntity.ProjectSprintId(sprintIdInt)})

	io.WriteString(w, string(html))
}
