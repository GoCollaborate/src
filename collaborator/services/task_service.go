package services

import (
	"bytes"
	"encoding/gob"
	. "github.com/GoCollaborate/src/artifacts/task"
)

func Encode(maps *map[int]*Task) (payload *TaskPayload, err error) {
	var maps_bytes bytes.Buffer

	enc := gob.NewEncoder(&maps_bytes)

	err = enc.Encode(maps)

	payload = &TaskPayload{
		Payload: maps_bytes.Bytes(),
	}
	return
}

func Decode(payload *TaskPayload) (maps *map[int]*Task, err error) {
	dec := gob.NewDecoder(bytes.NewReader(payload.GetPayload()))
	err = dec.Decode(&maps)
	return
}
