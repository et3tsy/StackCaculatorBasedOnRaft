// use 'go test -run Tsy'
// to run all tests

package calculator

import (
	"fmt"
	"testing"
)

func TestCreateAndDeleteByTsy(t *testing.T) {
	c, arr := CreateInit()
	err := c.Delete(arr[0])
	if err != nil {
		t.Error(err)
	}
	err = c.Delete(arr[0])
	if err == nil {
		t.Error("delete fail")
	}
}

func TestPushAndPopByTsy(t *testing.T) {
	c, arr := CreateInit()
	if err := c.Push(arr[0], 1); err != nil {
		t.Error(err)
	}

	v, err := c.Pop(arr[0])
	if err != nil {
		t.Errorf("%v", err)
	}
	if v != 1 {
		t.Error("pop error")
	}

	_, err = c.Pop(arr[0])
	if err == nil {
		t.Error("pop error")
	}

	c.Push(arr[2], 2)
	c.Push(arr[1], 3)
	c.Push(arr[2], 4)
	v, err = c.Pop(arr[2])
	if err != nil {
		t.Error("pop error")
	}
	if v != 4 {
		t.Error("pop error")
	}
}

func TestGetByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 4)
	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 4 {
		t.Error("get error")
	}
}

func TestAddByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 4)
	err := c.Add(arr[0])
	if err != nil {
		t.Error(err)
	}
	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 6 {
		t.Error("add error")
	}
}

func TestSubByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 4)
	err := c.Sub(arr[0])
	if err != nil {
		t.Error(err)
	}
	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 2 {
		t.Errorf("sub error, expected: %v, got: %v", 2, v)
	}
}

func TestMulByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 14)
	err := c.Mul(arr[0])
	if err != nil {
		t.Error(err)
	}
	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 28 {
		t.Error("mul error")
	}
}

func TestDivByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 13)
	err := c.Div(arr[0])
	if err != nil {
		t.Error(err)
	}
	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 6 {
		t.Error("div error")
	}
}

func TestIncByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 13)
	err := c.Inc(arr[0])
	if err != nil {
		t.Error(err)
	}

	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 14 {
		t.Error("div error")
	}
}

func TestDecByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	c.Push(arr[0], 13)
	err := c.Dec(arr[0])
	if err != nil {
		t.Error(err)
	}

	v, err := c.Get(arr[0])
	if err != nil {
		t.Error(err)
	}
	if v != 12 {
		t.Error("div error")
	}
}

func TestSthByTsy(t *testing.T) {
	c, arr := CreateInit()
	c.Push(arr[0], 2)
	fmt.Println(arr[0], arr[1])
	fmt.Println(c.Get(arr[0]))
	fmt.Println(c.Get(arr[1]))
}

func CreateInit() (c *Calculator, arr []int64) {
	c = Make(nil, nil)
	for i := 0; i < 3; i++ {
		v, _ := c.Create()
		arr = append(arr, v)
	}
	return
}
