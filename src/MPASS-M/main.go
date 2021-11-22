package main

import (
	"fmt"
	"mpass/dnc"
	"mpass/logger"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var SID, AuthToken string //Twilio

//init, runs initialisation procedures
func init() {
	//Load configs
	SID, AuthToken = os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTHTOKEN") //Twilio
	fmt.Println("SID/TOKEN:", SID, AuthToken)
}

//main entry point
func main() {
	logger.Log("INFO", 0, "MPASS start ...")
	OpenDB(Connstr)
	defer CloseDB()

	// dnc.InsertDNC("+001", "Expired", time.Now().Add(time.Hour*-1), DB)
	// dnc.InsertDNC("+991", "Expire1", time.Now().Add(time.Minute), DB)
	// dnc.InsertDNC("+995", "Expire5", time.Now().Add(time.Minute*5), DB)
	// dnc.InsertDNC("+000", "OPTOUT", time.Time{}, DB)

	go dnc.AutoDeleteExpiredDNC(DB)

	go ServeRoutes()             //API listener
	go OnTicker(time.Second, "") //Run jobs

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
	defer logger.Log("INFO", 0, "MPASS end ...")
}

//Ticker, runs funcs on tick.
func OnTicker(duration time.Duration, desc string) {
	ticker := time.NewTicker(duration)
	for ; true; <-ticker.C {
		JobManager()
	}
}
