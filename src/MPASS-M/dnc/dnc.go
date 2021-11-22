package dnc

import (
	"database/sql"
	"fmt"
	"mpass/logger"
	"time"
)

//IsDNC returns if a phone/email is in active DNC
func IsDNC(key string, db *sql.DB) bool {
	var isDNC bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM pegasus.DNC WHERE RID = ?", key).Scan(&isDNC)
	if err != nil {
		return true
	}
	return isDNC
}

//InsertDNC inserts a new active DNC phone/email
func InsertDNC(id, reasonCode string, expireTime time.Time, db *sql.DB) (int64, error) {
	sql := fmt.Sprintf("INSERT INTO pegasus.DNC (`RID`, `ReasonCode`, `ExpireTime` ) VALUES (?,?,?);")
	res, err := db.Exec(sql, id, reasonCode, expireTime.Format("2006-01-02 15:04:05"))
	logger.Log("INFO", 0, fmt.Sprintf("InsertDNC. Values = %s, %s, %s", id, reasonCode, expireTime.Format("2006-01-02 15:04:05")))
	if err != nil {
		return -1, err
	}
	rid, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return rid, nil
}

//DeleteExpiredDNC delete records that have expired
func DeleteExpiredDNC(db *sql.DB) int64 {
	sql := fmt.Sprintf("DELETE FROM pegasus.DNC WHERE NOW()>ExpireTime AND YEAR(ExpireTime)>1 ; ")
	result, err := db.Exec(sql)
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
		return -1
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
		return -1
	}
	if rowsAffected > 0 {
		logger.Log("INFO", 0, fmt.Sprintf("DeleteExpiredDNC. Rows affected = %d", rowsAffected))
	}
	return rowsAffected
}

//AutoDeleteExpiredSessions auto delete expired sessions every second
func AutoDeleteExpiredDNC(db *sql.DB) {
	for {
		DeleteExpiredDNC(db)
		time.Sleep(time.Second)
	}
}
