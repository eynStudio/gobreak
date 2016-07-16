package jsonx

import (
	"encoding/json"
	"io"
)

func Encode(v interface{}) io.Reader {
	reader, writer := io.Pipe()
	go func() {
		err := json.NewEncoder(writer).Encode(v)
		if err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()
	return reader
}
