package DomainService

type RenderHtmlServiceLink struct {
	Title     string
	PublicUrl string
}

type RenderHtmlServiceTeamMember struct {
	DisplayName    string
	FullName       string
	Office         string
	PublicImageUrl string
	Email          string
}

type RenderHtmlServiceSystemUpdatedNote struct {
	Title       string
	Type        string
	TextMessage string
	Status      string
	Assignees   []RenderHtmlServiceTeamMember
	Links       []RenderHtmlServiceLink
}

type RenderHtmlServiceSystemUpdated struct {
	SystemName string
	Notes      []RenderHtmlServiceSystemUpdatedNote
}

type RenderHtmlServiceHtmlRendered struct {
	HtmlContent string
}

type RenderHtmlService interface {
	Parse([]RenderHtmlServiceTeamMember, []RenderHtmlServiceSystemUpdated) RenderHtmlServiceHtmlRendered
}
