package export

import (
	"bytes"
	"fmt"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestMarshalXML(t *testing.T) {
	b := &bytes.Buffer{}
	x := newXMLExporter(b, Config{Namespace: "jobs", Element: "job"})
	err := x.Write(map[string]interface{}{
		"bool":    true,
		"int":     -123,
		"int8":    int8(-123),
		"int16":   int16(-655),
		"int32":   int32(-123456789),
		"int64":   int64(-123456789012),
		"uint":    uint(123),
		"uint8":   uint8(123),
		"uint16":  uint16(35000),
		"uint32":  uint32(123456789),
		"uint64":  uint64(123456789012),
		"float":   123.456,
		"float32": float32(123.456789),
		"float64": float64(123.456789012),
		"string":  "Oh, hello, this is a string. Also it has <reserved> characters in it. \"How about that!\"",
		"array":   []string{"A", "B", "C"},
		"map": map[string]interface{}{
			"sub":     123,
			"another": false,
		},
	})
	assert.Nil(t, err, fmt.Sprintf("%v", err))
	err = x.Flush()
	assert.Nil(t, err, fmt.Sprintf("%v", err))
	fmt.Println(string(b.Bytes()))
}
