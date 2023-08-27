package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	ApplicationUsecase "github.com/fergkz/jelease/src/Application/Usecase"
	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	DomainTool "github.com/fergkz/jelease/src/Domain/Tool"
	InfrastructureService "github.com/fergkz/jelease/src/Infrastructure/Service"
	"github.com/spf13/viper"
)

type testInstance struct {
	Config struct {
		Jira struct {
			Username     string
			AccessToken  string
			Hostname     string
			CustomFields []struct {
				TaskObjective string
				TaskType      string
			}
		}
		Server struct {
			Port string
		}
		Team []struct {
			Email       string
			DisplayName string
			Office      string
		}
	}
	JiraQueryService DomainService.JiraQueryService
	JiraTaskService  DomainService.TasksRequestService
	ReplaceMembers   map[string]DomainService.RenderHtmlServiceTeamMember
}

func getTestInstance() *testInstance {
	instance := new(testInstance)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	viper.Unmarshal(&instance.Config)

	instance.JiraQueryService = InfrastructureService.NewJiraQueryService(
		instance.Config.Jira.Username,
		instance.Config.Jira.AccessToken,
		instance.Config.Jira.Hostname,
	)

	instance.JiraTaskService = InfrastructureService.NewJiraTaskService(
		instance.Config.Jira.Username,
		instance.Config.Jira.AccessToken,
		instance.Config.Jira.Hostname,
		instance.Config.Jira.CustomFields[0].TaskObjective,
		instance.Config.Jira.CustomFields[0].TaskType,
	)

	instance.ReplaceMembers = map[string]DomainService.RenderHtmlServiceTeamMember{}
	for _, m := range instance.Config.Team {
		instance.ReplaceMembers[m.Email] = DomainService.RenderHtmlServiceTeamMember{
			DisplayName: m.DisplayName,
			FullName:    "",
			Office:      m.Office,
		}
	}

	return instance
}

func _TestBoardAndSprint(t *testing.T) {
	instance := getTestInstance()

	boardId := DomainEntity.ProjectBoardId(502)
	boardData := instance.JiraQueryService.GetBoardData(boardId)
	DomainTool.Pretty.Save(boardData, "debug-test.BoardData.json")

	sprintId := DomainEntity.ProjectSprintId(3319)
	sprintData := instance.JiraQueryService.GetSprints([]DomainEntity.ProjectSprintId{sprintId})
	DomainTool.Pretty.Save(sprintData, "debug-test.SprintData.json")
}

func TestHtmlV2(t *testing.T) {
	instance := getTestInstance()

	usecase := ApplicationUsecase.NewGenerateReleaseNotesSprintHtmlV2(
		instance.JiraTaskService,
		instance.JiraQueryService,
		instance.ReplaceMembers,
		"templatev2.twig",
		instance.Config.Jira.Hostname,
	)
	html := usecase.Run(
		DomainEntity.ProjectBoardId(502),
		[]DomainEntity.ProjectSprintId{DomainEntity.ProjectSprintId(3319)},
	)
	ioutil.WriteFile("debug-test.HTMLv2.html", []byte(html), 0644)
}
