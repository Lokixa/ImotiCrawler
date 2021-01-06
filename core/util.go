package core

import "reflect"

// GetStructFields returns fields in string format if table is a struct
func GetStructFields(table interface{}) []string {
	valType := reflect.ValueOf(table)
	result := make([]string, 0)

	stack := make([]reflect.StructField, 0)
	for i := 0; i < valType.NumField(); i++ {
		stack = append(stack, valType.Type().Field(i))
	}
	stackLen := len(stack)
	for stackLen > 0 {
		fieldType := stack[stackLen-1]
		stack = stack[:stackLen-1]
		stackLen = len(stack)

		if fieldType.Type.Kind() == reflect.Struct && fieldType.Type.NumField() > 0 {
			for i := 0; i < fieldType.Type.NumField(); i++ {
				stack = append(stack, fieldType.Type.Field(i))
			}
		} else {
			result = append(result, fieldType.Name)
		}
		stackLen = len(stack)
	}
	resLen := len(result)
	for i := 0; i < resLen/2; i++ {
		result[i], result[resLen-i-1] = result[resLen-i-1], result[i]
	}
	return result
}
