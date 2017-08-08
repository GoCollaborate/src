package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func ReflectFuncName(fun interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	return name
}

func StripRouteToAPIRoute(rt string) string {
	return strings.Replace(strings.TrimPrefix(rt, "github.com/GoCollaborate/"), ".", "/", -1)
}
