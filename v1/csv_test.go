package export

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCSV struct {
	Values []interface{}
	Expect string
	Error  error
}

func TestExportCSV(t *testing.T) {
	tests := []testCSV{
		{
			Values: []interface{}{
				struct {
					A string `csv:"a"`
					B int    `csv:"b"`
					C bool   `csv:"c"`
				}{
					"Hello",
					123,
					true,
				},
				struct {
					A string `csv:"a"`
					B int    `csv:"b"`
					C bool   `csv:"c"`
				}{
					"Another",
					987,
					false,
				},
			},
			Expect: "a,b,c\nHello,123,true\nAnother,987,false\n",
		},
	}

outer:
	for _, e := range tests {
		b := &bytes.Buffer{}
		x := newCSVExporter(b, Config{})
		for _, v := range e.Values {
			err := x.Write(v)
			if err != nil {
				assert.Equal(t, e.Error, err)
				continue outer
			}
		}

		err := x.Flush()
		if err != nil {
			assert.Equal(t, e.Error, err)
			continue outer
		}

		fmt.Println(e.Expect)
		fmt.Println(string(b.Bytes()))
		assert.Equal(t, e.Expect, string(b.Bytes()))
	}
}
