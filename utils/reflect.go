package utils

import (
	"github.com/GoCollaborate/src/constants"
	"github.com/fatih/structs"
	"reflect"
	"runtime"
	"strings"
)

func ReflectFuncName(fun interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
	return name
}

func StripRouteToAPIRoute(rt string) string {
	return strings.Replace(strings.TrimPrefix(rt, "_"+constants.PROJECT_DIR), ".", "/", -1)
}

func StripRouteToFunctName(rt string) string {
	return strings.Replace(strings.TrimPrefix(rt, "_"+constants.PROJECT_DIR), ".", "/", -1)
}

func Map(m interface{}) map[string]interface{} {
	return structs.Map(m)
}
