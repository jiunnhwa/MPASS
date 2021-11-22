/*

	Package implements the sending to TelcoA.

*/

package telcoA

import (
	"fmt"
	"mpass/api/sender"

	"time"
)

type TelcoA struct {
	accountSid string
	authToken  string
	URL        string
	RID        int64
	Info       sender.Info
}

//NewView constructs the view with parsing the files
func New(rid int64) *TelcoA {
	tc := &TelcoA{}
	tc.RID = rid
	return tc
}

//SendSMS post the sms to provider
func (tc *TelcoA) SendSMS(from, to, body string) *TelcoA {
	fmt.Println("SendSMS", from, to, body)
	tc.Info.StartTime = time.Now()

	//Mock Send
	tc.Info.ResponseCode = 200
	tc.Info.ResponseText = "sucess"
	fmt.Println("Response:", tc.Info.ResponseCode, "\n", tc.Info.ResponseText)

	time.Sleep(time.Second) //simulate work
	tc.Info.EndTime = time.Now()
	return tc
}
