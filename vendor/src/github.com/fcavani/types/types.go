// Copyright 2015 Felipe A. Cavani. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Start date:		2010-08-11

// Package types have functions to create an instatiation of one type from the type name.
package types

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var typemap map[string]reflect.Type

func init() {
	typemap = make(map[string]reflect.Type, 100)
	Insert(errors.New(""))
	InsertName("os.errorString", errors.New(""))
	Insert("")
	Insert(new(string))
	Insert(int(0))
	Insert(int8(0))
	Insert(int16(0))
	Insert(int32(0))
	Insert(int64(0))
	Insert(uint(0))
	Insert(uint8(0))
	Insert(uint16(0))
	Insert(uint32(0))
	Insert(uint64(0))
	Insert(new(int8))
	Insert(new(int16))
	Insert(new(int32))
	Insert(new(int64))
	Insert(new(uint))
	Insert(new(uint8))
	Insert(new(uint16))
	Insert(new(uint32))
	Insert(new(uint64))
	Insert(new(int))
	Insert(false)
	Insert(new(bool))
	Insert(float32(3.14))
	Insert(new(float32))
	Insert(float64(3.14))
	Insert(new(float64))
	Insert(time.Time{})
	Insert(&time.Time{})
	Insert(time.Duration(0))
}

// Dump the name and the type from the type base.
func Dump() {
	for key, t := range typemap {
		fmt.Println(key, t)
	}
}

func pkgname() (name string) {
	pc, _, _, ok := runtime.Caller(4)
	if !ok {
		return
	}
	f := runtime.FuncForPC(pc)
	s := strings.SplitN(f.Name(), ".", 2)
	if len(s) != 2 {
		return
	}
	name = s[0]
	return
}

func findpkgname(t reflect.Type) (name string) {
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		name = findpkgname(t.Elem())
		if name == "" {
			name = t.PkgPath()
		}
	default:
		name = t.PkgPath()
	}
	return
}

func replacepkgname(in string, t reflect.Type) (out string) {
	pkg := findpkgname(t)
	s := strings.Split(pkg, "/")
	if len(s) <= 0 {
		return
	}
	out = strings.Replace(in, s[len(s)-1], pkg, 1)
	return
}

func nameof(t reflect.Type) (name string) {
	n := t.Name()
	if t.Kind() == reflect.Interface || n == "" {
		name = replacepkgname(t.String(), t)
		if name == "" {
			name = t.String()
		}
	} else {
		pkg := t.PkgPath()
		if pkg == "" {
			name = n
		} else {
			name = pkg + "." + n
		}
	}
	return
}

// NameOf returns the package name and the name of the type
func NameOf(t reflect.Type) string {
	return nameof(t)
}

// Name accepts a variable of any type and returns the package
// name and the name of the type or a function
func Name(i interface{}) string {
	val := reflect.ValueOf(i)
	if !val.IsValid() {
		return ""
	}
	t := val.Type()
	switch t.Kind() {
	case reflect.Func:
		return runtime.FuncForPC(val.Pointer()).Name()
	default:
		return nameof(t)
	}
}

// InsertType insertes a type for future instantiation.
// Do this in the same package where the type was declared.
// The use of  init function is advised.
func InsertType(t reflect.Type) {
	tname := nameof(t)
	if _, found := typemap[tname]; !found {
		typemap[tname] = t
	}
}

// Insert type for future instantiation.
// Do this in the same package where the type was declared.
// The use of  init function is advised.
func Insert(i interface{}) {
	t := reflect.ValueOf(i).Type()
	tname := nameof(t)
	if _, found := typemap[tname]; !found {
		typemap[tname] = t
	}
}

// InsertName inserts a new type with the name.
func InsertName(tname string, i interface{}) {
	t := reflect.ValueOf(i).Type()
	if _, found := typemap[tname]; !found {
		typemap[tname] = t
	}
}

// Type returns the Type from the type name.
func Type(tname string) reflect.Type {
	if t, found := typemap[tname]; found {
		return t
	}
	panic("Type not found: " + tname)
}

// GetType return the type represented by tname
func GetType(tname string) (reflect.Type, error) {
	t, found := typemap[tname]
	if !found {
		return nil, errors.New("type not found: " + tname)
	}
	return t, nil

}

// IsEqualName compares the value type name with one name.
func IsEqualName(val reflect.Value, tname string) bool {
	return nameof(val.Type()) == tname
}

// MakeZero creates a zero value type for the type name.
func MakeZero(tname string) reflect.Value {
	return reflect.Zero(Type(tname))
}

// MakeNew create a new value from the type's name
func MakeNew(tname string, bufcap int) (val reflect.Value) {
	t := Type(tname)
	val = MakeNewType(t, bufcap)
	return
}

// MakeNewType creates a new value with type t.
func MakeNewType(t reflect.Type, bufcap int) (val reflect.Value) {
	switch t.Kind() {
	case reflect.Ptr:
		val = reflect.New(t).Elem()
		val.Set(reflect.New(val.Type().Elem()))
	case reflect.Chan:
		//typ := reflect.ChanOf(reflect.BothDir, t.Elem())
		val = reflect.New(t).Elem()
		val.Set(reflect.MakeChan(t, bufcap)) //TODO: set buf?
	case reflect.Slice:
		val = reflect.New(t).Elem()
		val.Set(reflect.MakeSlice(t, 0, bufcap))
		val.SetLen(bufcap)
	case reflect.Map:
		val = reflect.New(t).Elem()
		val.Set(reflect.MakeMap(t))
	default:
		val = reflect.New(t).Elem()
	}
	return
}

