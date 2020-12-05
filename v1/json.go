package export

import (
	"encoding/json"
	"io"
)

type jsonExporter struct {
	io.Writer
}

func newJSONExporter(w io.Writer, c Config) Exporter {
	return &jsonExporter{w}
}

func (e *jsonExporter) Write(record interface{}) error {
	err := json.NewEncoder(e.Writer).Encode(record)
	if err != nil {
		return err
	}
	return nil
}

func (e *jsonExporter) Flush() error {
	return nil
}
