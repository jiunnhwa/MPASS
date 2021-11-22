/*

This package provides simulated calls and data.

*/

package data

import (
	"bytes"
	"fmt"
	"mpass/api/job"
	"mpass/model/crm"
	"strings"
	"text/template"
)

//Order details
type Order struct {
	Name        string
	OrderID     string
	DeliverTime string
	DeliverAddr string
	TrackingURL string
}

//ExpandBody expandsand finalizes the body text from the template into the actual body
func ExpandBody(t *job.TaskInfo) string {
	if strings.Contains(t.Type, "oun") {
		return GetOrderByID(t.To, t.Body)
	}
	if strings.Contains(t.Type, "mm") {
		return GetCustInfoByID(t.To, t.Body)
	}
	return t.Body
}

//GetOrderByID search an order by TelNum, and returns the details
func GetOrderByID(ID, msg string) string {

	std1 := Order{"Vaness", "duyff55", "24 jul 10am", "1 orchard road", "https://myorder.io/duyff55"}

	tmp1 := template.New("Template_1")
	tmp1, _ = tmp1.Parse(msg)

	var outBytes bytes.Buffer
	err := tmp1.Execute(&outBytes, std1)

	if err != nil {
		fmt.Println(err)
	}
	return outBytes.String()
}

//GetCustInfoByID search an order by TelNum, and returns the details
func GetCustInfoByID(ID, msg string) string {

	cust1 := crm.CustomerInfo{Name: "Charlie", Phone: "+6598219019", Email: "charlie@gmail.com", HomeAddrress: "1 orchard road"}
	tmp1 := template.New("Template_1")
	tmp1, _ = tmp1.Parse(msg)

	var outBytes bytes.Buffer
	err := tmp1.Execute(&outBytes, cust1)

	if err != nil {
		fmt.Println(err)
	}
	return outBytes.String()
}
