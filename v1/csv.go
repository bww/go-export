package export

import (
	"encoding/csv"
	"io"
	"sort"
)

type csvExporter struct {
	*csv.Writer
	header []string
	safe   bool
}

func newCSVExporter(w io.Writer, c Config) Exporter {
	return &csvExporter{csv.NewWriter(w), nil, c.Safe}
}

func (e *csvExporter) Write(raw interface{}) error {
	record, err := convertRecord(CSV, raw)
	if err != nil {
		return err
	}

	if e.header == nil {
		h := make([]string, len(record))
		i := 0
		for k, _ := range record {
			h[i] = k
			i++
		}
		sort.Strings(h)
		err := e.Writer.Write(h)
		if err != nil {
			return err
		}
		e.header = h
	}

	values := make([]string, len(e.header))
	for i, x := range e.header {
		if v := record[x]; v == nil {
			values[i] = ""
		} else if e.safe {
			values[i] = escapeCSV(stringer(v))
		} else {
			values[i] = stringer(v)
		}
	}

	err = e.Writer.Write(values)
	if err != nil {
		return err
	}
	return nil
}

func (e *csvExporter) Flush() error {
	e.Writer.Flush()
	return e.Writer.Error()
}

func escapeCSV(v string) string {
	if len(v) > 0 {
		switch v[0] {
		case '=', '-', '+', '@':
			return "\t" + v
		}
	}
	return v
}
