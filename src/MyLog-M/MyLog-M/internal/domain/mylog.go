package domain

import "time"

type Status struct {
	Code int
	Text string
}

type Record struct {
	RID         int64
	UnixTime    int
	LocalTime   time.Time //for view
	LogType     string
	LogSeverity int
	LogText     string

	Status Status
}
