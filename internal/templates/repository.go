package templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/lucas-simao/go-gen-ca/internal/utils"
	goLog "github.com/lucas-simao/golog"
)

var re = regexp.MustCompile(`"([^"]+)"`)

//go:embed repository.tmpl
var repositoryFile string

func GenerateRepository(model, serviceName, projectName string) string {
	tmpl, err := utils.InitTemplate("repository", repositoryFile)
	if err != nil {
		goLog.Error(err)
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, map[string]string{
		"ProjectName": projectName,
		"ServiceName": serviceName,
		"InsertQuery": GenerateInsertQuery(model, serviceName),
		"InsertQueryArgs": func() string {
			var args string

			fields := GenerateFields(model)

			for i, v := range fields {
				args += fmt.Sprintf("%s.%s", utils.ToLower(serviceName), v)

				if i+1 < len(fields) {
					args += ", "
				}
			}
			return args
		}(),
		"GetQuery":    GenerateGetQuery(model, serviceName),
		"UpdateQuery": GenerateUpdateQuery(model, serviceName),
		"DeleteQuery": GenerateDeleteQuery(serviceName),
	})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return b.String()
}

func GenerateFields(model string) []string {
	var fieldsRegex []string

	for _, v := range re.FindAllString(model, -1) {
		v = strings.ReplaceAll(v, "\"", "")
		if !strings.Contains(v, "id") {
			v := strings.ToUpper(v[:1]) + v[1:]

			fieldsRegex = append(fieldsRegex, strings.ReplaceAll(v, "Id", "ID"))
		}
	}

	return fieldsRegex
}

func GenerateInsertQuery(model, serviceName string) string {
	var fieldsRegex []string

	for _, v := range re.FindAllString(model, -1) {
		if !strings.Contains(v, "id") {
			fieldsRegex = append(fieldsRegex, utils.ToSnakeCase(v))
		}
	}

	var fields string = strings.Join(fieldsRegex, ", ")

	fields = strings.ReplaceAll(fields, "\"", "")

	insert := fmt.Sprintf("INSERT INTO %s(%s)", strings.ToLower(serviceName), fields)

	var insertArgs string

	for i := 1; i <= len(fieldsRegex); i++ {
		insertArgs += fmt.Sprintf("$%v", i)

		if i < len(fieldsRegex) {
			insertArgs += ", "
		}
	}

	insert = fmt.Sprintf("%s\nVALUES(%s)\nRETURNING id", insert, insertArgs)

	return insert
}

func GenerateGetQuery(model, serviceName string) string {
	var fieldsRegex []string

	for _, v := range re.FindAllString(model, -1) {
		if !strings.Contains(v, "id") {
			fieldsRegex = append(fieldsRegex, utils.ToSnakeCase(v))
		}
	}

	var fields string = strings.Join(fieldsRegex, ", ")

	fields = strings.ReplaceAll(fields, "\"", "")

	selectQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", fields, strings.ToLower(serviceName))

	return selectQuery
}

func GenerateUpdateQuery(model, serviceName string) string {
	var fieldsRegex []string

	for _, v := range re.FindAllString(model, -1) {
		if !strings.Contains(v, "id") {
			fieldsRegex = append(fieldsRegex, utils.ToSnakeCase(v))
		}
	}

	update := fmt.Sprintf("UPDATE %s", strings.ToLower(serviceName))

	var updateArgs string

	for i := 0; i < len(fieldsRegex); i++ {
		updateArgs += fmt.Sprintf("%s = $%v", strings.ReplaceAll(fieldsRegex[i], "\"", ""), i+1)

		if i+1 < len(fieldsRegex) {
			updateArgs += ", "
		}
	}

	update = fmt.Sprintf("%s\nSET %s\nWHERE id = $%d", update, updateArgs, len(fieldsRegex)+1)

	return update
}

func GenerateDeleteQuery(serviceName string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = $1", strings.ToLower(serviceName))
}
