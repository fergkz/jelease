package InfrastructureController

import (
	"io"
	"net/http"
	"strconv"

	ApplicationUsecase "github.com/fergkz/jelease/src/Application/Usecase"
	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	DomainTool "github.com/fergkz/jelease/src/Domain/Tool"
	"github.com/gorilla/mux"
)

type releaseNotesControllerV2 struct {
	TasksRequestService DomainService.TasksRequestService
	RenderHtmlService   DomainService.RenderHtmlService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
}

func NewReleaseNotesControllerV2(
	TasksRequestService DomainService.TasksRequestService,
	RenderHtmlService DomainService.RenderHtmlService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
) *releaseNotesControllerV2 {
	controller := new(releaseNotesControllerV2)
	controller.TasksRequestService = TasksRequestService
	controller.RenderHtmlService = RenderHtmlService
	controller.ReplaceTeamMembers = ReplaceTeamMembers
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

	DomainTool.Pretty.Println("BOARD ID", boardIdInt)
	DomainTool.Pretty.Println("SPRINT ID", sprintId)
	// os.Exit(0)

	usecase := ApplicationUsecase.NewGenerateReleaseNotesSprintHtmlV2(
		controller.TasksRequestService,
		controller.RenderHtmlService,
		controller.ReplaceTeamMembers,
	)
	html := usecase.Run([]DomainEntity.ProjectSprintId{DomainEntity.ProjectSprintId(sprintIdInt)})

	io.WriteString(w, string(html))
}
