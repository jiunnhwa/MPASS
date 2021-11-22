/*

	Package implements Order Update Notification(OUN).


	OUN is to be pre-filled by calling system.
	OUN implements deferred, no-dnc calls.
	On sending, OUN will expand the template message, and filled it with order details (name, order id, deladdr,turl), from CRM for instance.

	Sample:
	Promotions are valid from 17 June- 14 July 2021. • <Beauty Buffet> Promotion is not valid for new items, promotional items, and promotional banded packs.

	[Lazada] Your 5% discount is expiring soon. shop now at https://kkkk.

*/

package oun

import (
	"database/sql"
	"io"
	"mpass/api/job"
	"mpass/filter/bloom"
	"mpass/provider"
	"time"
)

type OUN struct {
	*job.Job
}

//NewOUN creates a new OUN job
func NewOUN(db *sql.DB) *OUN {
	j := &OUN{Job: job.NewJob(db)}
	return j
}

//Insert a new OUN job
func (j *OUN) Insert(rdr io.Reader) *OUN {
	if j.Parse(rdr).HasErr {
		return j
	}

	//Set defaults, messageType is set from origin, sendTime can be a future time(default to 1 min for demo purposes)
	now, sendTime := time.Now().Local(), time.Now().Local().Add(time.Minute)
	j.SetMsgDefaults("", now, sendTime, "READY")

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
	// -- NONE YET --

	//Check Bad Content
	if bloom.HasBadWord(j.Message) {
		j.SetResultErr("Has bad content")
		return j
	}

	//Assign Provider, default to twilio if no provider for sms
	provider.AssignProvider(j.Job)

	//Run
	j.RunSQLInsert()
	return j
}
