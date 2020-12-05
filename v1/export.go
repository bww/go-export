package export

import (
	"errors"
	"fmt"
	"io"
)

var ErrUnsupported = errors.New("Unsupported export format")

type Config struct {
	Safe      bool
	Title     string // the title of the exported document
	Namespace string // the namespace for all exported elements
	Element   string // the name of a single export element
}

// An exporter exports a bunch of records
type Exporter interface {
	Write(interface{}) error
	Flush() error
}

// Create an exporter for the specified output type
func New(t string, w io.Writer) (Exporter, error) {
	return NewWithConfig(t, Config{}, w)
}

// Create an exporter for the specified output type
func NewWithConfig(t string, c Config, w io.Writer) (Exporter, error) {
	m, err := parseMIME(t)
	if err != nil {
		return nil, err
	}
	switch m {
	case CSV:
		return newCSVExporter(w, c), nil
	case XLSX:
		return newXLSXExporter(w, c), nil
	case XML:
		return newXMLExporter(w, c), nil
	case JSON:
		return newJSONExporter(w, c), nil
	default:
		return nil, ErrUnsupported
	}
}

// Convert to a string the best way possible
func stringer(v interface{}) string {
	switch c := v.(type) {
	case string:
		return c
	case *string:
		return *c
	case fmt.Stringer:
		return c.String()
	case error:
		return c.Error()
	default:
		return fmt.Sprint(v)
	}
}
