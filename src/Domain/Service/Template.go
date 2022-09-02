package DomainService

type TemplateParseServiceLink struct {
	Title     string
	PublicUrl string
}

type TemplateParseServiceTeamMember struct {
	DisplayName    string
	FullName       string
	Office         string
	PublicImageUrl string
}

type TemplateParseServiceSystemUpdatedNote struct {
	Title       string
	Type        string // news, improvement, bug, internal, note, spike
	Status      string // executing, done
	Assignees   []TemplateParseServiceTeamMember
	Links       []TemplateParseServiceLink
	TextMessage string
}

type TemplateParseServiceSystemUpdated struct {
	SystemName string
	Notes      []TemplateParseServiceSystemUpdatedNote
}

type TemplateParseServiceHtmlRendered struct {
	HtmlContent string
}

type TemplateParseService interface {
	Parse([]TemplateParseServiceTeamMember, []TemplateParseServiceSystemUpdated) TemplateParseServiceHtmlRendered
}
