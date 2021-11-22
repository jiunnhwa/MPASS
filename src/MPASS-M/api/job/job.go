/*

	Package implements the base methods for a job/message type.

*/

package job

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mpass/logger"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	StatusCode int
	Text       string
}

type Job struct {
	Db  *sql.DB
	RID int64

	From         string //sender phonenum
	Recipients   []string
	Provider     string //vendor
	OriginSystem string
	Message      string
	MessageType  string //sms,email,whatsapp,...
	CreateTime   time.Time
	CreateBy     string

	SendTime   time.Time //the time to send the message
	SendStatus string    //for otp, set tor ready, others leave empty

	HasErr    bool
	ErrorCode int
	ErrorDesc string

	Result     interface{}
	ResultText string

	Logs []string

	Response
}

//NewView constructs the view with parsing the files
// func NewJob(recipients []string) *Job {
// 	j := &Job{recipients, "Marketing", "Hi {{.Name}}, your order {{.OrderID}} has been scheduled to arrive at {{.DeliverAddr}} \n Track:{{.TrackingURL}}", "sms", time.Now(), "Super"}
// 	return j
// }

//NewJob constructs the view with parsing the files
func NewJob(db *sql.DB) *Job {
	j := &Job{Db: db, CreateTime: time.Now().Local()} //createTime := time.Now().Format("2006-01-02 15:04:05"), if not formatted, it is UTC
	return j
}

//func GetJobByID(rid int, db *sql.DB) (*TaskInfo, error) {
func (j *Job) GetJobByID(r *http.Request) *Job {
	JobID := -1
	if n, err := strconv.Atoi(strings.ToUpper(r.URL.Query().Get("jobid"))); err == nil {
		JobID = n
		log.Println("JobID:", JobID)

	}
	if JobID < 0 {
		j.SetResultErr("Invalid JobID")
		return j
	}

	sql := "SELECT `RID`, `Type`, `Providers`, `From`, `To`, `Body`, `CreateTime`, `CreateBy`, `OriginSystem`, `SendTime`, `WorkerName`, `StartTime`, `EndTime`, `Status`, `Result` FROM pegasus.message "
	sql += "WHERE RID = " + fmt.Sprint(JobID) + " "
	fmt.Println(sql)

	rows, err := j.Db.Query(sql)
	if err != nil {
		j.SetResultErr(err.Error())
		return j
	}
	defer rows.Close()

	item := TaskInfo{}
	for rows.Next() {

		//sql: Scan called without calling Next

		//if err := rows.Scan(&item.RID, &item.Providers, &item.From, &item.To, &item.Body); err != nil {
		if err := rows.Scan(&item.RID, &item.Type, &item.Providers, &item.From, &item.To, &item.Body, &item.CreateTime, &item.CreateBy, &item.OriginSystem, &item.SendTime, &item.WorkerName, &item.StartTime, &item.EndTime, &item.Status, &item.Result); err != nil {
			//sql: Scan error on column index 4, name "Body": converting NULL to string is unsupported
			j.SetResultErr(err.Error())
			return j
		}
	}
	if item.RID == 0 {
		j.SetResultErr(fmt.Sprintf("No records found for JobID %d", JobID))
		return j
	}
	j.Response.StatusCode = 200
	j.Result = item
	response, _ := json.Marshal(item)
	j.ResultText = string(response)
	return j
}

//Parse the json into the job struct, and sets error
func (j *Job) Parse(rdr io.Reader) *Job {
	err := json.NewDecoder(rdr).Decode(&j)
	log.Println("Parse:", err)
	log.Printf("%#v\n", j)
	if err != nil {
		j.HasErr = true
		j.SetResultErr("Bad json")
	}
	return j
}

//HasCommonErrors checks the basic required fields are not empty, including valid phone number
func (j *Job) HasCommonErrors() bool {
	//Check empty fields
	if j.IsInvalid((len(strings.TrimSpace(j.From)) == 0), "From cannot be empty") {
		return true
	}

	if j.IsInvalid((len(j.Recipients) == 0), "Recipients cannot be empty") {
		return true
	}

	if j.IsInvalid((len(strings.TrimSpace(j.Message)) == 0), "Message cannot be empty") {
		return true
	}

	if j.IsInvalid((len(strings.TrimSpace(j.CreateBy)) == 0), "CreateBy cannot be empty") {
		return true
	}

	if j.IsInvalid((len(strings.TrimSpace(j.OriginSystem)) == 0), "OriginSystem cannot be empty") {
		return true
	}

	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	for k, v := range j.Recipients {
		if j.IsInvalid(!re.MatchString(v), "Invalid phone number for item "+fmt.Sprint(k)) {
			return true
		}
	}
	return false
}

const phoneNumberRegex string = `^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`

func (j *Job) HasCommonErrorsV1() error {
	if strings.TrimSpace(j.From) == "" {
		return errors.New("From cannot be empty")
	}

	if len(j.Recipients) == 0 {
		return errors.New("Recipients cannot be empty")
	}

	if strings.TrimSpace(j.Message) == "" {
		return errors.New("Message cannot be empty")
	}

	if strings.TrimSpace(j.CreateBy) == "" {
		return errors.New("CreateBy cannot be empty")
	}

	if strings.TrimSpace(j.OriginSystem) == "" {
		return errors.New("OriginSystem cannot be empty")
	}

	re := regexp.MustCompile(phoneNumberRegex)
	for k, v := range j.Recipients {
		if !re.MatchString(v) {
			return fmt.Errorf("Invalid phone number for item %s", fmt.Sprint(k))
		}
	}
	return nil
}

