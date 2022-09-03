package DomainEntity

import (
	"time"
)

type ProjectUserExecutor struct {
	User             ProjectUser
	ContributionPerc float32
}

type ProjectEpic struct {
	Key     string
	Summary string
	Status  string
}

type ProjectUser struct {
	Id        string
	Email     string
	Name      string
	AvatarUrl string
}

type ProjectComment struct {
	Body      string
	CreatedAt time.Time
	Public    bool
	User      ProjectUser
}

type ProjectSprintId int

type ProjectSprint struct {
	Id           ProjectSprintId
	Name         string
	State        string
	StartDate    time.Time
	EndDate      time.Time
	CompleteDate time.Time
}

type ProjectTask struct {
	Key               string
	Summary           string
	Type              string
	Status            string
	Assignees         []ProjectUserExecutor
	Reporter          ProjectUser
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Planned           bool
	TimeEstimateHours int
	TimeSpentHours    int
	Epic              ProjectEpic
	Comments          []ProjectComment
	PublicHtmlUrl     string
}
