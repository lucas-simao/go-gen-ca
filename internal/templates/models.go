package templates

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/lucas-simao/go-gen-ca/internal/utils"
	goLog "github.com/lucas-simao/golog"
)

//go:embed models.tmpl
var modelsFile string

func GenerateModel(serviceName, jsonParse string) string {
	tmpl, err := utils.InitTemplate("models", modelsFile)
	if err != nil {
		goLog.Error(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, map[string]string{
		"Model":       generateStructFromJson(serviceName, jsonParse),
		"ServiceName": serviceName,
	})
	if err != nil {
		goLog.Error(err)
	}

	return buf.String()
}

func generateStructFromJson(nameParse, jsonParse string) string {
	if len(nameParse) == 0 {
		goLog.Error("should pass name")
		return ""
	}

	b, err := json.MarshalIndent(jsonParse, "", "")
	if err != nil {
		goLog.Error(err)
		return ""
	}

	if !json.Valid(b) {
		goLog.Error("invalid json")
		return ""
	}

	j := fmt.Sprintf("--json=%s", jsonParse)
	n := fmt.Sprintf("--name=%s", nameParse)

	cmd := exec.Command("node", "json-to-go.js", n, j)
	output, err := cmd.Output()
	if err != nil {
		goLog.Error(err)
		return ""
	}

	return string(output)
}
