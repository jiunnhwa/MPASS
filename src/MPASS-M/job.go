package main

import (
	"mpass/api/job"
	"mpass/api/sender"
	"mpass/data"
	"mpass/dnc"
	"mpass/logger"
	"mpass/provider/twilio"
	"strings"
	"time"
)

const (
	WorkerID      = "W1"
	JobCount      = 1
	ClientRetries = 3
)

//JobManager handles concurrent job sending
func JobManager() {
	s := sender.NewSender(DB)
	jobs, _ := s.LockJobs(WorkerID, JobCount).GetSendingJob(WorkerID, "LOCKED", JobCount)
	for _, v := range *jobs {
		go worker(&v, s) //async send
	}
}

//worker runs the task
func worker(j *job.TaskInfo, s *sender.Sender) {
	//handles more specific types first, each block uses the same standard steps

	if strings.Contains(j.Type, "otp") {
		s.UpdateSendStatus(j.RID, "SENDING") //set intermediate status to prevent re-entrant select.
		tw := twilio.New(j.RID, SID, AuthToken).TrySendSMS(j.From, j.To, data.ExpandBody(j), ClientRetries)
		res := "[" + strings.TrimSuffix(strings.TrimSpace(tw.Info.ResponseText), ",") + "]"
		tw.AppendLog(res)
		logger.Log("INFO", 0, strings.Join(tw.Logs, "\n"))
		s.UpdateJobStatus(j.RID, &tw.Info, res)
		return
	}

	if strings.Contains(j.Type, "mm") {
		//Do
		s.UpdateSendStatus(j.RID, "SENDING") //set intermediate status to prevent re-entrant select.
		recipients := strings.Split(j.To, ",")
		responses := "" //  provider response, manual concat

		tw := twilio.New(j.RID, SID, AuthToken)
		for _, recipientTel := range recipients {
			//Check for DNC
			if dnc.IsDNC(recipientTel, s.Job.Db) {
				responses += `{ "Result": "DNC" }` + ","
			} else {
				//TrySend
				tw.TrySendSMS(j.From, recipientTel, data.ExpandBody(j), ClientRetries)
				responses += tw.Info.ResponseText + ","

				//Set DNC +1 minute
				dnc.InsertDNC(recipientTel, "JustSent", time.Now().Add(time.Minute), s.Job.Db)
			}
		}
		res := "[" + strings.TrimSuffix(strings.TrimSpace(responses), ",") + "]"
		logger.Log("INFO", 0, res)
		s.UpdateJobStatus(j.RID, &tw.Info, res)

		return
	}

	if strings.Contains(j.Type, "whatsapp") {
		s.UpdateSendStatus(j.RID, "SENDING") //set intermediate status to prevent re-entrant select.
		recipients := strings.Split(j.To, ",")
		responses := "" //  provider response, manual concat

		tw := twilio.New(j.RID, SID, AuthToken)
		for _, recipientTel := range recipients {
			tw.SendWhatsApp(j.From, recipientTel, data.ExpandBody(j))
			responses += tw.Info.ResponseText + ","
		}
		res := "[" + strings.TrimSuffix(strings.TrimSpace(responses), ",") + "]"

		logger.Log("INFO", 0, res)
		s.UpdateJobStatus(j.RID, &tw.Info, res)
		return
	}

	//if specified, or catch all/sms
	if j.Providers == "twilio" || (len(strings.TrimSpace(j.Providers)) == 0 && strings.Contains(j.Type, "sms")) {
		s.UpdateSendStatus(j.RID, "SENDING") //set intermediate status to prevent re-entrant select.
		recipients := strings.Split(j.To, ",")
		responses := "" //  provider response, manual concat

		tw := twilio.New(j.RID, SID, AuthToken)
		for _, recipientTel := range recipients {
			tw.TrySendSMS(j.From, recipientTel, data.ExpandBody(j), ClientRetries)
			responses += tw.Info.ResponseText + ","

		}
		res := "[" + strings.TrimSuffix(strings.TrimSpace(responses), ",") + "]"

		logger.Log("INFO", 0, res)
		s.UpdateJobStatus(j.RID, &tw.Info, res)
		return
	}

}
