package generator

import uuid "github.com/satori/go.uuid"

// UUID ...
func UUID() string {
	u4 := uuid.NewV4()
	return u4.String()
}
