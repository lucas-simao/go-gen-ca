package templates

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/lucas-simao/go-gen/internal/utils"
)

//go:embed controllers.tmpl
var controllersFile string

func GenerateController(serviceName, projectName string) string {
	tmpl, err := utils.InitTemplate("controllers", controllersFile)
	if err != nil {
		log.Panic(err)
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, map[string]string{
		"ProjectName": projectName,
		"ServiceName": serviceName,
	})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return b.String()
}
