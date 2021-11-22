/*

	Package implements the sending to Twilio.

*/

package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mpass/api/sender"
	"mpass/logger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type myResponse struct {
	*http.Response
	Bytes      []byte
	Status     string
	StatusCode int
	ErrCode    int
}

type Twilio struct {
	accountSid string
	authToken  string
	URL        string
	RID        int64
	Info       sender.Info

	Bytes      []byte
	Status     string
	StatusCode int
	ErrCode    int

	Logs []string
}

//NewView constructs the view with parsing the files
func New(rid int64, accountSid, authToken string) *Twilio {
	tw := &Twilio{}
	tw.RID = rid
	tw.accountSid = accountSid
	tw.authToken = authToken
	fmt.Println("twilio:", tw.accountSid, tw.authToken)
	tw.URL = "https://api.twilio.com/2010-04-01/Accounts/" + tw.accountSid + "/Messages.json"
	return tw
}

const (
	AUTHORIZED_FROM_SMS      = "+18606070103"
	AUTHORIZED_FROM_WHATSAPP = "+14155238886"
)

//SendWhatsApp post the whatsapp message
func (tw *Twilio) SendWhatsApp(from, to, body string) *Twilio {
	fmt.Println("SendWhatsApp", from, to, body)
	tw.Info.StartTime = time.Now()

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", "whatsapp:"+to) //NUMBER_TO
	//msgData.Set("From", from) //NUMBER_FROM
	msgData.Set("From", "whatsapp:"+"+14155238886")           //NUMBER_FROM override, with authorised number
	msgData.Set("From", "whatsapp:"+AUTHORIZED_FROM_WHATSAPP) //NUMBER_FROM override, with authorised number
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", tw.URL, &msgDataReader)
	req.SetBasicAuth(tw.accountSid, tw.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	tw.Info.ResponseCode = resp.StatusCode
	tw.Info.ResponseText = string(*RespReadBytes(resp))
	fmt.Println("Response:", resp.StatusCode, "\n", tw.Info.ResponseText)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)

		if err == nil {
			//fmt.Println(data["sid"])
			tw.Logs = append(tw.Logs, fmt.Sprint(data["sid"]))

		}
	} else {
		//fmt.Println("twilio", resp.Status)
		tw.Logs = append(tw.Logs, fmt.Sprint("Response Status:", resp.Status))
	}
	tw.Info.EndTime = time.Now()
	return tw
}

func RespReadBytes(resp *http.Response) *[]byte {
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logger.Log("ERROR", 1, readErr.Error())
	}
	return &body
}

//

//TrySendSMS post the sms to provider
//with Retries, client timeout, exponential backoff
func (tw *Twilio) TrySendSMS(from, to, body string, MaxTries int) *Twilio {
	timeout := time.Duration(10) * time.Second
	backoff := GetNextNum(0)
	for i := 0; i < MaxTries; i++ {
		tw.Post(AUTHORIZED_FROM_SMS, to, body, timeout)
		if tw.StatusCode == 200 || tw.StatusCode == 201 || tw.StatusCode == 202 {
			break
		}
		//handle generic non-retry error
		if tw.StatusCode == 404 {
			logger.Log("DEBUG", 0, fmt.Sprintf("Status: %s\tBody : %s\n", tw.Status, tw.Bytes))
			return tw
		}
		if tw.ErrCode > 0 {
			logger.Log("DEBUG", 0, fmt.Sprintf("Status: %s\tBody : %s\n", tw.Status, tw.Bytes))
			//handle any specific non-retry error

			//do retry
			if tw.ErrCode == 1000 && strings.Contains(tw.Status, "Timeout") {
				logger.Log("DEBUG", 0, fmt.Sprintf("Status: %s\tBody : %s\n", tw.Status, tw.Bytes))
				timeout = timeout * 2
				time.Sleep(time.Duration(backoff()+RandMinMax(0, 3)) * time.Second)
				continue
			}
		}
	}
	logger.Log("DEBUG", 0, fmt.Sprintf("Status: %s\tBody : %s\n", tw.Status, tw.Bytes))
	return tw
}

func (tw *Twilio) Post(from, to, body string, timeout time.Duration) *Twilio {
	// Prepare
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	c := http.Client{Timeout: timeout}
	req, err := http.NewRequest("POST", tw.URL, &msgDataReader)
	if err != nil {
		tw.ErrCode, tw.Status = 1000, err.Error()
		return tw
	}
	req.SetBasicAuth(tw.accountSid, tw.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//Post
	resp, err := c.Do(req)
	if err != nil {
		tw.Logs = append(tw.Logs, err.Error())
		tw.ErrCode, tw.Status = 1000, err.Error()
		return tw
	}
	defer resp.Body.Close()
	tw.Status = resp.Status
	tw.StatusCode = resp.StatusCode
	bytes, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		tw.Logs = append(tw.Logs, err.Error())
		tw.ErrCode, tw.Status = 1000, errRead.Error()
		return tw
	}
	tw.Bytes = bytes

	//Process
	tw.Info.ResponseCode = resp.StatusCode
	tw.Info.ResponseText = string(*RespReadBytes(resp))
	tw.Logs = append(tw.Logs, fmt.Sprint("StatusCode:", resp.StatusCode, "ResponseText:", tw.Info.ResponseText))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			//fmt.Println(data["sid"])
			tw.Logs = append(tw.Logs, fmt.Sprint(data["sid"]))
		}
	} else {
		//fmt.Println("twilio", resp.Status)
		tw.Logs = append(tw.Logs, fmt.Sprint("Response Status:", resp.Status))
	}
	tw.Info.EndTime = time.Now()
	return tw
}

//AppendLog adds text to the log object
func (tw *Twilio) AppendLog(logText string) *Twilio {
	tw.Logs = append(tw.Logs, logText)
	return tw
}

//GetNextNum uses a closure to generate sequences, doubling for every call
func GetNextNum(startNum int) func() int {
	currNum := startNum
	return func() int {
		currNum += 1
		return currNum * 2
	}
}

func RandMinMax(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(max-min+1) + min)
}
