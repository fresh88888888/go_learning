package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

// func main() {

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello World!")
// 	})

// 	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
// 		t := time.Now()
// 		time_str := fmt.Sprintf("{\"time\":\"%s\"}", t)
// 		w.Write([]byte(time_str))
// 	})

//		http.ListenAndServe(":8080", nil)
//	}
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

func CreateFBS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var fbs []int
	var err error
	for i := 0; i < 100000; i++ {
		fbs, err = getFibonacci(50)
		if err != nil {
			break
		}
	}
	w.Write([]byte(fmt.Sprintf("%v", fbs)))
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintln(w, " welcome...")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, " hello, %s\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)
	router.GET("/fb", CreateFBS)

	log.Fatal(http.ListenAndServe(":8081", router))
}
