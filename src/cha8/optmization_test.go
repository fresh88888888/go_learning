package optmiz

import "testing"

func TestCreateRequest(t *testing.T) {
	str := createRequest()
	t.Log(str)
}

func TestProcessRequest(t *testing.T) {
	reqs := []string{}
	reqs = append(reqs, createRequest())
	t.Log(reqs)
	reps := processRequest(reqs)
	t.Log(reps)
}

func BenchmarkProcessRequestOld(b *testing.B) {
	reqs := []string{}
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequestOld(reqs)
	}
	b.StopTimer()
}

func BenchmarkProcessRequest(b *testing.B) {
	reqs := []string{}
	reqs = append(reqs, createRequest())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processRequest(reqs)
	}
	b.StopTimer()
}
