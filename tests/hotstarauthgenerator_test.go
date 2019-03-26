package tests

import (
	"github.com/Gotham25/hotstar-dl/utils"
	"testing"
)

func TestGenerate(t *testing.T) {
	expectedHotstarAuth := "st=1551575749~exp=1551581749~acl=/*~hmac=c29d490c8c13024b8e61c608ff1baa1899a092fff9433d1112f367b5dd9a334d"
	actualHotstarAuth := utils.Generate(1551575749)

	if expectedHotstarAuth != actualHotstarAuth {
		t.Error("Expected", expectedHotstarAuth, " but got", actualHotstarAuth)
	}
}
