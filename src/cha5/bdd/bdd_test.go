package bdd

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given 2 even numbers", t, func() {
		a := 2
		b := 4

		Convey("When add the two numbers", func() {
			c := a + b

			Convey("Then the result is still even", func() {
				So(c%2, ShouldEqual, 0)
			})
		})
	})
}

func checkType(i interface{}) {
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		fmt.Println("float....")
	case reflect.Int, reflect.Int32, reflect.Int64:
		fmt.Println("int....")
	default:
		fmt.Println("unknown")
	}
}
func TestBasicType(t *testing.T) {
	var f float64 = 2
	checkType(f)
}

func TestTypeAndValue(t *testing.T) {
	var f float64 = 10
	t.Log(reflect.TypeOf(f), reflect.ValueOf(f))
	t.Log(reflect.ValueOf(f).Type())
}

type Employe struct {
	EmployeID string
	Name      string `format:"normal"`
	Age       int
}

func (e *Employe) UpdateAge(age int) {
	e.Age = age
}

type Customer struct {
	CookieID string
	Name     string
	Age      int
}

func TestInvokeByName(t *testing.T) {
	e := Employe{"100", "Mike", 23}
	t.Logf("Name: value(%[1]v), type(%[1]T)", reflect.ValueOf(e).FieldByName("Name"))
	if nameField, ok := reflect.TypeOf(e).FieldByName("Name"); !ok {
		t.Error("Failed to get 'Name' field.")
	} else {
		t.Log("Tag:format", nameField.Tag.Get("format"))
	}

	reflect.ValueOf(&e).MethodByName("UpdateAge").Call([]reflect.Value{reflect.ValueOf(1)})
	t.Log("update age :", e)
}

func TestDeepEqual(t *testing.T) {
	a := map[int]string{1: "one", 2: "two", 3: "three"}
	b := map[int]string{1: "one", 2: "two", 3: "three"}

	t.Log(reflect.DeepEqual(a, b))

	s1 := []int{1, 2, 3}
	s2 := []int{2, 1, 3}
	s3 := []int{3, 1, 2}

	t.Log("s1 == s2?", reflect.DeepEqual(s1, s2))
	t.Log("s1 == s3?", reflect.DeepEqual(s1, s3))
}

func fillBySettings(st interface{}, settings map[string]interface{}) error {
	if reflect.TypeOf(st).Kind() != reflect.Ptr {
		if reflect.TypeOf(st).Elem().Kind() != reflect.Struct {
			return errors.New("The first parameter should be a point to struct type ")
		}
	}

	if settings == nil {
		return errors.New("settings is nil.")
	}
	var (
		field reflect.StructField
		ok    bool
	)

	for k, v := range settings {
		if field, ok = reflect.ValueOf(st).Elem().Type().FieldByName(k); !ok {
			continue
		}
		if field.Type == reflect.TypeOf(v) {
			vstr := reflect.ValueOf(st).Elem()
			vstr.FieldByName(k).Set(reflect.ValueOf(v))
		}
	}

	return nil
}

func TestFillNameAndAge(t *testing.T) {
	settings := map[string]interface{}{"Name": "Mike", "Age": 32}
	e := Employe{}
	if err := fillBySettings(&e, settings); err != nil {
		t.Fatal(err)
	}

	t.Log(e)
	c := new(Customer)
	if err := fillBySettings(c, settings); err != nil {
		t.Fatal(err)
	}
	t.Log(*c)
}

func TestUnsafe(t *testing.T) {
	i := 10
	f := *(*float64)(unsafe.Pointer(&i))
	t.Log(unsafe.Pointer(&i))
	t.Log(f)
}

type my_int int

func TestConvert(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := *(*[]my_int)(unsafe.Pointer(&a))
	t.Log(b)
}

func TestAtomic(t *testing.T) {
	var shardBufferPtr unsafe.Pointer
	writeDatafn := func() {
		data := []int{}
		for i := 0; i < 100; i++ {
			data = append(data, i)
		}
		atomic.StorePointer(&shardBufferPtr, unsafe.Pointer(&data))
	}

	readDataFn := func() {
		data := atomic.LoadPointer(&shardBufferPtr)
		fmt.Println(data, *(*[]int)(data))
	}

	var wg sync.WaitGroup
	writeDatafn()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				writeDatafn()
				time.Sleep(time.Microsecond * 100)
			}
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				readDataFn()
				time.Sleep(time.Microsecond * 100)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

type BasicInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type JobInfo struct {
	Skills []string `json:"skills"`
}

type Employee struct {
	BasicInfo BasicInfo `json:"basic_info"`
	JobInfo   JobInfo   `json:"job_info"`
}

var jsonStr = `{
	"basic_info":{
		"name":"Mike",
		"age": 20
	},
	"job_info": {
		"skills": ["Java", "c++", "Go"]
	}
}`

func TestEmbeddedJson(t *testing.T) {
	e := new(Employee)
	err := json.Unmarshal([]byte(jsonStr), e)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(*e)
	if v, err := json.Marshal(e); err == nil {
		fmt.Println(string(v))
	} else {
		t.Error(err)
	}
}
