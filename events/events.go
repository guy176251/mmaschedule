package events

import (
	"bytes"
	"encoding/gob"
)

func decode(data []byte, value any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(value)
}

func encode(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	return buf.Bytes(), err
}

func (e *DbEvent) ReadData(value any) error {
	return decode(e.Data, value)
}

func (e *DbEvent) WriteData(value any) error {
	b, err := encode(value)
	if err == nil {
		e.Data = b
	}
	return err
}

func (e *DbEvent) ReadHistory(value any) error {
	return decode(e.History, value)
}

func (e *DbEvent) WriteHistory(value any) error {
	b, err := encode(value)
	if err == nil {
		e.History = b
	}
	return err
}
