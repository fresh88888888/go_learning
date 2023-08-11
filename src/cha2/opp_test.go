package opp_test

import (
	"errors"
	"fmt"
	"testing"
)

func doSomething(f interface{}) {
	switch v := f.(type) {
	case int:
		fmt.Println("Integer ", v)
	case string:
		fmt.Println("string ", v)
	default:
		fmt.Println("Unknown")
	}
}

func TestEmptyInterfaceAssertion(t *testing.T) {
	doSomething(10)
	doSomething("123")
}

var lessThanTwoError = errors.New("n is not less than 2")
var largeThanThousandError = errors.New("n is not large than 1000")

func getFibonacci(n int) ([]int, error) {
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

func TestGetFibonacci(t *testing.T) {
	if l, err := getFibonacci(1); err == nil {
		t.Log(l)
	} else {
		switch err {
		case lessThanTwoError:
			t.Log("lessThanTwoError: " + err.Error())
		case largeThanThousandError:
			t.Log("largeThanThousandError: " + err.Error())
		}

	}
}

func TestRecover(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover from...", err)
		}
	}()
	fmt.Println("start...")
	panic(errors.New("something wrong!"))
}
