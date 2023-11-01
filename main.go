package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucas-simao/go-gen/internal/templates"
)

type Request struct {
	ServiceName string `json:"serviceName"`
	ProjectName string `json:"projectName"`
	Json        string `json:"json"`
}

type Response struct {
	Controller string `json:"controller"`
	UseCase    string `json:"useCase"`
	Model      string `json:"model"`
	Repository string `json:"repository"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p *Request

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		model := templates.GenerateModel(p.ServiceName, p.Json)

		var resp = Response{
			Controller: templates.GenerateController(p.ServiceName, p.ProjectName),
			UseCase:    templates.GenerateUsecases(p.ServiceName, p.ProjectName),
			Model:      model,
			Repository: templates.GenerateRepository(model, p.ServiceName, p.ProjectName),
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Panic(err)
	}
}
