package main

import (
	"fmt"
	"testing"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
	DomainTool "github.com/fergkz/jelease/src/Domain/Tool"
	InfrastructureService "github.com/fergkz/jelease/src/Infrastructure/Service"
	"github.com/spf13/viper"
)

type TestConfig struct {
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

func getTestConfig() *TestConfig {
	config := new(TestConfig)

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	viper.Unmarshal(config)

	return config
}

func TestBoardAndSprint(t *testing.T) {
	config := getTestConfig()

	service := InfrastructureService.NewJiraQueryService(
		config.Jira.Username,
		config.Jira.AccessToken,
		config.Jira.Hostname,
	)

	boardId := DomainEntity.ProjectBoardId(502)
	boardData := service.GetBoardData(boardId)
	DomainTool.Pretty.Save(boardData, "debug-test.BoardData.json")

	sprintId := DomainEntity.ProjectSprintId(3319)
	sprintData := service.GetSprints([]DomainEntity.ProjectSprintId{sprintId})
	DomainTool.Pretty.Save(sprintData, "debug-test.SprintData.json")
}
