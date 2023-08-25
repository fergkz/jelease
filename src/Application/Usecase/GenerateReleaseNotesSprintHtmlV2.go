package ApplicationUsecase

import (
	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
)

type generateReleaseNotesSprintHtmlV2 struct {
	TasksRequestService DomainService.TasksRequestService
	RenderHtmlService   DomainService.RenderHtmlService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
}

func NewGenerateReleaseNotesSprintHtmlV2(
	TasksRequestService DomainService.TasksRequestService,
	RenderHtmlService DomainService.RenderHtmlService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
) *generateReleaseNotesSprintHtmlV2 {
	usecase := new(generateReleaseNotesSprintHtmlV2)
	usecase.TasksRequestService = TasksRequestService
	usecase.RenderHtmlService = RenderHtmlService
	usecase.ReplaceTeamMembers = ReplaceTeamMembers
	return usecase
}

func (usecase *generateReleaseNotesSprintHtmlV2) Run(sprintIds []DomainEntity.ProjectSprintId) string {

	tasks := []DomainEntity.ProjectTask{}
	tasks, _ = usecase.TasksRequestService.GetTasksFromSprints(sprintIds)

	members := []DomainService.RenderHtmlServiceTeamMember{}
	systems := []DomainService.RenderHtmlServiceSystemUpdated{}

	systemNotes := map[string][]DomainService.RenderHtmlServiceSystemUpdatedNote{}
	membersMap := map[string]DomainService.RenderHtmlServiceTeamMember{}
	taskNoteIsset := map[string]bool{}

	for _, task := range tasks {

		note := usecase.taskToNote(task, task.Summary, task.Objective, task.TaskType, "")

		for _, m := range note.Assignees {
			membersMap[m.Email] = m
		}

		systemNotes[task.Epic.Summary] = append(systemNotes[task.Epic.Summary], note)

		taskNoteIsset[task.Key] = true

		if _, ok := taskNoteIsset[task.Key]; !ok {
			note := usecase.taskToNote(task, "", "", "", "")
			for _, m := range note.Assignees {
				membersMap[m.Email] = m
			}
			systemNotes[task.Epic.Summary] = append(systemNotes[task.Epic.Summary], note)
		}
	}

	for system, notes := range systemNotes {
		systems = append(systems, DomainService.RenderHtmlServiceSystemUpdated{
			SystemName: system,
			Notes:      notes,
		})
	}

	for _, member := range membersMap {
		members = append(members, member)
	}

	rendered := usecase.RenderHtmlService.Parse(members, systems)

	return rendered.HtmlContent
}

func (usecase *generateReleaseNotesSprintHtmlV2) taskToNote(task DomainEntity.ProjectTask, cTitle, cDescr, cType, cSystem string) (note DomainService.RenderHtmlServiceSystemUpdatedNote) {
	note.Title = task.Summary
	if cTitle != "" {
		note.Title = cTitle
	}
	note.Type = task.Type
	if cType != "" {
		note.Type = cType
	}
	note.TextMessage = ""
	if cDescr != "" {
		note.TextMessage = cDescr
	}

	note.Status = task.Status

	assigneeIsset := map[string]bool{}
	note.Assignees = []DomainService.RenderHtmlServiceTeamMember{}
	for _, assignee := range task.Assignees {
		if !assigneeIsset[assignee.User.Email] && assignee.User.Email != "" {

			newDisplayName := assignee.User.Name
			newFullName := assignee.User.Name
			newOffice := assignee.User.Email
			newPublicImageUrl := assignee.User.AvatarUrl

			if data, ok := usecase.ReplaceTeamMembers[assignee.User.Email]; ok {
				if data.DisplayName != "" {
					newDisplayName = data.DisplayName
				}
				if data.Office != "" {
					newOffice = data.Office
				}
			}

			member := DomainService.RenderHtmlServiceTeamMember{
				DisplayName:    newDisplayName,
				FullName:       newFullName,
				Office:         newOffice,
				PublicImageUrl: newPublicImageUrl,
				Email:          assignee.User.Email,
			}
			note.Assignees = append(note.Assignees, member)
			assigneeIsset[assignee.User.Email] = true
		}
	}

	note.Links = []DomainService.RenderHtmlServiceLink{}
	note.Links = append(note.Links, DomainService.RenderHtmlServiceLink{
		Title:     task.Key,
		PublicUrl: task.PublicHtmlUrl,
	})

	return note
}
