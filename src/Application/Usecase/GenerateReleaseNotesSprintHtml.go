package ApplicationUsecase

import (
	"strings"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
)

type generateReleaseNotesSprintHtml struct {
	TasksRequestService DomainService.TasksRequestService
	RenderHtmlService   DomainService.RenderHtmlService
	ReplaceTeamMembers  map[string]DomainService.RenderHtmlServiceTeamMember
}

func NewGenerateReleaseNotesSprintHtml(
	TasksRequestService DomainService.TasksRequestService,
	RenderHtmlService DomainService.RenderHtmlService,
	ReplaceTeamMembers map[string]DomainService.RenderHtmlServiceTeamMember,
) *generateReleaseNotesSprintHtml {
	usecase := new(generateReleaseNotesSprintHtml)
	usecase.TasksRequestService = TasksRequestService
	usecase.RenderHtmlService = RenderHtmlService
	usecase.ReplaceTeamMembers = ReplaceTeamMembers
	return usecase
}

func (usecase *generateReleaseNotesSprintHtml) Run(sprintIds []DomainEntity.ProjectSprintId) string {

	tasks := []DomainEntity.ProjectTask{}
	// DomainTool.Pretty.SetCache("cache-tasks.json", &tasks)
	// DomainTool.Pretty.GetCache("cache-tasks.json", &tasks)
	// if len(tasks) == 0 {
	tasks, _ = usecase.TasksRequestService.GetTasksFromSprints(sprintIds)
	// DomainTool.Pretty.SetCache("cache-tasks.json", &tasks)
	// }

	members := []DomainService.RenderHtmlServiceTeamMember{}
	systems := []DomainService.RenderHtmlServiceSystemUpdated{}

	systemNotes := map[string][]DomainService.RenderHtmlServiceSystemUpdatedNote{}
	membersMap := map[string]DomainService.RenderHtmlServiceTeamMember{}
	taskNoteIsset := map[string]bool{}

	for _, task := range tasks {
		for _, comment := range task.Comments {

			comment.Body = strings.ReplaceAll(comment.Body, "*RELEASE NOTES*", "RELEASE NOTES")
			comment.Body = strings.ReplaceAll(comment.Body, "*Title*", "Title")
			comment.Body = strings.ReplaceAll(comment.Body, "*Description*", "Description")
			comment.Body = strings.ReplaceAll(comment.Body, "*Type*", "Type")

			if strings.HasPrefix(comment.Body, "RELEASE NOTES") {

				cTitle, cDescr, cType, cSystem := usecase.parseComment(comment.Body)

				note := usecase.taskToNote(task, cTitle, cDescr, cType, cSystem)

				for _, m := range note.Assignees {
					membersMap[m.Email] = m
				}

				if cSystem != "" {
					systemNotes[cSystem] = append(systemNotes[cSystem], note)
				} else {
					systemNotes[task.Epic.Summary] = append(systemNotes[task.Epic.Summary], note)
				}

				taskNoteIsset[task.Key] = true
			}
		}

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

func (usecase *generateReleaseNotesSprintHtml) taskToNote(task DomainEntity.ProjectTask, cTitle, cDescr, cType, cSystem string) (note DomainService.RenderHtmlServiceSystemUpdatedNote) {
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

func (usecase *generateReleaseNotesSprintHtml) parseComment(comment string) (title string, description string, ctype string, csystem string) {

	sc := strings.Split(comment, "\n")

	for _, t := range sc {
		if strings.HasPrefix(t, "Title:") {
			sctitle := strings.Split(t, "Title:")
			title = strings.Trim(sctitle[1], " ")
		}
		if strings.HasPrefix(t, "Description:") {
			scdescr := strings.Split(t, "Description:")
			description = strings.Trim(scdescr[1], " ")
		}
		if strings.HasPrefix(t, "Type:") {
			sctype := strings.Split(t, "Type:")
			ctype = strings.Trim(sctype[1], " ")
		}
		if strings.HasPrefix(t, "System:") {
			scsystem := strings.Split(t, "System:")
			csystem = strings.Trim(scsystem[1], " ")
		}
	}

	return title, description, ctype, csystem
}
