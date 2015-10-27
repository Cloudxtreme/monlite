// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2010-09-21
// Last modification:	2010-

package types

import (
	"errors"
	"reflect"
	"testing"
)

type testitem struct {
	t    interface{}
	name string
	val  interface{}
}

type testvec0 []int
type testvec1 []testitem
type teststr string

var boolean bool = true
var ptrbool *bool = &boolean
var interger int = 42
var ptrint *int = &interger
var interger32 int32 = 42
var ptrint32 *int32 = &interger32
var interger8 uint8 = 42
var ptrint8 *uint8 = &interger8
var cadeia string = "outrastring"
var ptrcadeia *string = &cadeia

type channel chan bool

var ptrChannel *channel

var tests []testitem = []testitem{
	{true, "bool", true},
	{new(bool), "*bool", ptrbool},
	{int(0), "int", 42},
	{new(int), "*int", ptrint},
	{int32(0), "int32", int32(64)},
	{new(int32), "*int32", ptrint32},
	{byte(0), "uint8", uint8(1)},
	{new(byte), "*uint8", ptrint8},
	{"string", "string", "istoeumastring"},
	{new(string), "*string", ptrcadeia},
	{struct{}{}, "struct {}", struct{}{}},
	{&struct{}{}, "*struct {}", &struct{}{}},
	{testitem{}, "github.com/fcavani/types.testitem", testitem{int(0), "int", 0}},
	{&testitem{}, "*github.com/fcavani/types.testitem", &testitem{int(0), "int", 1}},
	{[3]int{}, "[3]int", [3]int{1, 2, 3}},
	{[]int{}, "[]int", []int{1, 2, 3}},
	{[2]testitem{}, "[2]github.com/fcavani/types.testitem", [2]testitem{{int(9), "int", 0}, {int(9), "int", 2}}},
	{[]testitem{}, "[]github.com/fcavani/types.testitem", []testitem{{int(7), "int", 0}, {int(8), "int", 2}}},
	{testvec0{}, "github.com/fcavani/types.testvec0", testvec0{1, 2, 3, 4}},
	{map[string]string{}, "map[string]string", map[string]string{"a": "1", "b": "2"}},
	{map[string]string{}, "map[string]string", map[string]string{}},
	//{&testvec0{}, "*serialization/types.testvec0", &testvec0{1,2,3,4}},
	{teststr("oi"), "github.com/fcavani/types.teststr", teststr("oi")},
	{make(channel), "github.com/fcavani/types.channel", make(channel)},
	{new(channel), "*github.com/fcavani/types.channel", new(channel)},
}

func TestNameOf(t *testing.T) {
	//print("\nTestNameOf\n\n")
	for i, test := range tests {
		s := NameOf(reflect.ValueOf(test.t).Type())
		if s != test.name {
			t.Fatalf("NameOf %v failed: %v != %v", i, test.name, s)
		}
	}
}

func TestIsert(t *testing.T) {
	for _, typ := range tests {
		Insert(typ.t)
	}
}

func TestMake(t *testing.T) {
	for i, typ := range tests {
		val := MakeNew(typ.name, 0)
		//println("can set:", val.CanSet())
		name := NameOf(val.Type())
		//println(name)
		if name != typ.name {
			t.Fatalf("type name differ in %v: %v != %v", i, name, typ.name)
		}
		//println(i)
		val.Set(reflect.ValueOf(typ.val))
		if !reflect.DeepEqual(val.Interface(), typ.val) {
			t.Fatalf("not equal: %v", i)
		}
	}
}

func TestCopyBool(t *testing.T) {
	b := true
	cp := Copy(reflect.ValueOf(b))
	b = false
	if cp.Bool() != true {
		t.Fatal("copy failed")
	}
}

func TestCopyInt(t *testing.T) {
	var i int = -42
	cp := Copy(reflect.ValueOf(i))
	i = 0
	if cp.Interface().(int) != -42 {
		t.Fatal("copy failed")
	}
	var i8 int8 = 1<<7 - 1
	cp = Copy(reflect.ValueOf(i8))
	i8 = 0
	if cp.Interface().(int8) != 1<<7-1 {
		t.Fatal("copy failed")
	}
	if cp.Kind() != reflect.Int8 {
		t.Fatal("wrong kind")
	}
}

func TestCopyUint(t *testing.T) {
	var ui uint = 42
	cp := Copy(reflect.ValueOf(ui))
	ui = 0
	if cp.Interface().(uint) != 42 {
		t.Fatal("copy failed")
	}
}

func TestCopyFloat32(t *testing.T) {
	var f float32 = 3.1415
	cp := Copy(reflect.ValueOf(f))
	f = 0.0
	if cp.Interface().(float32) != 3.1415 {
		t.Fatal("copy failed")
	}
}

func TestCopyComplex(t *testing.T) {
	var c complex64 = complex(1, 1)
	cp := Copy(reflect.ValueOf(c))
	c = 0.0
	if cp.Interface().(complex64) != complex(1, 1) {
		t.Fatal("copy failed")
	}
}

func TestCopyPtr(t *testing.T) {
	var ui uint = 42
	pui := &ui
	cp := Copy(reflect.ValueOf(pui))
	*pui = 0
	if *(cp.Interface().(*uint)) != 42 {
		t.Fatal("copy failed")
	}
}

