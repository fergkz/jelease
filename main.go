package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/fcgi"

	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	InfrastructureController "github.com/fergkz/jelease/src/Infrastructure/Controller"
	InfrastructureService "github.com/fergkz/jelease/src/Infrastructure/Service"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	config := new(struct {
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
	})

	viper.Unmarshal(config)

	replaceMembers := map[string]DomainService.RenderHtmlServiceTeamMember{}
	for _, m := range config.Team {
		replaceMembers[m.Email] = DomainService.RenderHtmlServiceTeamMember{
			DisplayName: m.DisplayName,
			FullName:    "",
			Office:      m.Office,
		}
	}

	jiraTaskService := InfrastructureService.NewJiraTaskService(
		config.Jira.Username,
		config.Jira.AccessToken,
		config.Jira.Hostname,
		config.Jira.CustomFields[0].TaskObjective,
		config.Jira.CustomFields[0].TaskType,
	)

	jiraQueryService := InfrastructureService.NewJiraQueryService(
		config.Jira.Username,
		config.Jira.AccessToken,
		config.Jira.Hostname,
	)

	router := mux.NewRouter()
	releaseNotesController := InfrastructureController.NewReleaseNotesController(
		jiraTaskService,
		InfrastructureService.NewRenderHtmlService("template.twig"),
		replaceMembers,
	)
	router.HandleFunc("/sprint/{SprintId:[0-9]+}", releaseNotesController.Get).Methods("GET")

	releaseNotesControllerV2 := InfrastructureController.NewReleaseNotesControllerV2(
		jiraTaskService,
		jiraQueryService,
		replaceMembers,
		"templatev2.twig",
		config.Jira.Hostname,
	)
	router.HandleFunc("/v2/board/{BoardId:[0-9]+}/sprint/{SprintId:[0-9]+}", releaseNotesControllerV2.Get).Methods("GET")

	if viper.GetString("server.method") == "http" {
		fmt.Printf("Server started at port %s", config.Server.Port)
		log.Fatal(http.ListenAndServe("127.0.0.1:"+config.Server.Port, router))
	} else {
		fcgi.Serve(nil, router)
	}

}
