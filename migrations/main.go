package migrations

import (
	"fmt"
	"reflect"
)

const (
	ACTION_CREATE = "create"
	ACTION_UP     = "up"
	ACTION_DOWN   = "down"
	ACTION_UP_ALL = "up_all"
)

type Migration struct {
	UpFunc   map[string]interface{}
	DownFunc map[string]interface{}
}

var migrations = &Migration{
	UpFunc:   make(map[string]interface{}),
	DownFunc: make(map[string]interface{}),
}

func RegisterUpFunc(name string, function interface{}) {
	migrations.UpFunc[name] = function
}

func RegisterDownFunc(name string, function interface{}) {
	migrations.DownFunc[name] = function
}

func Run(name string, action string, args ...interface{}) error {
	// Check if the function exists
	var (
		fn interface{}
		ok bool
	)
	switch action {
	case ACTION_UP:
		fn, ok = migrations.UpFunc[name]
		if !ok {
			return fmt.Errorf("function %s not found", name)
		}
	case ACTION_DOWN:
		fn, ok = migrations.DownFunc[name]
		if !ok {
			return fmt.Errorf("function %s not found", name)
		}
	default:
		return fmt.Errorf("function %s not supported", name)
	}

	var inputArgs []reflect.Value
	for _, arg := range args {
		inputArgs = append(inputArgs, reflect.ValueOf(arg))
	}

	result := reflect.ValueOf(fn).Call(inputArgs)

	// Check for errors in the result
	for _, r := range result {
		if r.Type().String() == "error" && !r.IsNil() {
			fmt.Printf("error run migrate %s: %s\n", name, r.Interface())

			continue
		}
	}

	return nil // Function call was successful
}
