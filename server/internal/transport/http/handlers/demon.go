package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/DimaKropachev/calculate-web-server/server/internal/entities"
	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
)

type DemonHandler struct {
	rq *entities.ResultsQueue
	tq *entities.TasksQueue
}

func NewDemonHandler(tq *entities.TasksQueue, rq *entities.ResultsQueue) *DemonHandler {
	return &DemonHandler{
		tq: tq,
		rq: rq,
	}
}
func (h *DemonHandler) GiveTask(w http.ResponseWriter, r *http.Request) {
	var task *models.FinalTask

	for {
		task = h.tq.Get()
		if task == nil {
			return
		} else {
			break
		}
	}

	data, err := json.Marshal(task)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (h *DemonHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res := &models.Result{}
	err = json.Unmarshal(data, res)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.rq.Add(res)
}
