package gcfriendly

import "testing"

const num_of_elements = 1000

type Content struct {
	Detail [10000]int
}

func withValue(arr [num_of_elements]Content) int {
	return 0
}

func withReference(arr *[num_of_elements]Content) int {
	return 0
}

func TestFn(t *testing.T) {
	var arr [num_of_elements]Content
	withValue(arr)
	withReference(&arr)
}

func BenchmarkPassingArrayWithValue(b *testing.B) {
	var arr [num_of_elements]Content
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		withValue(arr)
	}
	b.StopTimer()
}

func BenchmarkPassingArrayWithReference(b *testing.B) {
	var arr [num_of_elements]Content
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		withReference(&arr)
	}
	b.StopTimer()
}
