package main

import (
	"fmt"
	"log"
	"mpass/api/job"
	"mpass/api/job/mm"
	"mpass/api/job/otp"
	"mpass/api/job/oun"
	html "mpass/util/html"
	"mpass/util/response"
	"net/http"
	"strings"
	"time"
)

const (
	IP   = ""
	PORT = "88"
)

//ServeRoutes handles the api endpoints
func ServeRoutes() error {

	//VIEWS
	http.HandleFunc("/", home)

	//API
	http.HandleFunc("/api/job", myJob)
	http.HandleFunc("/api/otp", myOTP)
	http.HandleFunc("/api/oun", myOUN)
	http.HandleFunc("/api/mm", myMM)

	//REPORTS
	http.HandleFunc("/reports/daily", myDaily)

	//fileServer :=
	http.Handle("/report/", http.StripPrefix("/report", http.FileServer(http.Dir("./data/report/"))))

	//path := "certs\\"
	//log.Fatal(http.ListenAndServeTLS(":8080", path+"ssl.cert", path+"ssl.key", nil))

	fmt.Println("ServiceBroker listening at", IP, ":", PORT)
	if err := http.ListenAndServe(IP+":"+PORT, nil); err != nil {
		return err
	}

	return nil
}

//handles home page, updates the view data and serve
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Messaging-Platform-As-A-Service.")
}

//ViewData is a collection of data for the view
type ViewData struct {
	PageTitle  string
	ReportDate string
	Records    []job.TaskInfo
	RowCount   int
}

var tplDir string = "./html/templates"

//myView provides a report view of the list of transactions by parameter date
func myDaily(w http.ResponseWriter, r *http.Request) {
	tmpl := html.LoadTemplate(tplDir, "view.gohtml")
	if r.Method == http.MethodGet {
		recs := *Get(r)
		reportDate := r.URL.Query().Get("date")
		if reportDate == "" {
			reportDate = time.Now().Local().AddDate(0, 0, -1).Format("2006-01-02")
		}
		viewData := &ViewData{PageTitle: "View - Daily Transactions", ReportDate: reportDate, Records: recs, RowCount: len(recs)}
		tmpl.Execute(w, viewData)
		return
	}
	response.AsJSONError(w, http.StatusMethodNotAllowed, "Invalid action")
}

func Get(r *http.Request) *[]job.TaskInfo {
	var result []job.TaskInfo
	reportDate, err := time.Parse("2006-01-02", strings.ToLower(r.URL.Query().Get("date")))
	if err != nil {
		reportDate = time.Now().Local().AddDate(0, 0, -1)

	}
	log.Println("RID:", reportDate)
	sql := "SELECT RID, Providers, `From`, `To`, `Body`, SendTime, Status FROM pegasus.message WHERE Status = 'END' AND date(SendTime) = '" + reportDate.Format("2006-01-02") + "' "
	sql += "ORDER BY RID DESC "
	// if rid > 0 {
	// 	sql += "WHERE RID = " + fmt.Sprint(rid) + " "
	// }
	fmt.Println(sql)

	rows, err := DB.Query(sql)
	if err != nil {
		item := job.TaskInfo{}
		//item.Status.Code, item.Status.Text = -100, err.Error()
		result = append(result, item)
		return &result
	}
	defer rows.Close()
	for rows.Next() {
		item := job.TaskInfo{}
		if err := rows.Scan(&item.RID, &item.Providers, &item.From, &item.To, &item.Body, &item.SendTime, &item.Status); err != nil {
			fmt.Println(err)
			return &result
		}
		//fmt.Println((item))
		result = append(result, item)
	}
	return &result
}

//myJob is the base implentation for Read and Insert a new job
func myJob(w http.ResponseWriter, r *http.Request) {
	j := job.NewJob(DB)
	if r.Method == http.MethodGet {
		response.AsJSON(w, j.StatusCode, j.GetJobByID(r).Result)
		return
	}
	if r.Method == http.MethodPost {
		response.AsJSON(w, j.StatusCode, j.Insert(r.Body).Result)
		return
	}
	response.AsJSONError(w, http.StatusBadRequest, "Invalid action")
}

//myOTP provides read and insert of an OTP message type
func myOTP(w http.ResponseWriter, r *http.Request) {
	j := otp.NewOTP(DB)
	if r.Method == http.MethodGet {
		response.AsJSON(w, j.StatusCode, j.GetJobByID(r).Result)
		return
	}
	if r.Method == http.MethodPost {
		response.AsJSON(w, j.StatusCode, j.Insert(r.Body).Result)
		return
	}
	response.AsJSONError(w, http.StatusBadRequest, "Invalid action")
}

//myOUN provides read and insert of an OUN message type
func myOUN(w http.ResponseWriter, r *http.Request) {
	j := oun.NewOUN(DB)
	if r.Method == http.MethodGet {
		response.AsJSON(w, j.StatusCode, j.GetJobByID(r).Result)
		return
	}
	if r.Method == http.MethodPost {
		response.AsJSON(w, j.StatusCode, j.Insert(r.Body).Result)
		return
	}
	response.AsJSONError(w, http.StatusBadRequest, "Invalid action")
}

//myMM provides read and insert of a MM message type
func myMM(w http.ResponseWriter, r *http.Request) {
	j := mm.NewMM(DB)
	if r.Method == http.MethodGet {
		response.AsJSON(w, j.StatusCode, j.GetJobByID(r).Result)
		return
	}
	if r.Method == http.MethodPost {
		response.AsJSON(w, j.StatusCode, j.Insert(r.Body).Result)
		return
	}
	response.AsJSONError(w, http.StatusBadRequest, "Invalid action")
}
