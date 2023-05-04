package main

import "reflect"

type typ interface {
	exec() interface{}
}

type scalar struct {
	resolver reflect.Value
}

type object struct {
	fields map[string]typ
}

type parseError string

type schema struct {
	types map[string]*object
}

func nativeGqlServer() {

}
