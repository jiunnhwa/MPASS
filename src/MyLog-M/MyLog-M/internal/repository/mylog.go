package repository

import (
	"MyLog-M/internal/domain"
	"database/sql"
	"fmt"
	"time"
)

type LogRepo struct {
	db *sql.DB
}

func NewLogRepo(db *sql.DB) *LogRepo {
	return &LogRepo{db: db}
}

func (this *LogRepo) Get(id int64) (*domain.Data, error) {
	item := &domain.Data{}

	sql := "SELECT `RID`, `UnixTime`, `Type`, `Severity`, `Log` FROM Log "
	sql += "WHERE RID = " + fmt.Sprint(id) + " "
	fmt.Println(sql)

	rows, err := this.db.Query(sql)
	if err != nil {
		return item, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&item.RID, &item.UnixTime, &item.LogType, &item.LogSeverity, &item.LogText); err != nil {
			return item, err
		}
	}
	if item.RID == 0 {
		fmt.Println(item)
		return item, fmt.Errorf(fmt.Sprintf("No records found for ID %v", item.RID))
	}
	fmt.Println(item)
	return item, nil
}

func (this *LogRepo) Insert(data *domain.Data) (int64, error) {
	sql := "INSERT INTO Log(`UnixTime`,`Type`,`Severity`,`Log`) VALUES (strftime('%s','now'),?,?,?)"
	res, err := this.db.Exec(sql, data.LogType, data.LogSeverity, data.LogText)
	if err != nil {
		return -1, err
	}
	rid, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return rid, nil
}

//GetLogTail shows the last n records
func (this *LogRepo) Tail(limit int64) (*[]domain.Data, error) {
	var result []domain.Data
	sql := "SELECT `RID`, `UnixTime`, `Type`, `Severity`, `Log` FROM Log "
	sql += "ORDER BY RID DESC "
	if limit > 0 {
		sql += fmt.Sprintf("LIMIT %d ", limit)
	}
	fmt.Println(sql)

	rows, err := this.db.Query(sql)
	if err != nil {
		return &result, err
	}

	defer rows.Close()
	for rows.Next() {
		item := domain.Data{}
		if err := rows.Scan(&item.RID, &item.UnixTime, &item.LogType, &item.LogSeverity, &item.LogText); err != nil {
			fmt.Println(err)
			return &result, err
		}
		item.LocalTime = time.Unix(int64(item.UnixTime), 0)
		fmt.Println((item))
		result = append(result, item)
	}
	return &result, nil
}
