package gordb

import "reflect"

type Attr struct {
	Name string
	reflect.Kind
}

type Schema []Attr
