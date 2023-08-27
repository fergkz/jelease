package DomainService

import (
	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
)

type JiraQueryService interface {
	Query(JQL string) []interface{}
	GetSprints(SprintIds []DomainEntity.ProjectSprintId) []DomainEntity.ProjectSprint
	GetBoardData(BoardId DomainEntity.ProjectBoardId) DomainEntity.ProjectBoard
	Api(method string, url string, responseFull interface{}) (err error)
	ApiStr(method string, url string) (response []byte, err error)
}
