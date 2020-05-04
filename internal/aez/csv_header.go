package aez

import "reflect"

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func csvHeader(obj interface{}, exclude ...string) []string {
	var fields []string
	val := reflect.Indirect(reflect.ValueOf(obj))
	t := val.Type()
	for i, l := 0, val.NumField(); i < l; i++ {
		if name := t.Field(i).Name; !stringInSlice(name, exclude) {
			fields = append(fields, t.Field(i).Name)
		}
	}
	return fields
}
