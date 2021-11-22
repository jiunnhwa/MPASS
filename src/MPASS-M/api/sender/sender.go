/*

	Package implements the base methods for sender.

*/

package sender

import (
	"database/sql"
	"fmt"
	"mpass/api/job"
	"mpass/logger"
	"strings"
	"time"
)

//====================== Feed ==========================
type Info struct {
	StartTime, EndTime time.Time

	URL          string
	ResponseText string
	ResponseCode int
}

type Sender struct {
	*job.Job

	Info
}

//NewJob constructs the view with parsing the files
func NewSender(db *sql.DB) *Sender {
	s := &Sender{Job: job.NewJob(db)}
	return s
}

//LockJobs locks the job based on Status and SendTime
func (s *Sender) LockJobs(workerID string, limit int) *Sender {
	newStatus, oldStatus := "LOCKED", "READY"
	sql := "UPDATE pegasus.message "
	sql += "SET WorkerName = '" + workerID + "', Status = '" + newStatus + "'  " //LOCKED
	sql += "WHERE Status = '" + oldStatus + "' "                                 //READY
	sql += "AND NOW() >= SendTime "
	sql += "ORDER BY RID Desc "
	if limit > 0 {
		sql += fmt.Sprintf("LIMIT %d ", limit)
	}
	sql += "; "
	s.ExecSQL(sql)
	return s
}

//GetSendingJob returns a list of jobs where status is 'LOCKED' by this worker, and has reached SendTime
func (s *Sender) GetSendingJob(workerID, status string, limit int) (*[]job.TaskInfo, error) {
	//func GetSendingJob(workerID, status string, limit int) (*[]TaskInfo, error) {
	//fmt.Println(time.Now(), "GetJob(", limit, ")")
	sql := "SELECT RID, Type, Providers, `From`, `To`, `Body`, SendTime, Status FROM pegasus.message " //need to specify fields
	sql += "WHERE Status = '" + status + "' " + "AND WorkerName = '" + workerID + "' "
	sql += " AND NOW() >= SendTime "

	if limit > 0 {
		sql += fmt.Sprintf("LIMIT %d ", limit)
	}

	//fmt.Println(sql)

	rows, err := s.Db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []job.TaskInfo
	for rows.Next() {
		item := job.TaskInfo{}
		if err := rows.Scan(&item.RID, &item.Type, &item.Providers, &item.From, &item.To, &item.Body, &item.SendTime, &item.Status); err != nil {
			//sql: Scan error on column index 4, name "Body": converting NULL to string is unsupported
			fmt.Println(err)
			return nil, err
		}
		result = append(result, item)
	}
	return &result, nil
}

//UpdateJobStatus
func (s *Sender) UpdateJobStatus(rid int64, info *Info, responses string) *Sender {
	//if len(strings.TrimSpace(info.ResponseText)) > 0 {
	if len(strings.TrimSpace(responses)) > 0 {
		Status := "END"
		if info.ResponseCode == 400 {
			Status = "ERROR"
		}
		// if responses == "DNC" {
		// 	Status = "DNC"
		// }
		sql := "UPDATE pegasus.message "
		sql += "SET Status = '" + Status + "', "
		sql += "StartTime = '" + info.StartTime.Format("2006-01-02 15:04:05") + "', " + "EndTime = '" + info.EndTime.Format("2006-01-02 15:04:05") + "', "
		sql += "Result = '" + StringEscape(responses) + "' "
		//sql += "Result = '" + StringEscape(info.ResponseText) + "' "
		//sql += "WHERE Status = '" + oldStatus + "' " + " AND RID = " + fmt.Sprint(rid) + " "
		sql += "WHERE RID = " + fmt.Sprint(rid) + " "
		sql += "; "
		s.ExecSQL(sql)

	}

	return s
}

func (s *Sender) UpdateJobStatus1(rid int64, info *Info, responses string) error {
	if len(strings.TrimSpace(responses)) > 0 {
		Status := "END"
		if info.ResponseCode == 400 {
			Status = "ERROR"
		}
		res, err := s.Db.Exec(`
			UPDATE pegasus.message 
			SET Status = ?, StartTime = ?, Result = ?
			WHERE RID = ?
		`, Status, info.StartTime.Format("2006-01-02 15:04:05"), responses, rid)
		if err != nil {
			return err
		}
		fmt.Print(res)
	}
	return nil
}

//UpdateSendStatus
func (s *Sender) UpdateSendStatus(rid int64, status string) *Sender {
	sql := "UPDATE pegasus.message "
	sql += "SET Status = '" + status + "' "
	sql += "WHERE RID = " + fmt.Sprint(rid) + " "
	sql += "; "
	s.ExecSQL(sql)
	return s
}

//*************************************************************
// SQL Sanitizer
//*************************************************************
func StringEscape(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	return str
}

//ExecSQL prepares and executes
func (s *Sender) ExecSQL(sql string) (int64, error) {
	stmt, err := s.Db.Prepare(sql)
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
	}
	res, err := stmt.Exec()
	if err != nil {
		//Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'WHERE ID='3333'' at line 1
		logger.Log("ERROR", 1, err.Error())
	}
	defer stmt.Close()
	lastId, err := res.LastInsertId()
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
	}
	return lastId, nil
}
