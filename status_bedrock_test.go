package mcstatus_test

import (
	"testing"

	"github.com/mcstatus-io/go-mcstatus"
)

func TestGetBedrockStatus(t *testing.T) {
	result, err := mcstatus.GetBedrockStatus("demo.mcstatus.io", 19132)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", result)
}
