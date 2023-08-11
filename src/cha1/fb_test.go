package fb_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"
)

const (
	Monday  = 1
	Tuesday = 2
	Wednesday
)

const (
	Readable = 1 << iota
	Writable
	Executeable
)

func TestConstTry(t *testing.T) {
	t.Log(Monday, Tuesday)
	a := 1
	t.Log(a&Readable == Readable, a&Writable == Writable, a&Executeable == Executeable)
}

type my_int int64

func TestImplicit(t *testing.T) {
	var a int32 = 1
	var b int64
	var c my_int

	b = int64(a)
	c = my_int(b)

	t.Log(a, b, c)
}

func TestPoint(t *testing.T) {
	var a int = 3
	p := &a
	t.Log(a, p)
	t.Logf("%T, %T", a, p)
}

func TestString(t *testing.T) {
	var str string
	t.Log("*" + str + "*")
	t.Log(len(str))
}

func TestBitClear(t *testing.T) {
	a := 7 //0111
	a = a &^ Writable
	a = a &^ Executeable

	t.Log(a&Readable == Readable, a&Writable == Writable, a&Executeable == Executeable)
}

func TestWileLoop(t *testing.T) {
	n := 1
	for n < 6 {
		t.Log(n)
		n++
	}
}

func TestIf(t *testing.T) {
	if a := true; a {
		t.Log(a)
	}
}

func TestSwitchMultiCase(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch i {
		case 0, 2:
			t.Log("Evevn")
		case 1, 3:
			t.Log("Odd")
		default:
			t.Log("It is not 0~3!")
		}
	}
}

func TestSwitchCaseCondition(t *testing.T) {
	for i := 0; i < 5; i++ {
		switch {
		case i%2 == 0:
			t.Log("I is a even")
		case i%2 == 1:
			t.Log("I is a odd")
		default:
			t.Log("I is not in 0~5")
		}
	}
}

func TestArrayInit(t *testing.T) {
	var arr [3]int
	arr1 := [4]int{1, 2, 3, 4}
	arr2 := [...]int{1, 2, 3, 4, 5}
	arr1[1] = 5
	t.Log(arr[1], arr[2])
	t.Log(arr1, arr2)
}

func TestTravel(t *testing.T) {
	arr3 := [...]int{1, 2, 3, 4, 5}
	for i := 0; i < len(arr3); i++ {
		t.Log(arr3[i])
	}

	for _, e := range arr3 {
		t.Log(e)
	}
}

func TestArraySection(t *testing.T) {
	var arr []int = []int{1, 2, 3, 4, 5}
	arr3 := arr[:3]
	arr4 := arr[3:]
	arr5 := arr[:]

	t.Log(arr3, arr4, arr5)
}

func TestSliceInit(t *testing.T) {
	var s0 []int
	t.Log(len(s0), cap(s0))
	s0 = append(s0, 1)
	t.Log(len(s0), cap(s0))

	s1 := []int{1, 2, 3, 4, 5}
	t.Log(len(s1), cap(s1))

	s2 := make([]int, 3, 5)
	t.Log(len(s2), cap(s2))
	s2 = append(s2, 1)
	t.Log(s2[0], s2[1], s2[2], s2[3])
	t.Log(len(s2), cap(s2))
}

func TestSliceGrowing(t *testing.T) {
	var arr = []int{}
	for i := 0; i < 10; i++ {
		arr = append(arr, i)
		t.Log(len(arr), cap(arr))
	}
}

func TestSliceShareMemory(t *testing.T) {
	year := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "NoV", "Dec"}
	Q2 := year[3:6]
	t.Log(len(Q2), cap(Q2))
	summary := year[5:8]
	t.Log(len(summary), cap(summary))
	summary[0] = "Unkown"
	t.Log(Q2)
	t.Log(year)
}

func TestMapInit(t *testing.T) {
	map1 := map[int]int{1: 2, 2: 4, 4: 8, 3: 9}
	t.Log(map1[3])
	t.Logf("map1 len=%d", len(map1))

	map2 := map[int]int{}
	map2[1] = 2
	t.Logf("map2 len=%d", len(map2))

	map3 := make(map[int]int, 10)
	t.Logf("map3 len=%d", len(map3))
}

func TestAccessNoExistkey(t *testing.T) {
	m1 := map[int]int{}
	t.Log(m1[1])
	m1[2] = 0
	t.Log(m1[2])
	m1[3] = 1
	if v, ok := m1[3]; ok {
		t.Logf("key 3 is exist, v=%d", v)
	} else {
		t.Log("key 3 is not exist")
	}
}

func TestTravelMap(t *testing.T) {
	m1 := map[int]int{1: 2, 2: 4, 4: 8, 3: 9}
	for k, v := range m1 {
		t.Logf("key= %d, value= %d", k, v)
	}
}

func TestMapWithFunction(t *testing.T) {
	map1 := map[int]func(int) int{}
	map1[0] = func(ops int) int { return ops }
	map1[1] = func(ops int) int { return ops * ops }
	map1[2] = func(ops int) int { return ops * ops * ops }

	t.Log(map1[0](2), map1[1](3), map1[2](4))
}

