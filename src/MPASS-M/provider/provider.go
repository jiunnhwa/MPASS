/*

	Package implements the base methods for providers.

*/

package provider

import (
	"mpass/api/job"
	"strings"
)

//AssignProvider sets the gateway
//AssignProvider defaults to Twilio if no provider for sms,
//AssignProvider will also set the provider based on recipient number at this stage for single recipient, otherwise late bind upon sending
func AssignProvider(j *job.Job) {
	if strings.ToUpper(strings.TrimSpace(j.MessageType)) == "SMS" {
		prov := strings.TrimSpace(j.Provider)
		if len(prov) == 0 {
			j.Provider = "twilio"
			return
		}

		if len(j.Recipients) == 1 {
			j.Provider = GetTelco(j.Provider, j.Recipients[0])
		}

	}
}

//AssignProviderOTP defaults to Twilio if no provider specified for OTP
func AssignProviderOTP(j *job.Job) {
	if len(strings.TrimSpace(j.Provider)) == 0 {
		ss := strings.Split(strings.ToUpper(strings.TrimSpace(j.MessageType)), "/") //eg: otp/sms
		if ss[0] == "OTP" {
			j.Provider = "twilio"
		}
	}
}

//GetTelco returns the telco based on the recipient phone number
func GetTelco(telco, phonenum string) string {
	if strings.HasPrefix(phonenum, "+65123") {
		return "TelcoA"
	}
	if strings.HasPrefix(phonenum, "+65456") {
		return "TelcoB"
	}
	if len(telco) == 0 {
		return "twilio"
	}
	return telco
}