//SetDefaults set the defaults per the message type
func (j *Job) SetMsgDefaults(messageType string, createTime, sendTime time.Time, sendStatus string) *Job {
	//override if not empty
	if len(messageType) > 0 {
		j.MessageType = messageType //"otp"
	}

	j.CreateTime = createTime //time.Now().Local()
	j.SendTime = sendTime     //j.CreateTime

	//override if not empty
	if len(sendStatus) > 0 {
		j.SendStatus = sendStatus // "READY"
	}

	return j
}

//Insert a new job
func (j *Job) Insert(rdr io.Reader) *Job {
	json.NewDecoder(rdr).Decode(&j)

	//Check
	if len(j.Recipients) > 0 {
		for k, v := range j.Recipients {
			//TODO: validate as phone num
			if len(strings.TrimSpace(v)) == 0 {
				j.SetResultErr(fmt.Sprintf("Recipient at index %d is empty", k))
				return j
			}
		}
	} else {
		j.SetResultErr("Recipients cannot be empty")
		return j
	}

	if !j.IsFieldValid(len(strings.TrimSpace(j.Message)) > 0, "Message cannot be empty") {
		return j
	}

	if !j.IsFieldValid(len(strings.TrimSpace(j.MessageType)) > 0, "MessageType cannot be empty") {
		return j
	}

	if len(strings.TrimSpace(j.CreateBy)) == 0 {
		j.SetResultErr("CreateBy cannot be empty")
		return j
	}

	j.RunSQLInsert()
	return j

}

func (j *Job) RunSQLInsert() *Job {

	//Exec SQL
	stmt, err := j.Db.Prepare("INSERT INTO pegasus.message (`Type`, `Providers`, `From`, `To`,  `Body`, `CreateBy`, `CreateTime`, `SendTime`, `Status`, `OriginSystem`) VALUES (?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	//createTimeLocal := time.Now().Format("2006-01-02 15:04:05") //over-ride the NewJob.CreateTime which is UTC

	res, err := stmt.Exec(j.MessageType, j.Provider, j.From, strings.Join(j.Recipients, ","), j.Message, j.CreateBy /*j.CreateTime*/, j.SetDefaultCreateTime(), j.SetDefaultSendTime(), j.SetDefaultStatus("NEW"), j.OriginSystem)
	if err != nil {
		j.SetResultErr(err.Error())
		return j
	}
	rid, err := res.LastInsertId()
	if err != nil {
		j.SetResultErr(err.Error())
		return j
	}
	j.RID = (rid)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		j.SetResultErr(err.Error())
		return j
	}

	//Log and return result
	j.SetResultOK(rid, rowsAffected)
	log.Printf("Job: %#v\n", j)
	logger.Log("INFO", 0, "LastInsertId(pegasus.message):"+fmt.Sprint(rid))
	return j
}

func (j *Job) SetDefaultCreateTime() string {
	createTimeLocal := time.Now().Format("2006-01-02 15:04:05")
	if !j.CreateTime.IsZero() {
		createTimeLocal = j.CreateTime.Format("2006-01-02 15:04:05")
	}
	return createTimeLocal
}

func (j *Job) SetDefaultSendTime() string {
	sendTimeLocal := "0-0-0 00:00:00"
	if !j.SendTime.IsZero() {
		sendTimeLocal = j.SendTime.Format("2006-01-02 15:04:05")
	}
	return sendTimeLocal
}

func (j *Job) SetDefaultStatus(status string) string {
	if len(strings.TrimSpace(j.SendStatus)) > 0 {
		status = j.SendStatus
	}
	return status
}

//SetResultErr
func (j *Job) SetResultErr(errorDesc string) *Job {
	j.HasErr = true
	j.ErrorDesc = errorDesc
	j.Response.StatusCode = http.StatusBadRequest
	obj := struct {
		Error string
	}{
		Error: j.ErrorDesc,
	}
	j.Result = obj
	response, _ := json.Marshal(obj)
	j.ResultText = (string(response))
	log.Printf("%s\n", j.ErrorDesc)
	logger.Log("ERROR", 1, j.ErrorDesc)
	return j
}

//SetResultOK
func (j *Job) SetResultOK(rid, rowsAffected int64) *Job {
	j.Response.StatusCode = 200
	obj := struct {
		RID             int64
		RecordsAffected interface{}
	}{
		RID:             rid,
		RecordsAffected: rowsAffected,
	}
	j.Result = obj
	response, _ := json.Marshal(obj)
	j.ResultText = string(response)
	return j
}

func (j *Job) IsFieldValid(ok bool, errDesc string) bool {
	if !ok {
		j.SetResultErr(errDesc)
		return false
	}
	return true
}

func (j *Job) IsInvalid(invalid bool, errDesc string) bool {
	if invalid {
		j.SetResultErr(errDesc)
		return true
	}
	return false
}

//TaskInfo is job struct interface to db, due to null values(eg. sql.NullTime) considerations
type TaskInfo struct {
	RID int64

	Type      string //# sms,email
	Providers string // ordered preferences of twilio,starhub

	From          string
	To            string
	Body          string
	CreateTime    sql.NullTime
	CreateTimeStr string
	CreateBy      string

	OriginSystem sql.NullString

	SendTime sql.NullTime // #time to send

	WorkerName sql.NullString
	StartTime  sql.NullTime //#actual send time
	EndTime    sql.NullTime // 	#done time

	Status sql.NullString //  #new,approved,validated,queued,locked,finished
	Result sql.NullString
}
