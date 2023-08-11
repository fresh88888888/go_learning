package service

import "errors"

var lessThanTwoError = errors.New("n is not less than 2")
var largeThanThousandError = errors.New("n is not large than 1000")

func GetFibonacci(n int) ([]int, error) {
	if n < 2 {
		return nil, lessThanTwoError
	}
	if n > 1000 {
		return nil, largeThanThousandError
	}

	fib_list := []int{1, 1}
	for i := 2; i < n; i++ {
		fib_list = append(fib_list, fib_list[i-2]+fib_list[i-1])
	}

	return fib_list, nil
}

func Square(n int) int {
	return n * n
}
