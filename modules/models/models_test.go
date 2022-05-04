package models

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SettingUp("test")
	exitVal := m.Run()
	TearDown("test")
	os.Exit(exitVal)
}

type Creator interface {
	Create() (any, error)
}

func CreateRecords(records ...Creator) {
	for _, v := range records {
		v.Create()
	}
}
