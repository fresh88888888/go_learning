package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

var jsonStr = `{
	"basic_info":{
		"name":"Mike",
		"age": 20
	},
	"job_info": {
		"skills": ["Java", "c++", "Go"]
	}
}`

func TestEasyJson(t *testing.T) {
	e := &Employee{}
	e.UnmarshalJSON([]byte(jsonStr))
	fmt.Println(e)
	if v, err := e.MarshalJSON(); err != nil {
		t.Error(err)
	} else {
		fmt.Println(string(v))
	}
}

func BenchmarkEmbeddedjson(b *testing.B) {
	b.ResetTimer()
	e := new(Employee)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsonStr), e)
		if err != nil {
			b.Error(err)
		}
		if _, err := json.Marshal(e); err != nil {
			b.Error(err)
		}
	}

	b.StartTimer()
}

func BenchmarkEasyjson(b *testing.B) {
	b.ResetTimer()
	e := &Employee{}
	for i := 0; i < b.N; i++ {
		e.UnmarshalJSON([]byte(jsonStr))
		if _, err := e.MarshalJSON(); err != nil {
			b.Error(err)
		}
	}

	b.StopTimer()
}
