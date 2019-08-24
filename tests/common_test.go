package tests

import (
	"reflect"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
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

func TestMakeRange(t *testing.T) {
	expectedRange := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	actualRange := utils.MakeRange(1, 13)

	if !reflect.DeepEqual(expectedRange, actualRange) {
		t.Error("Expected", expectedRange, " but got", actualRange)
	}
}
