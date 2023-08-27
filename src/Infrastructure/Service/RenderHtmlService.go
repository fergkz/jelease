package InfrastructureService

import (
	"bytes"
	"log"
	"path/filepath"

	DomainService "github.com/fergkz/jelease/src/Domain/Service"
	"github.com/tyler-sommer/stick"
)

type renderHtmlService struct {
	templateFilename string
}

func NewRenderHtmlService(
	templateFilename string,
) *renderHtmlService {
	service := new(renderHtmlService)
	service.templateFilename = templateFilename
	return service
}

func (service renderHtmlService) Parse(members []DomainService.RenderHtmlServiceTeamMember, systems []DomainService.RenderHtmlServiceSystemUpdated) DomainService.RenderHtmlServiceHtmlRendered {

	fullPath, _ := filepath.Abs(service.templateFilename)

	dirname, filename := filepath.Split(fullPath)

	render := stick.New(stick.NewFilesystemLoader(dirname))

	nParams := map[string]stick.Value{}
	nParams["members"] = members
	nParams["systems"] = systems

	mapSysTypeNote := map[string]map[string][]DomainService.RenderHtmlServiceSystemUpdatedNote{}
	mapSysKeys := map[string]DomainService.RenderHtmlServiceSystemUpdated{}

	for _, system := range systems {
		mapSysKeys[system.SystemName] = system

		if len(mapSysTypeNote[system.SystemName]) == 0 {
			mapSysTypeNote[system.SystemName] = map[string][]DomainService.RenderHtmlServiceSystemUpdatedNote{}
		}

		for _, note := range system.Notes {
			if len(mapSysTypeNote[system.SystemName][note.Type]) == 0 {
				mapSysTypeNote[system.SystemName][note.Type] = []DomainService.RenderHtmlServiceSystemUpdatedNote{}
			}

			mapSysTypeNote[system.SystemName][note.Type] = append(mapSysTypeNote[system.SystemName][note.Type], note)
		}
	}

	nParams["mapSysTypeNote"] = mapSysTypeNote
	nParams["mapSysKeys"] = mapSysKeys

	var b bytes.Buffer
	err := render.Execute(filename, &b, nParams)
	if err != nil {
		log.Fatalln(err)
	}

	rendered := DomainService.RenderHtmlServiceHtmlRendered{}
	rendered.HtmlContent = b.String()

	return rendered
}
