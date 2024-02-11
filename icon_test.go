package mcstatus_test

import (
	"testing"

	"github.com/mcstatus-io/go-mcstatus"
)

func TestGetIcon(t *testing.T) {
	result, err := mcstatus.GetIcon("demo.mcstatus.io", 25565)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", result)
}
