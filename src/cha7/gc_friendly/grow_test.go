package gcfriendly

import "testing"

const times = 1000

func TestAutoGrow(t *testing.T) {
	for i := 0; i < times; i++ {
		s := []int{}
		for j := 0; j < num_of_elements; j++ {
			s = append(s, j)
		}
	}
}

func TestProperInit(t *testing.T) {
	for i := 0; i < times; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < num_of_elements; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkAutoGrow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := []int{}
		for j := 0; j < num_of_elements; j++ {
			s = append(s, j)
		}
	}
	b.StopTimer()
}

func BenchmarkProperInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, num_of_elements)
		for j := 0; j < num_of_elements; j++ {
			s = append(s, j)
		}
	}
	b.StopTimer()
}

func BenchmarkOverSizeInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, num_of_elements*8)
		for j := 0; j < num_of_elements; j++ {
			s = append(s, j)
		}
	}
	b.StopTimer()
}
