package util

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestFormatQuota(t *testing.T) {
	var s = "awdawdwd'123'123'1(23()"
	s = FormatQuota(s)
	t.Fatal(s)
}


func TestGetUuid(t *testing.T) {
	s := uuid.NewV5(uuid.NewV4(), "123")
	t.Log(s)
}
