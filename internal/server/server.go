package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/lucas-simao/go-gen/internal/templates"
	"github.com/lucas-simao/go-gen/internal/utils"
)

type Response struct {
	Controller string `json:"controller"`
	UseCase    string `json:"useCase"`
	Model      string `json:"model"`
	Repository string `json:"repository"`
}

//go:embed html/index.html
var indexHTML string

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", getIndex())
	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("internal/server/css"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("internal/server/js"))))

	return mux
}

func getIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("index").Parse(indexHTML)
		if err != nil {
			log.Panic(err)
		}

		projectName := r.URL.Query().Get("projectName")
		serviceName := r.URL.Query().Get("serviceName")
		jsonParsed := r.URL.Query().Get("json")

		t.Execute(w, func() map[string]string {
			if len(projectName) > 0 && len(serviceName) > 0 && len(jsonParsed) > 0 {
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
