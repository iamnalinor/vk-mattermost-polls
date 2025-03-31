package repo

import (
	"fmt"
	"reflect"
)

func asInt(val any) int {
	return int(reflect.ValueOf(val).Int())
}

func asStringSlice(val any) []string {
	slice := val.([]any)
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprint(v)
	}
	return result
}