func isRecursive(v reflect.Value) bool {
	ind := reflect.Indirect(v)
	t := v.Type()
	if ind.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < ind.Type().NumField(); i++ {
		ft := ind.Field(i).Type()
		return ft == t || ft == reflect.PtrTo(t) || t.AssignableTo(ft) || reflect.PtrTo(t).AssignableTo(ft)
	}
	return false
}

//AllocStructPtrs find pointers in a struct and alloc than recursivily.
func AllocStructPtrs(v reflect.Value) {
	val := reflect.Indirect(v)
	t := val.Type()

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := val.Field(i)
			switch field.Type().Kind() {
			case reflect.Ptr:
				v := MakeNewType(field.Type(), 0)
				if isRecursive(v) {
					//panic(fmt.Sprintf("struct %v have a field of the same type of this struct", NameOf(v.Type())))
					continue
				}
				AllocStructPtrs(v)
				if field.CanSet() {
					field.Set(v)
				}
			case reflect.Slice:
				v := MakeNewType(field.Type(), 0)
				if field.CanSet() {
					field.Set(v)
				}
			default:
				continue
			}
		}
	}
	return
}

// Make instantiate a value of t type and allocate pointer and slices.
func Make(t reflect.Type) (val reflect.Value) {
	val = MakeNewType(t, 0)
	AllocStructPtrs(val)
	return
}

type deepcopy map[reflect.Value]reflect.Value

func (d deepcopy) copy(src reflect.Value) (dst reflect.Value) {
	switch src.Kind() {
	case reflect.Bool:
		dst = reflect.New(src.Type()).Elem()
		dst.SetBool(src.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst = reflect.New(src.Type()).Elem()
		dst.SetInt(src.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dst = reflect.New(src.Type()).Elem()
		dst.SetUint(src.Uint())
	case reflect.Uintptr:
		panic("Uintptr isn't supported")
	case reflect.Float32, reflect.Float64:
		dst = reflect.New(src.Type()).Elem()
		dst.SetFloat(src.Float())
	case reflect.Complex64, reflect.Complex128:
		dst = reflect.New(src.Type()).Elem()
		dst.SetComplex(src.Complex())
	case reflect.Array:
		dst = reflect.New(src.Type()).Elem()
		for i := 0; i < src.Type().Len(); i++ {
			dst.Index(i).Set(d.copy(src.Index(i)))
		}
	case reflect.Chan:
		dst = reflect.New(src.Type()).Elem()
		dst.Set(reflect.MakeChan(src.Type(), src.Cap()))
	case reflect.Func:
		dst = src
	case reflect.Interface:
		dst = reflect.New(src.Type()).Elem()
		if !src.Elem().IsValid() {
			return
		}
		dst.Set(d.copy(src.Elem()))
	case reflect.Map:
		dst = reflect.New(src.Type()).Elem()
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			dst.SetMapIndex(d.copy(key), d.copy(src.MapIndex(key)))
		}
	case reflect.Ptr:
		if dst, found := d[src]; found {
			return dst
		}
		dst = reflect.New(src.Type()).Elem()
		d[src] = dst
		if !src.Elem().IsValid() {
			return
		}
		dst.Set(reflect.New(src.Type().Elem()))
		val := d.copy(src.Elem())
		dst.Elem().Set(val)
	case reflect.Slice:
		dst = reflect.New(src.Type()).Elem()
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(d.copy(src.Index(i)))
		}
	case reflect.String:
		dst = reflect.New(src.Type()).Elem()
		dst.SetString(src.String())
	case reflect.Struct:
		dst = reflect.New(src.Type()).Elem()
		for i := 0; i < src.Type().NumField(); i++ {
			if dst.Field(i).CanSet() {
				dst.Field(i).Set(d.copy(src.Field(i)))
			}
		}
	case reflect.UnsafePointer:
		panic("UnsafePointer isn't supported")
	case reflect.Invalid:
		return
	default:
		panic(fmt.Sprintf("kind %v is not supported", src.Kind()))
	}
	return
}

// Copy create a new value with all data of src copied into it.
func Copy(src reflect.Value) reflect.Value {
	d := make(deepcopy)
	return d.copy(src)
}

//AnySettableValue find if exist one value that you can set.
func AnySettableValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Array:
		for i := 0; i < val.Type().Len(); i++ {
			if AnySettableValue(val.Index(i)) {
				return true
			}
		}
	case reflect.Interface:
		return AnySettableValue(val.Elem())
	case reflect.Map:
		for _, key := range val.MapKeys() {
			if AnySettableValue(val.MapIndex(key)) {
				return true
			}
		}
	case reflect.Ptr:
		return AnySettableValue(val.Elem())
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if AnySettableValue(val.Index(i)) {
				return true
			}
		}
	case reflect.Struct:
		for i := 0; i < val.Type().NumField(); i++ {
			if AnySettableValue(val.Field(i)) {
				return true
			}
		}
	default:
		return val.CanSet()
	}
	return false
}
