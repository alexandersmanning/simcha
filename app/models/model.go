//Generic functions and types to be used by all models

package models

import (
	"fmt"
)

type modelError struct {
	fieldName string
	errorText string
}

func (m *modelError) Error() string {
	return fmt.Sprintf("%s %s", m.fieldName, m.errorText)
}
