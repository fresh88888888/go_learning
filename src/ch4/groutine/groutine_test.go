package groutine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			t.Log(i)
		}(i)
	}
	time.Sleep(time.Microsecond * 100)
}

func TestThreadSafe(t *testing.T) {
	count := 0
	var lock sync.Mutex
	for i := 0; i < 5000; i++ {
		go func() {
			defer func() {
				lock.Unlock()
			}()

			lock.Lock()
			count++
		}()
	}

	time.Sleep(1 * time.Second)
	t.Logf("count: %d", count)
}

func TestThreadWaitGroup(t *testing.T) {
	count := 0
	var lock sync.Mutex
	var vg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		vg.Add(1)
		go func() {
			defer func() {
				lock.Unlock()
			}()

			lock.Lock()
			count++
			vg.Done()
		}()
	}
	vg.Wait()
	t.Logf("count: %d", count)
}

func service() string {
	time.Sleep(time.Microsecond * 50)
	return "Done..."
}

func otherTask() {
	fmt.Println("working on something else.")
	time.Sleep(time.Microsecond * 100)
	fmt.Println("Task is done...")
}

func TestService(t *testing.T) {
	fmt.Println(service())
	otherTask()
}

func asyncService() chan string {
	rech := make(chan string, 1)
	go func() {
		ret := service()
		fmt.Println("returned result...")
		rech <- ret
		fmt.Println("service exit...")
	}()

	return rech
}

func TestAsyncService(t *testing.T) {
	rech := asyncService()
	otherTask()
	fmt.Println(<-rech)
	//time.Sleep(10 * time.Microsecond)
}

// func TestSelect(t *testing.T) {
// 	select {
// 	case ret := <-asyncService():
// 		t.Log(ret)
// 	case <-time.After(time.Microsecond * 200):
// 		t.Error("timeout...")
// 	}
// }

func dataProducer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()
}

func dataReceiver(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for {
			if data, ok := <-ch; ok {
				fmt.Println(data)
			} else {
				break
			}
		}
		wg.Done()
	}()
}

func TestCloseChannel(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	dataProducer(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Add(1)
	dataReceiver(ch, &wg)
	wg.Wait()
}

func isCancelled(cancelChan chan struct{}) bool {
	select {
	case <-cancelChan:
		return true
	default:
		return false
	}
}

func cancel_1(cancelChan chan struct{}) {
	cancelChan <- struct{}{}
}

func cancel_2(cancelChan chan struct{}) {
	close(cancelChan)
}

func TestCancel(t *testing.T) {
	cancelChan := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func(i int, ch chan struct{}) {
			for {
				if isCancelled(ch) {
					break
				}
				time.Sleep(time.Microsecond * 5)
			}
		}(i, cancelChan)
		fmt.Println(i, "canceled...")
	}
	cancel_2(cancelChan)
	time.Sleep(time.Second * 1)
}

func isContextCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 5; i++ {
		go func(i int, ctx context.Context) {
			for {
				if isContextCancelled(ctx) {
					break
				}
				time.Sleep(time.Microsecond * 5)
			}
		}(i, ctx)
		fmt.Println(i, "canceled...")
	}

	cancel()
	time.Sleep(time.Second * 1)
}

type Singleton struct {
}

var obj *Singleton
var once sync.Once

func getSingleton() *Singleton {
	once.Do(func() {
		fmt.Println("create singleton obj...")
		obj = new(Singleton)
	})

	return obj
}

func TestSingleton(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("i: %d, singleton obj: %p\n", i, getSingleton())
			wg.Done()
		}(i)
	}
	wg.Wait()
	// time.Sleep(time.Microsecond * 100)
}

func runTask(id int) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("The result is from %d", id)
}

func firstResponse() string {
	numOfRunner := 12
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}
	return <-ch
}

func allResponse() string {
	numOfRunner := 12
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func(i int) {
			ret := runTask(i)
			ch <- ret
		}(i)
	}
	final_resp := ""
	for i := 0; i < numOfRunner; i++ {
		final_resp += <-ch + "\n"
	}
	return final_resp
}

func TestFirstResponse(t *testing.T) {
	t.Log("before: ", runtime.NumGoroutine())
	t.Log(firstResponse())
	time.Sleep(time.Second * 1)
	t.Log("after: ", runtime.NumGoroutine())
}

func TestAllResponse(t *testing.T) {
	t.Log("before: ", runtime.NumGoroutine())
	t.Log(allResponse())
	time.Sleep(time.Second * 1)
	t.Log("after: ", runtime.NumGoroutine())
}

type ReusableObj struct {
}

type ObjPool struct {
	bufChan chan *ReusableObj //buffer channel for reuse obj
}

func newObjPool(numOfObj int) *ObjPool {
	obj_pool := ObjPool{}
	obj_pool.bufChan = make(chan *ReusableObj, numOfObj)
	for i := 0; i < numOfObj; i++ {
		obj_pool.bufChan <- &ReusableObj{}
	}

	return &obj_pool
}

func (p *ObjPool) getReuseableObj(timeout time.Duration) (*ReusableObj, error) {
	select {
	case ret := <-p.bufChan:
		return ret, nil
	case <-time.After(timeout): // control time
		return nil, errors.New("timeout...")
	}
}

func (p *ObjPool) releaseReuseableObj(obj *ReusableObj) error {
	select {
	case p.bufChan <- obj:
		return nil
	default:
		return errors.New("overflow...")
	}
}

func TestObjPool(t *testing.T) {
	pool := newObjPool(10)
	for i := 0; i < 11; i++ {
		if v, err := pool.getReuseableObj(time.Second * 1); err != nil {
			t.Error(err)
		} else {
			fmt.Printf("obj: %T\n", v)
			if err := pool.releaseReuseableObj(v); err != nil {
				t.Error(err)
			} else {
				t.Logf("release obj: %d success!", i)
			}
		}
	}
	fmt.Println("done...")
}

func TestSyncPool(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new pool...")
			return 100
		},
	}

	v := pool.Get().(int)
	fmt.Println(v)
	pool.Put(3)
	runtime.GC() // gc not clear sync pool.
	v1, _ := pool.Get().(int)
	fmt.Println(v1)
}

func TestSyncPoolmultiGroutine(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Create a new pool...")
			return 10
		},
	}

	pool.Put(100)
	pool.Put(100)
	pool.Put(100)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			fmt.Println(pool.Get())
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// func TestErrorInCode(t *testing.T) {
// 	fmt.Println("start...")
// 	t.Error("error...")
// 	fmt.Println("end...")
// }

// func TestFailInCode(t *testing.T) {
// 	fmt.Println("start...")
// 	t.Fatal("fail...")
// 	fmt.Println("end...")
// }

func TestAssert(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestCancatStringByAdd(t *testing.T) {
	assert := assert.New(t)
	elems := []string{"1", "2", "3", "4", "5"}
	ret := ""
	for _, elem := range elems {
		ret += elem
	}

	assert.Equal("12345", ret)
}

func TestConcatStringByByteBuffer(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	elems := []string{"1", "2", "3", "4", "5"}
	for _, elem := range elems {
		buf.WriteString(elem)
	}
	assert.Equal("12345", buf.String())
}

func BenchmarkCancatStringByAdd(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret := ""
		for _, elem := range elems {
			ret += elem
		}
	}

	b.StopTimer()
}

func BenchmarkConcatStringByByteBuffer(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for _, elem := range elems {
			buf.WriteString(elem)
		}
	}

	b.StopTimer()
}