func TestCopyPtrNil(t *testing.T) {
	var ui uint = 42
	var pui *uint
	cp := Copy(reflect.ValueOf(pui))
	pui = &ui
	if cp.Interface().(*uint) != nil {
		t.Fatal("copy failed")
	}
}

type IntForCopy interface {
	Foo(x string) string
}

type CopyTextStruct struct {
	Name  string
	Val   float32
	Ref   *CopyTextStruct
	dummy complex64
}

func (c *CopyTextStruct) Foo(x string) string {
	c.Name = x
	return x
}

func TestCopyStruct1(t *testing.T) {
	s := CopyTextStruct{
		Name: "foo",
		Val:  2.2,
		Ref:  nil,
	}
	cp := Copy(reflect.ValueOf(s))
	s.Name = "blá"
	n := cp.Interface().(CopyTextStruct)
	if n.Name != "foo" || n.Val != 2.2 || n.Ref != nil {
		t.Fatalf("copy failed: %#v %#v", cp.Interface(), s)
	}
}

func TestCopyStruct2(t *testing.T) {
	ps := &CopyTextStruct{
		Name: "foo",
		Val:  2.2,
		Ref:  nil,
	}
	ps.Ref = ps
	cp := Copy(reflect.ValueOf(ps))
	if !reflect.DeepEqual(cp.Interface(), ps) {
		t.Fatalf("copy failed: %#v", cp.Interface())
	}
}

func TestCopyArray(t *testing.T) {
	a := [3]int{1, 2, 3}
	cp := Copy(reflect.ValueOf(a))
	a[0] = 0
	acp := cp.Interface().([3]int)
	if acp[0] != 1 || acp[1] != 2 || acp[2] != 3 {
		t.Fatal("copy failed")
	}
}

func TestCopyMap(t *testing.T) {
	m := map[string]string{
		"foo":   "bar",
		"test1": "test2",
		"blá":   "blá",
	}
	cp := Copy(reflect.ValueOf(m))
	m["foo"] = "catoto"
	mcp := cp.Interface().(map[string]string)
	if mcp["foo"] != "bar" || mcp["test1"] != "test2" || mcp["blá"] != "blá" {
		t.Fatal("copy failed")
	}
}

func TestCopySlice(t *testing.T) {
	s := []int{1, 2, 3}
	cp := Copy(reflect.ValueOf(s))
	s[0] = 0
	scp := cp.Interface().([]int)
	if scp[0] != 1 || scp[1] != 2 || scp[2] != 3 {
		t.Fatal("copy failed")
	}
}

func TestCopyInterface1(t *testing.T) {
	ps := &CopyTextStruct{
		Name: "foo",
		Val:  2.2,
		Ref:  nil,
	}
	ps.Ref = ps
	var i IntForCopy = ps
	cp := Copy(reflect.ValueOf(i))
	if cp.Kind() != reflect.Ptr {
		t.Fatal("wrong kind", cp.Kind())
	}
	if !reflect.DeepEqual(cp.Interface(), i) {
		t.Fatalf("copy failed: %#v", cp.Interface())
	}
	i.Foo("new name")
	cpps := cp.Interface().(*CopyTextStruct)
	if cpps.Name != "foo" || cpps.Val != 2.2 {
		t.Fatalf("copy failed: %#v", cp.Interface())
	}
}

func TestCopyInterface2(t *testing.T) {
	var i IntForCopy
	cp := Copy(reflect.ValueOf(i))
	if cp.IsValid() {
		t.Fatal("very wrong")
	}
}

type TestInterface struct {
	Name string
	Int  interface{}
}

func TestCopyInterface3(t *testing.T) {
	s := &TestInterface{
		Name: "foo",
	}
	cp := Copy(reflect.ValueOf(s))
	if !reflect.DeepEqual(cp.Interface(), s) {
		t.Fatalf("copy failed: %#v", cp.Interface())
	}
}

func TestAnySettableValue(t *testing.T) {
	ti := &TestInterface{
		Name: "foo",
	}
	if !AnySettableValue(reflect.ValueOf(ti)) {
		t.Fatal("AnySettableValue failed.")
	}
	err := errors.New("foo")
	if AnySettableValue(reflect.ValueOf(err)) {
		t.Fatal("AnySettableValue failed.")
	}
}

type NotRecursive1 int

type NotRecursive2 struct {
	A string
}

type Recursive1 struct {
	ptr *Recursive1
}

func TestIsRecursive(t *testing.T) {
	if isRecursive(reflect.ValueOf(NotRecursive1(0))) {
		t.Fatal("recursive")
	}
	if isRecursive(reflect.ValueOf(NotRecursive2{})) {
		t.Fatal("recursive")
	}
	if isRecursive(reflect.ValueOf(&NotRecursive2{})) {
		t.Fatal("recursive")
	}
	if !isRecursive(reflect.ValueOf(Recursive1{})) {
		t.Fatal("not recursive")
	}
	if !isRecursive(reflect.ValueOf(&Recursive1{})) {
		t.Fatal("not recursive")
	}
}
