package optmiz

import (
	"encoding/json"
	"strconv"
	"strings"
)

func createRequest() string {
	payload := make([]int, 100, 100)
	for i := 0; i < 100; i++ {
		payload[i] = i
	}
	req := Request{"demo_transaction", payload}
	v, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}

	return string(v)
}

func processRequestOld(reqs []string) []string {
	reps := []string{}
	for _, req := range reqs {
		reqObj := &Request{}
		json.Unmarshal([]byte(req), reqObj)
		ret := ""
		for _, e := range reqObj.Payload {
			ret += strconv.Itoa(e) + ","
		}

		repObj := &Response{reqObj.TransactionID, ret}
		repJson, err := json.Marshal(&repObj)
		if err != nil {
			panic(err)
		}
		reps = append(reps, string(repJson))
	}
	return reps
}

func processRequest(reqs []string) []string {
	reps := []string{}
	for _, req := range reqs {
		reqObj := &Request{}
		reqObj.UnmarshalJSON([]byte(req))
		var buf strings.Builder
		for _, e := range reqObj.Payload {
			buf.WriteString(strconv.Itoa(e))
			buf.WriteString(",")
		}

		repObj := &Response{reqObj.TransactionID, buf.String()}
		repJson, err := repObj.MarshalJSON()
		if err != nil {
			panic(err)
		}
		reps = append(reps, string(repJson))
	}
	return reps
}
