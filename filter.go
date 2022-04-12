package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func main() {
	s := struct {
		A string
		B float64
	}{A: "qwe", B: 1.1}
	i, _ := FilterStructByFields(s, []string{"a", "c"})
	fmt.Printf("%+v", i)
}

func FilterStructByFields(s interface{}, f []string) (interface{}, error) {
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("parameter must be struct type")
	}
	m := make(map[string]struct{}, 0)
	for _, field := range f {
		m[strings.ToLower(field)] = struct{}{}
	}

	t := reflect.TypeOf(struct{}{})
	fields := reflect.VisibleFields(t)

	//создаём поля
	for i := 0; i < typ.NumField(); i++ {
		fName := typ.Field(i).Name
		fType := typ.Field(i).Type
		if _, ok := m[strings.ToLower(fName)]; !ok {
			continue
		}
		tag := strings.ToLower(fName)
		fields = append(fields, reflect.StructField{
			Name: fName,
			Type: fType,
			Tag:  reflect.StructTag(`json:"` + tag + `"`),
		})
	}

	structOf := reflect.StructOf(fields)
	n := reflect.New(structOf).Elem()

	//заполняем структурку
	for i := 0; i < n.NumField(); i++ {
		_, ok := typ.FieldByName(structOf.Field(i).Name)
		if !ok {
			continue
		}
		fieldByIndex := n.Field(i)
		switch fieldByIndex.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldByIndex.SetInt(val.FieldByName(structOf.Field(i).Name).Int())
		case reflect.String:
			fieldByIndex.SetString(val.FieldByName(structOf.Field(i).Name).String())
		case reflect.Float32, reflect.Float64:
			fieldByIndex.SetFloat(val.FieldByName(structOf.Field(i).Name).Float())
		}

	}

	strct := n.Interface()

	return strct, nil
}
