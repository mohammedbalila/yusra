package internal

import "reflect"

type Column struct {
	Name     string
	DataType reflect.Type
}
