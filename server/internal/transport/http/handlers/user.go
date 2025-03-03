package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
	"github.com/gorilla/mux"
)

type Service interface {
	Add(expression string) string
	Get(Id string) (*models.Expression, error)
	GetAll() []*models.Expression
}

type UserHandler struct {
	service Service
	idChan  chan string
}

func NewUserHandler(service Service, idChan chan string) *UserHandler {
	return &UserHandler{
		service: service,
		idChan:  idChan,
	}
}

func (h *UserHandler) AddExpr(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	expression := r.FormValue("expression")

	id := h.service.Add(expression)

	h.idChan <- id

	w.WriteHeader(201)

	tmpl := template.Must(template.ParseFiles("./server/web/id.html"))
	err = tmpl.Execute(w, struct {
		Expression string
		Id         string
	}{
		Expression: expression,
		Id:         id,
	})
	if err != nil {
		resp := &models.AddResp{
			ID: id,
		}
		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.Write(data)
	}
}

func (h *UserHandler) GetAllExprs(w http.ResponseWriter, r *http.Request) {
	result := h.service.GetAll()

	w.WriteHeader(200)

	tmpl := template.Must(template.ParseFiles("./server/web/expressions.html"))
	err := tmpl.Execute(w, result)
	if err != nil {
		_, err := json.Marshal(struct {
			Expr []*models.Expression `json:"expressions"`
		}{
			Expr: result,
		})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}

func (h *UserHandler) GetExprById(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	res, err := h.service.Get(id)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.WriteHeader(200)

	tmpl := template.Must(template.ParseFiles("./server/web/expression.html"))
	err = tmpl.Execute(w, res)
	if err != nil {
		data, err := json.Marshal(struct {
			Expr *models.Expression `json:"expression"`
		}{
			Expr: res,
		})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Write(data)
	}
}
