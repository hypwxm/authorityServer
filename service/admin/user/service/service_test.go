package service

import (
	"log"
	"testing"
)

func TestModels(t *testing.T) {
	err := InitAdmin()
	if err != nil {
		log.Fatal(err)
	}
}
