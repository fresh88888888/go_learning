package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const row int = 100
const col int = 100

func fillMatrix(m *[row][col]int) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			m[i][j] = s.Intn(10000)
		}
	}
}

func caculate(m *[row][col]int) {
	for i := 0; i < row; i++ {
		temp := 0
		for j := 0; j < col; j++ {
			temp += m[i][j]
		}
	}
}

/*
 *  multiple profile of go
 *
 *  https://golang.org//src/runtime/pprof/pprof.go
 */

func main() {
	// create output of file
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("Could not create CPU profile...", err)
	}

	// get sys info
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Could not start CPU profile...", err)
	}

	defer pprof.StopCPUProfile()

	// main logic section
	x := [row][col]int{}
	fillMatrix(&x)
	caculate(&x)

	f1, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("Could not creat memory profile...", err)
	}

	runtime.GC()

	if err := pprof.WriteHeapProfile(f1); err != nil {
		log.Fatal("Could not write memory profile...", err)
	}

	f1.Close()

	f2, err := os.Create("goroutine.prof")
	if err != nil {
		log.Fatal("Could not create groutine profile...", err)
	}

	if grof := pprof.Lookup("goroutine"); grof == nil {
		log.Fatal("Could not write groutine profile...")
	} else {
		grof.WriteTo(f2, 0)
	}

	f2.Close()
}
