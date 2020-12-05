package export

import (
	"mime"
)

type mimetype string

const (
	JSON = mimetype("application/json")
	XML  = mimetype("text/xml")
	CSV  = mimetype("text/csv")
	XLSX = mimetype("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
)

func parseMIME(v string) (mimetype, error) {
	m, _, err := mime.ParseMediaType(v)
	if err != nil {
		return "", err
	}
	return mimetype(m), nil
}
