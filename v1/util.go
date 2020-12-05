package export

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	errUnsupportedType   = fmt.Errorf("Unsupported type")
	errUnsupportedFormat = fmt.Errorf("Unsupported output format")
)

// Convert a value to a single record for the provided type.
func convertRecord(m mimetype, v interface{}) (map[string]interface{}, error) {
	x := reflect.Indirect(reflect.ValueOf(v))
	switch x.Kind() {
	case reflect.Struct:
		return convertStruct(m, x)
	case reflect.Map:
		return convertMap(m, x)
	default:
		return nil, errUnsupportedType
	}
}

func convertStruct(m mimetype, v reflect.Value) (map[string]interface{}, error) {
	var tagname string
	switch m {
	case JSON:
		tagname = "json"
	case CSV:
		tagname = "csv"
	case XLSX:
		tagname = "csv"
	default:
		return nil, errUnsupportedFormat
	}

	r := make(map[string]interface{})
	n := v.NumField()
	t := v.Type()

	for i := 0; i < n; i++ {
		f := t.Field(i)
		d := v.Field(i)
		if d.IsValid() && d.CanInterface() {
			tag := f.Tag.Get(tagname)
			if tag == "" {
				tag = f.Name
			} else if x := strings.Index(tag, ","); x >= 0 {
				tag = tag[:x]
			}
			r[strings.TrimSpace(tag)] = d.Interface()
		}
	}

	return r, nil
}

func convertMap(m mimetype, v reflect.Value) (map[string]interface{}, error) {
	t := v.Type()
	if t.Key().Kind() != reflect.String {
		return nil, fmt.Errorf("Invalid map key type")
	}

	r := make(map[string]interface{})
	iter := v.MapRange()
	for iter.Next() {
		k := iter.Key()
		d := iter.Value()
		if d.IsValid() && d.CanInterface() {
			r[k.String()] = d.Interface()
		}
	}

	return r, nil
}
