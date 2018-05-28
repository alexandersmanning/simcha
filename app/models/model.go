package models

import (
	"fmt"
	"encoding/gob"
	"time"
)

type ModelAction interface {
	Timestamps() (time.Time, time.Time)
	SetID(id int)
	SetTimestamps()
}
//TODO Move to models package
type ModelError struct {
	FieldName string
	ErrorText string
}

func (m *ModelError) Error() string {
	return fmt.Sprintf("%s %s", m.FieldName, m.ErrorText)
}

func init() {
	gob.Register(&User{})
}

