package export

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testJSON struct {
	Values []interface{}
	Expect string
	Error  error
}

func TestExportJSON(t *testing.T) {
	tests := []testJSON{
		{
			Values: []interface{}{
				struct {
					A string `json:"a"`
					B int    `json:"b"`
					C bool   `json:"c"`
				}{
					"Hello",
					123,
					true,
				},
				struct {
					A string `json:"a"`
					B int    `json:"b"`
					C bool   `json:"c"`
				}{
					"Another",
					987,
					false,
				},
			},
			Expect: `{"a":"Hello","b":123,"c":true}
{"a":"Another","b":987,"c":false}
`,
		},
	}

outer:
	for _, e := range tests {
		b := &bytes.Buffer{}
		x := newJSONExporter(b, Config{})
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

		fmt.Println(string(e.Expect))
		fmt.Println(string(b.Bytes()))
		assert.Equal(t, e.Expect, string(b.Bytes()))
	}
}
