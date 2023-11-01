package templates

import (
	"bytes"
	_ "embed"

	goLog "github.com/lucas-simao/golog"

	"github.com/lucas-simao/go-gen-ca/internal/utils"
)

//go:embed controllers.tmpl
var controllersFile string

func GenerateController(serviceName, projectName string) string {
	tmpl, err := utils.InitTemplate("controllers", controllersFile)
	if err != nil {
		goLog.Error(err)
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, map[string]string{
		"ProjectName": projectName,
		"ServiceName": serviceName,
	})
	if err != nil {
		goLog.Error(err)
		return ""
	}

	return b.String()
}
