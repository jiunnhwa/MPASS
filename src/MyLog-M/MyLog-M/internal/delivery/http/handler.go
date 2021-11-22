package http

import (
	"MyLog-M/internal/domain"
	"MyLog-M/pkg/response"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:generate mockgen -source=./handler.go -destination=./mock/mock.go -package=mock
type service interface {
	Get(id int64) (*domain.Data, error)
	Insert(data domain.Data) (int64, error)
	Tail(limit int64) (*[]domain.Data, error)
}

type Handler struct {
	service service
}

//New Handler
func New(service service) *Handler {
	return &Handler{
		service: service,
	}
}

//MyLog gets a record by id, or insert a posted record
func (h Handler) MyLog(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ID := -1
		if n, err := strconv.Atoi(strings.ToUpper(r.URL.Query().Get("id"))); err == nil {
			ID = n
		}
		if ID < 0 {
			response.AsJSONError(w, http.StatusMethodNotAllowed, "Invalid id")
			return
		}
		resp, err := h.service.Get(int64(ID))
		log.Println(resp)
		if err != nil {
			response.AsJSONError(w, http.StatusMethodNotAllowed, err.Error())
			return
		}
		response.AsJSON(w, 200, resp)
		return
	}
	if r.Method == http.MethodPost {
		var data domain.Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			response.AsJSONError(w, 200, err.Error())
			return
		}
		LastID, err := h.service.Insert(data)
		if err != nil {
			response.AsJSONError(w, 200, err.Error())
			return
		}
		resp, err := h.service.Get(LastID)
		if err != nil {
			response.AsJSONError(w, 200, err.Error())
			return
		}
		log.Println(resp)
		response.AsJSON(w, 200, resp)
		return
	}
	response.AsJSONError(w, http.StatusMethodNotAllowed, "Invalid action")
}

//MyTail returns as json the last limit number of records, with default limit=1
func (h *Handler) MyTail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		limit := 1 //default show lastest 1
		if n, err := strconv.Atoi(strings.ToUpper(r.URL.Query().Get("limit"))); err == nil {
			limit = n
			log.Println("limit:", limit)
		}
		resp, err := h.service.Tail(int64(limit))
		if err != nil {
			response.AsJSONError(w, 200, err.Error())
			return
		}
		log.Println(resp)
		response.AsJSON(w, 200, resp)
		return
	}
	response.AsJSONError(w, http.StatusMethodNotAllowed, "Invalid action")
}

//MyView shows the last limit number of records, with default limit=1
func (h Handler) MyView(w http.ResponseWriter, r *http.Request) {
	tmpl := LoadTemplate(tplDir, "view.gohtml")
	if r.Method == http.MethodGet {
		limit := 1 //default show lastest 1
		if n, err := strconv.Atoi(strings.ToUpper(r.URL.Query().Get("limit"))); err == nil {
			limit = n
		}
		recs, err := h.service.Tail(int64(limit))
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		viewData := &ViewData{PageTitle: "VIEW", Records: *recs, RowCount: len(*recs)}
		tmpl.Execute(w, viewData)
		return

	}
	response.AsJSONError(w, http.StatusMethodNotAllowed, "Invalid action")
}

var tplDir string = "./html/templates"

//LoadTemplate loads the template tmplName from tplDir
func LoadTemplate(tplDir, tmplName string) *template.Template {
	t, err := template.New(tmplName).ParseFiles(tplDir + "/" + tmplName)
	if err != nil {
		panic(err)
	}
	return t
}

//ViewData is a collection of data for the view
type ViewData struct {
	PageTitle string
	Records   []domain.Data
	RowCount  int
}

// HealthCheck returns the status 200 OK
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `Ok`)
}

//handles home page
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to MyLog-As-A-Service.")
}
