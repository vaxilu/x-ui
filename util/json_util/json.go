package json_util

import (
	"encoding/json"
	"reflect"
	"x-ui/util/reflect_util"
)

/*
MarshalJSON 特殊处理 json.RawMessage

当 json.RawMessage 不为 nil 且 len() 为 0 时，MarshalJSON 将会解析报错
*/
func MarshalJSON(i interface{}) ([]byte, error) {
	m := map[string]interface{}{}
	t := reflect.TypeOf(i).Elem()
	v := reflect.ValueOf(i).Elem()
	fields := reflect_util.GetFields(t)
	for _, field := range fields {
		key := field.Tag.Get("json")
		if key == "" || key == "-" {
			continue
		}
		fieldV := v.FieldByName(field.Name)
		value := fieldV.Interface()
		switch value.(type) {
		case json.RawMessage:
			value := value.(json.RawMessage)
			if len(value) > 0 {
				m[key] = value
			}
		default:
			m[key] = value
		}
	}
	return json.Marshal(m)
}
