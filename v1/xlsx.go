package export

import (
	"io"
	"sort"
	"time"

	"github.com/tealeg/xlsx"
)

type xlsxExporter struct {
	writer io.Writer
	file   *xlsx.File
	sheet  *xlsx.Sheet
	header []string
	conf   Config
}

func newXLSXExporter(w io.Writer, c Config) Exporter {
	f := xlsx.NewFile()
	return &xlsxExporter{w, f, nil, nil, c}
}

func (e *xlsxExporter) Title() string {
	if e.conf.Title != "" {
		return e.conf.Title
	} else {
		return "Report"
	}
}

func (e *xlsxExporter) Write(raw interface{}) error {
	record, err := convertRecord(XLSX, raw)
	if err != nil {
		return err
	}

	if e.sheet == nil {
		e.sheet, err = e.file.AddSheet(e.Title())
		if err != nil {
			return err
		}
	}

	if e.header == nil {
		h := make([]string, len(record))
		i := 0
		for k, _ := range record {
			h[i] = k
			i++
		}
		sort.Strings(h)
		r := e.sheet.AddRow()
		for _, x := range h {
			r.AddCell().Value = x
		}
		e.header = h
	}

	r := e.sheet.AddRow()
	for _, x := range e.header {
		if v := record[x]; v == nil {
			r.AddCell()
		} else {
			c := r.AddCell()
			switch t := v.(type) {
			case int8:
				c.SetInt(int(t))
			case uint8:
				c.SetInt(int(t))
			case int16:
				c.SetInt(int(t))
			case uint16:
				c.SetInt(int(t))
			case int32:
				c.SetInt(int(t))
			case uint32:
				c.SetInt(int(t))
			case int:
				c.SetInt(t)
			case uint:
				c.SetInt(int(t))
			case int64:
				c.SetInt64(t)
			case uint64:
				c.SetInt64(int64(t))
			case time.Time:
				c.SetDate(t)
			default:
				c.Value = stringer(v)
			}
		}
	}

	return nil
}

func (e *xlsxExporter) Flush() error {
	if e.sheet == nil {
		_, err := e.file.AddSheet(e.Title())
		if err != nil {
			return err
		}
	}
	return e.file.Write(e.writer)
}
