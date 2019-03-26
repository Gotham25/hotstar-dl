package tests

import (
	"github.com/Gotham25/hotstar-dl/utils"
	"reflect"
	"testing"
)

func TestCopyMap(t *testing.T) {
	originalMap := map[string]string{
		"Name":        "Scott",
		"Designation": "Software Evangelist",
	}

	clonedMap := utils.CopyMap(originalMap)

	if !reflect.DeepEqual(originalMap, clonedMap) {
		t.Error("Expected", originalMap, " but got", clonedMap)
	}
}

func TestGetDateStr_NonPadded(t *testing.T) {

	expectedDate := "2018-04-27 21:30:00 +0530 IST"
	actualDate := utils.GetDateStr(1524844800)

	if expectedDate != actualDate {
		t.Error("Expected", expectedDate, " but got", actualDate)
	}

}

func TestGetDateStr_Padded(t *testing.T) {

	expectedDate := "2018-04-27 17:45:00 +0530 IST"
	actualDate := utils.GetDateStr(1524831300000)

	if expectedDate != actualDate {
		t.Error("Expected", expectedDate, " but got", actualDate)
	}

}
