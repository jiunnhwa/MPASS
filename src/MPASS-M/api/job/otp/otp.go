/*

	Package implements One-Time-Password(OTP).


	OTP is to be pre-filled by calling system.
	OTP implements instant, no-dnc calls.


	Sample:

*/

package otp

import (
	"database/sql"
	"io"
	"mpass/api/job"
	"mpass/filter/bloom"
	"mpass/provider"
	"strings"
	"time"
)

//OTP job type
type OTP struct {
	*job.Job
}

//NewJob constructs the view with parsing the files
func NewOTP(db *sql.DB) *OTP {
	j := &OTP{Job: job.NewJob(db)}
	return j
}

//Insert a new OTP job
func (j *OTP) Insert(rdr io.Reader) *OTP {
	if j.Parse(rdr).HasErr {
		return j
	}

	//Set defaults,
	now := time.Now().Local()
	j.SetMsgDefaults("otp/sms", now, now, "READY")

	//Check basic errors - empty fields, validity
	if err := j.HasCommonErrorsV1(); err != nil {
		j.SetResultErr(err.Error())
		return j
	}

	//Check specific errors
	if len(j.Recipients) != 1 {
		j.SetResultErr("Recipient can only be one")
		return j
	}

	//Check Rules
	if strings.ToUpper(strings.TrimSpace(j.OriginSystem)) != "FINANCE" {
		j.SetResultErr("Invalid OriginSystem")
		return j
	}

	//Check Bad Content
	if bloom.HasBadWord(j.Message) {
		j.SetResultErr("Has bad content")
		return j
	}

	//Assign Provider, default to twilio if no provider for OTP
	provider.AssignProviderOTP(j.Job)

	//Run
	j.RunSQLInsert()
	return j
}
