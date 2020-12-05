package export

import (
	"encoding"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

var typeOfTextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

// XML exporter
type xmlExporter struct {
	*xml.Encoder
	root, element xml.Name
	init, safe    bool
}

// Create an XML exporter
func newXMLExporter(w io.Writer, c Config) Exporter {
	var elem string
	if c.Element != "" {
		elem = c.Element
	} else {
		elem = "item"
	}
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	return &xmlExporter{e, xml.Name{Local: c.Namespace}, xml.Name{Local: elem}, false, c.Safe}
}

// Write a record
func (e *xmlExporter) Write(raw interface{}) error {
	record, err := convertRecord(XML, raw)
	if err != nil {
		return err
	}

	if !e.init && e.root.Local != "" {
		err := e.EncodeToken(xml.StartElement{Name: e.root})
		if err != nil {
			return err
		}
		e.init = true
	}
	return e.writeElement(e.element, reflect.ValueOf(record))
}

// Write a record
func (e *xmlExporter) writeElement(name xml.Name, val reflect.Value) error {
	if !val.IsValid() {
		return nil // do nothing
	}
	err := e.EncodeToken(xml.StartElement{Name: name})
	if err != nil {
		return err
	}
	err = e.writeValue(name, val)
	if err != nil {
		return err
	}
	err = e.EncodeToken(xml.EndElement{Name: name})
	if err != nil {
		return err
	}
	return e.Encoder.Flush()
}

// Write a record
func (e *xmlExporter) writeValue(name xml.Name, val reflect.Value) error {
	var err error
	if !val.IsValid() {
		// nothing
	} else if val.Type().Implements(typeOfTextMarshaler) {
		err = e.writeCustom(name, val)
	} else {
		switch val.Kind() {
		case reflect.Invalid: // ignore
		case reflect.Ptr:
			err = e.writeValue(name, reflect.Indirect(val))
		case reflect.Interface:
			err = e.writeValue(name, val.Elem())
		case reflect.String:
			err = e.EncodeToken(xml.CharData(val.String()))
		case reflect.Bool:
			err = e.EncodeToken(xml.CharData(strconv.FormatBool(val.Bool())))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = e.EncodeToken(xml.CharData(strconv.FormatInt(val.Int(), 10)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = e.EncodeToken(xml.CharData(strconv.FormatUint(val.Uint(), 10)))
		case reflect.Float32, reflect.Float64:
			err = e.EncodeToken(xml.CharData(strconv.FormatFloat(val.Float(), 'f', -1, 64)))
		case reflect.Map:
			err = e.writeMap(name, val)
		case reflect.Slice, reflect.Array:
			err = e.writeSlice(name, val)
		default:
			err = e.writeCustom(name, val)
		}
	}
	return err
}

// Write a map
func (e *xmlExporter) writeMap(name xml.Name, val reflect.Value) error {
	for _, k := range val.MapKeys() {
		if k.Kind() != reflect.String {
			return fmt.Errorf("Invalid key: %v", k.Type())
		}
		err := e.writeElement(xml.Name{Local: k.Interface().(string)}, val.MapIndex(k))
		if err != nil {
			return err
		}
	}
	return nil
}

// Write a slice
func (e *xmlExporter) writeSlice(name xml.Name, val reflect.Value) error {
	l := val.Len()
	for i := 0; i < l; i++ {
		err := e.writeElement(xml.Name{Local: "elem"}, val.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

// Write a type with a custom marshaler
func (e *xmlExporter) writeCustom(name xml.Name, val reflect.Value) error {
	t := val.Type()
	if t.Implements(typeOfTextMarshaler) {
		v, err := val.Interface().(encoding.TextMarshaler).MarshalText()
		if err != nil {
			return err
		}
		return e.writeValue(name, reflect.ValueOf(string(v)))
	}
	return fmt.Errorf("Unsupported type: %v", val.Type())
}

// Flush
func (e *xmlExporter) Flush() error {
	if e.init && e.root.Local != "" {
		err := e.EncodeToken(xml.EndElement{Name: e.root})
		if err != nil {
			return err
		}
	}
	return e.Encoder.Flush()
}
