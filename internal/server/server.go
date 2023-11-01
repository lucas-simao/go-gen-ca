package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	goLog "github.com/lucas-simao/golog"

	"github.com/lucas-simao/go-gen-ca/internal/templates"
	"github.com/lucas-simao/go-gen-ca/internal/utils"
)

type Response struct {
	Controller string `json:"controller"`
	UseCase    string `json:"useCase"`
	Model      string `json:"model"`
	Repository string `json:"repository"`
}

//go:embed html/index.html
var indexHTML string

//go:embed css/index.css
var indexCSS string

//go:embed css/prism.css
var prismCSS string

//go:embed js/prism.js
var prismJS string

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", getIndex())
	mux.Handle("/css/index.css", getIndexCSS())
	mux.Handle("/css/prism.css", getPrismCSS())
	mux.Handle("/js/prism.js", getPrismJS())

	return mux
}

func getIndexCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, indexCSS)
	}
}

func getPrismCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fmt.Fprint(w, prismCSS)
	}
}

func getPrismJS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/text")
		fmt.Fprint(w, prismJS)
	}
}

func getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").Parse(indexHTML)
		if err != nil {
			goLog.Error(err)
		}

		projectName := r.URL.Query().Get("projectName")
		serviceName := r.URL.Query().Get("serviceName")
		jsonParsed := r.URL.Query().Get("json")

		t.Execute(w, func() map[string]string {
			if len(projectName) > 0 && len(serviceName) > 0 && len(jsonParsed) > 0 && json.Valid([]byte(jsonParsed)) {
				goLog.Info("new request", map[string]interface{}{
					"projectName": projectName,
					"serviceName": serviceName,
					"jsonParsed":  jsonParsed,
				})

				model := templates.GenerateModel(serviceName, jsonParsed)

				return map[string]string{
					"ModelPath":      fmt.Sprintf("/%s/src/domain/models/%s.go", utils.ToLower(projectName), utils.ToLower(serviceName)),
					"Model":          model,
					"ControllerPath": fmt.Sprintf("/%s/src/application/controllers/%s_controller.go", utils.ToLower(projectName), utils.ToLower(serviceName)),
					"Controller":     templates.GenerateController(serviceName, projectName),
					"UseCasePath":    fmt.Sprintf("/%s/src/domain/usecases/%s_usecases.go", utils.ToLower(projectName), utils.ToLower(serviceName)),
					"UseCase":        templates.GenerateUsecases(serviceName, projectName),
					"RepositoryPath": fmt.Sprintf("/%s/src/infra/repositories/%s_repository.go", utils.ToLower(projectName), utils.ToLower(serviceName)),
					"Repository":     templates.GenerateRepository(model, serviceName, projectName),
				}
			}
			return map[string]string{}
		}())
	}
}
