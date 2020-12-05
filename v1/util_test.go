package export

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type testConvert struct {
	Mime   mimetype
	Value  interface{}
	Expect map[string]interface{}
	Error  error
}

func TestConvertStruct(t *testing.T) {
	tests := []testConvert{
		{
			Mime: JSON,
			Value: struct {
				A string `json:"a"`
				B int    `json:"b"`
			}{
				"Hello",
				123,
			},
			Expect: map[string]interface{}{
				"a": "Hello",
				"b": 123,
			},
		},
		{
			Mime: CSV,
			Value: struct {
				A string
				B int
			}{
				"Hello",
				123,
			},
			Expect: map[string]interface{}{
				"A": "Hello",
				"B": 123,
			},
		},
		{
			Mime: CSV,
			Value: struct {
				A string `csv:"a_csv"`
				B int    `csv:"b_csv"`
			}{
				"Hello",
				123,
			},
			Expect: map[string]interface{}{
				"a_csv": "Hello",
				"b_csv": 123,
			},
		},
		{
			Mime: JSON,
			Value: map[string]string{
				"a_map": "Hello",
				"b_map": "123",
			},
			Expect: map[string]interface{}{
				"a_map": "Hello",
				"b_map": "123",
			},
		},
		{
			Mime: JSON,
			Value: map[string]bool{
				"a_map": true,
				"b_map": false,
			},
			Expect: map[string]interface{}{
				"a_map": true,
				"b_map": false,
			},
		},
		{
			Mime:  JSON,
			Value: 123,
			Error: errUnsupportedType,
		},
		{
			Mime:  mimetype("image/gif"),
			Value: struct{}{},
			Error: errUnsupportedFormat,
		},
		{
			Mime:   mimetype("image/gif"),
			Value:  map[string]bool{}, // don't care about maps
			Expect: map[string]interface{}{},
		},
	}

	for _, e := range tests {
		r, err := convertRecord(e.Mime, e.Value)
		spew.Dump(e.Value)
		spew.Dump(r)
		if assert.Equal(t, e.Error, err) {
			assert.Equal(t, e.Expect, r)
		}
	}
}
