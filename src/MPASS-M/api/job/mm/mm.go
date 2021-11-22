/*

	Package implements Marketing Messages (MM).

	MM is to be pre-filled by calling system.
	MM implements deferred, no-dnc calls.
	On sending, MM will expand the template message, and filled it with order details (name, order id, deladdr,turl), from CRM for instance.

	Sample:
	Promotions are valid from 17 June- 14 July 2021. • <Beauty Buffet> Promotion is not valid for new items, promotional items, and promotional banded packs.

	[Lazada] Your 5% discount is expiring soon. shop now at https://kkkk.


*/

package mm

import (
	"database/sql"
	"io"
	"mpass/api/job"
	"mpass/filter/bloom"
	"strings"
	"time"
)

type MM struct {
	*job.Job
}

//NewMM creates a new MM job
func NewMM(db *sql.DB) *MM {
	j := &MM{Job: job.NewJob(db)}
	return j
}

//Insert a new MM job
func (j *MM) Insert(rdr io.Reader) *MM {
	if j.Parse(rdr).HasErr {
		return j
	}

	//Set defaults, messageType is set from origin, sendTime can be a future time(default to 1 min for demo purposes)
	now, sendTime := time.Now().Local(), time.Now().Local().Add(time.Minute)
	j.SetMsgDefaults("", now, sendTime, "READY")

	//Check basic errors - empty fields, validity
	if j.HasCommonErrors() {
		return j
	}

	//Check specific errors

	//Check Rules
	if j.IsInvalid((strings.ToUpper(strings.TrimSpace(j.OriginSystem)) == "FINANCE"), "Invalid OriginSystem") {
		return j
	}

	//Check Bad Content
	if bloom.HasBadWord(j.Message) {
		j.SetResultErr("Has bad content")
		return j
	}

	//Assign Provider, default to twilio if no provider for sms
	//Late bind, sender will determine as it is multi-recipients/phone
	//provider.AssignProvider(j.Job)

	//Run
	j.RunSQLInsert()
	return j
}
