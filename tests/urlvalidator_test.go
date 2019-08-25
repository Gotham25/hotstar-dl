package tests

import (
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestIsValidHotstarURL_ValidURL1(t *testing.T) {
	expectedVideoID := "1100003795"

	isValid, actualVideoID := utils.IsValidHotstarURL("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/1100003795")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}

func TestIsValidHotstarURL_ValidURL2(t *testing.T) {
	expectedVideoID := "1100020335"

	isValid, actualVideoID := utils.IsValidHotstarURL("https://uk.hotstar.com/tv/vijay-television-awards/s-203/the-main-event/1100020335")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}

func TestIsValidHotstarURL_InvalidURL1(t *testing.T) {
	expectedVideoID := ""

	isValid, actualVideoID := utils.IsValidHotstarURL("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/123")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}

func TestIsValidHotstarURL_InvalidURL2(t *testing.T) {
	expectedVideoID := ""

	isValid, actualVideoID := utils.IsValidHotstarURL("http://www.google.com")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}
