package response

import (
	"encoding/json"
	"log"
	"net/http"
)

//AsJSON writes out the header and body for a json payload.
func AsJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	log.Println("AsJSON:", code, string(response), err)
	if err != nil {
		AsJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if response != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
		return
	}
	AsJSONError(w, http.StatusBadRequest, "Bad json")
}

//AsJSON writes out the header and body for a json payload.
func AsJSONError(w http.ResponseWriter, code int, errorDesc string) {
	obj := struct {
		Error string
	}{
		Error: errorDesc,
	}
	response, _ := json.Marshal(obj)
	log.Println(string(response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
