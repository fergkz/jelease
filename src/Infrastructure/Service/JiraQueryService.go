package InfrastructureService

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	DomainEntity "github.com/fergkz/jelease/src/Domain/Entity"
)

type jiraQueryService struct {
	Username    string
	AccessToken string
	Hostname    string
}

func NewJiraQueryService(
	Username string,
	AccessToken string,
	Hostname string,
) *jiraQueryService {
	service := new(jiraQueryService)
	service.Username = Username
	service.AccessToken = AccessToken
	service.Hostname = Hostname
	return service
}

func (service jiraQueryService) Query(JQL string) (rows []interface{}) {

	JQL = url.QueryEscape(JQL)

	var responseFull struct {
		Issues []interface{}
	}

	maxResults := 100
	currentIndex := 0

	for {
		url := service.Hostname + "/rest/api/2/search?jql=" + JQL + "&maxResults=" + strconv.Itoa(maxResults) + "&startAt=" + strconv.Itoa(currentIndex) + "&expand=changelog&fields=*all"

		service.Api("GET", url, &responseFull)

		if len(responseFull.Issues) == 0 {
			break
		}

		rows = append(rows, responseFull.Issues...)

		currentIndex = currentIndex + maxResults
	}

	return rows
}

func (service jiraQueryService) GetSprints(SprintIds []DomainEntity.ProjectSprintId) (rows []DomainEntity.ProjectSprint) {

	for _, i := range SprintIds {
		url := service.Hostname + "/rest/agile/1.0/sprint/" + strconv.Itoa(int(i))
		var response struct {
			StartDate     interface{}
			EndDate       interface{}
			CompleteDate  interface{}
			Goal          string
			Id            int
			Name          string
			OriginBoardId int
			State         string
		}
		service.Api("GET", url, &response)

		sprint := DomainEntity.ProjectSprint{
			Id:            DomainEntity.ProjectSprintId(response.Id),
			Name:          response.Name,
			State:         response.State,
			OriginBoardId: response.OriginBoardId,
			Goal:          response.Goal,
		}

		if response.StartDate != nil {
			sprint.StartDate, _ = time.Parse("2006-01-02T15:04:05.000Z", response.StartDate.(string))
		}
		if response.EndDate != nil {
			sprint.EndDate, _ = time.Parse("2006-01-02T15:04:05.000Z", response.EndDate.(string))
		}
		if response.CompleteDate != nil {
			sprint.CompleteDate, _ = time.Parse("2006-01-02T15:04:05.000Z", response.CompleteDate.(string))
		}

		rows = append(rows, sprint)
	}

	return rows
}

func (service jiraQueryService) GetBoardData(BoardId DomainEntity.ProjectBoardId) DomainEntity.ProjectBoard {
	url := service.Hostname + "/rest/agile/1.0/board/" + strconv.Itoa(int(BoardId))

	response := new(DomainEntity.ProjectBoard)
	service.Api("GET", url, &response)

	return *response
}

func (service jiraQueryService) Test() {

}

func (service jiraQueryService) Api(method string, url string, responseFull interface{}) (err error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:http.NewRequest\n", err)
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(service.Username + ":" + service.AccessToken))

	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err == nil && res.StatusCode != 200 {
		err = errors.New("STATUS RESPONSE: " + strconv.Itoa(res.StatusCode) + " - " + http.StatusText(res.StatusCode))
		log.Panicln("ERROR: ", err)
	}
	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:client.Do\n", err)
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:ioutil.ReadAll\n", err)
		return err
	}

	json.Unmarshal(body, &responseFull)

	return nil
}

func (service jiraQueryService) ApiStr(method string, url string) (response []byte, err error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:http.NewRequest\n", err)
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(service.Username + ":" + service.AccessToken))

	req.Header.Add("Authorization", "Basic "+auth)

	res, err := client.Do(req)
	if err == nil && res.StatusCode != 200 {
		err = errors.New("STATUS RESPONSE: " + strconv.Itoa(res.StatusCode) + " - " + http.StatusText(res.StatusCode))
		log.Panicln("ERROR: ", err)
	}
	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:client.Do\n", err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicln("TaskerInfrastructureService.TaskService.jiraApi", "ERROR:ioutil.ReadAll\n", err)
	}

	return body, nil
}
