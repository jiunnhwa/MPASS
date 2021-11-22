package domain

import "time"

type Data struct {
	RID         int64
	UnixTime    int
	LocalTime   time.Time //for view
	LogType     string
	LogSeverity int
	LogText     string

	Status Status
}