func TestMapForSet(t *testing.T) {
	mySet := map[int]bool{}
	mySet[1] = true
	n := 1
	if mySet[n] {
		t.Logf("%d is existing!", n)
	} else {
		t.Logf("%d is not existing", n)
	}
	t.Logf("myset len=%d", len(mySet))
	delete(mySet, 1)
	if mySet[n] {
		t.Logf("%d is existing!", n)
	} else {
		t.Logf("%d is not existing", n)
	}
}

func TestStringInit(t *testing.T) {
	var s string
	t.Log(s)
	s = "hello"
	t.Log(len(s))
	// s[1] = '2'  string not mutiple byte slice
	s = "\xE4\xB8\xA5"
	t.Log(s)

	s = "串"
	t.Log(len(s))

	c := []rune(s)
	t.Log(len(c))
	t.Logf("串 unicode %x", c[0])
	t.Logf("串 utf8 %x", s)
}

func TestStringToRune(t *testing.T) {
	s := "条件分支"
	for _, v := range s {
		t.Logf("%[1]c, %[1]x", v)
	}
}

func TestStringSplitFunc(t *testing.T) {
	s := "A, B, C, D"
	parts := strings.Split(s, ",")
	for _, part := range parts {
		t.Log(strings.Trim(part, " "))
	}
	t.Log(strings.Join(parts, "-"))
}

func TestStringConv(t *testing.T) {
	s := strconv.Itoa(9)
	t.Logf("str = %s", s)
	if i, error := strconv.Atoi("10"); error == nil {
		t.Log(10 + i)
	}
}

func returnMultivalues() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

func TestMultiCase(t *testing.T) {
	a, _ := returnMultivalues()
	t.Log(a)
}

type IntConv func(ops int) int

func timeSpent(inner IntConv) IntConv {
	return func(ops int) int {
		start := time.Now()
		ret := inner(ops)
		fmt.Println("time spent: ", time.Since(start).Seconds())

		return ret
	}
}

func slowFunc(ops int) int {
	time.Sleep(time.Second * 1)
	return ops
}

func TestFunc(t *testing.T) {
	a, _ := returnMultivalues()
	t.Log(a)
	sf := timeSpent(slowFunc)
	t.Log(sf(10))
}

func sum(nums ...int) int {
	var sum int
	for _, n := range nums {
		sum += n
	}

	return sum
}

func TestVarParm(t *testing.T) {
	t.Log(sum(1, 2, 3, 4, 5))
	t.Log(sum(1, 2, 3, 4, 5, 6))
}

func clear() {
	fmt.Println("Clear resource!")
}
func TestDefer(t *testing.T) {
	defer clear()
	t.Log("start.....")
	// panic("err")
}

type Employee struct {
	Id   string
	Name string
	Age  int
}

func TestCreateEmployeeObj(t *testing.T) {
	e := Employee{"0", "Bob", 20}
	e1 := Employee{Name: "Mike", Age: 30}
	e2 := new(Employee)
	e2.Id = "1"
	e2.Name = "Jack"
	e2.Age = 28

	t.Log(e)
	t.Log(e1)
	t.Log(e1.Id)
	t.Log(e2)
	t.Logf("e is %T", e)
	t.Logf("e2 is %T", e2)
}

// func (e Employee) String() string {
// 	fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
// 	return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.Id, e.Name, e.Age)
// }

func (e *Employee) String() string {
	fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
	return fmt.Sprintf("ID:%s-Name:%s-Age:%d", e.Id, e.Name, e.Age)
}
func TestStructOperations(t *testing.T) {
	e := Employee{"1", "Bob", 20}
	fmt.Printf("Address is %x\n", unsafe.Pointer(&e.Name))
	t.Log(e.String())
}

type Code string
type Programer interface {
	writeHelloWorld() Code
}

type GoProgramer struct {
	Id string
}

func (*GoProgramer) writeHelloWorld() Code {
	return "fmt.Println(\"Hello World! \")"
}

type JavaProgramer struct {
	Id string
}

func (*JavaProgramer) writeHelloWorld() Code {
	return "System.out.Println(\"Hello World! \")"
}

func TestClient(t *testing.T) {
	var p Programer = new(GoProgramer)
	t.Log(p.writeHelloWorld())
}

type Pet struct {
}

func (p *Pet) speak() {
	fmt.Println("pet of speak func.")
}

func (p *Pet) speakTo(host string) {
	p.speak()
	fmt.Println("pet of speakTo func.")
}

type Dog struct {
	Pet
}

func (d *Dog) speakTo(host string) {
	d.speak()
	fmt.Println("dog of speakTo func.")
}

func TestDog(t *testing.T) {
	dog := new(Dog)
	dog.speakTo("chao")
}

func writeFirstProgramer(p Programer) {
	fmt.Printf("%T, %v\n", p, p.writeHelloWorld())
}

func TestPolymorphism(t *testing.T) {
	gor := new(GoProgramer)
	java := new(JavaProgramer)
	writeFirstProgramer(gor)
	writeFirstProgramer(java)
}
