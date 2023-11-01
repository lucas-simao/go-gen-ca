package templates

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/lucas-simao/go-gen/internal/utils"
)

//go:embed models.tmpl
var modelsFile string

func GenerateModel(serviceName, jsonParse string) string {
	structGenerated, err := generateStructFromJson(serviceName, jsonParse)
	if err != nil {
		log.Panic(err)
	}

	tmpl, err := utils.InitTemplate("models", modelsFile)
	if err != nil {
		log.Panic(err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, map[string]string{
		"Model":       structGenerated,
		"ServiceName": serviceName,
	})
	if err != nil {
		log.Panic(err)
	}

	return buf.String()
}

func generateStructFromJson(nameParse, jsonParse string) (string, error) {
	if len(nameParse) == 0 {
		return "", errors.New("should pass name")
	}

	b, err := json.MarshalIndent(jsonParse, "", "")
	if err != nil {
		return "", err
	}

	if !json.Valid(b) {
		return "", errors.New("invalid json")
	}

	j := fmt.Sprintf("--json=%s", jsonParse)
	n := fmt.Sprintf("--name=%s", nameParse)

	cmd := exec.Command("node", "json-to-go.js", n, j)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
