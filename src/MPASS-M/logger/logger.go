package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//Record is the struct required for json post
type Record struct {
	RID         int64
	UnixTime    int
	LogType     string
	LogSeverity int
	LogText     string
}

var client *http.Client = &http.Client{}

const URL = "http://localhost:558/api/log"

//Log provides logging similar to a Syslog format
func Log(LogType string, LogSeverity int, LogText string) {
	rec := &Record{LogText: LogText}
	rec.LogType = LogType
	rec.LogSeverity = LogSeverity
	rec.LogText = LogText
	bytes, err := json.Marshal(rec)
	if err != nil {
		//fallback logging options, default to filelog
		log.Println(err)
		go Post(err.Error())
		return
	}
	//multi-write, can post to multiple destination if needed for redundancy
	log.Printf("%#v %#v %#v\n", LogType, LogSeverity, LogText)
	go Post(string(bytes))
}

func Post(json string) ([]byte, error) {
	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
	return bytes, nil
}
