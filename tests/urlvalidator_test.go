package tests

import (
	"github.com/Gotham25/hotstar-dl/utils"
	"testing"
)

func TestIsValidHotstarUrl_ValidUrl(t *testing.T) {
	expectedVideoId := "1100003795"

	isValid, actualVideoId := utils.IsValidHotstarUrl("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/1100003795")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoId != actualVideoId {
		t.Error("Expected", expectedVideoId, "but got", actualVideoId)
	}
}

func TestIsValidHotstarUrl_InvalidUrl1(t *testing.T) {
	expectedVideoId := ""

	isValid, actualVideoId := utils.IsValidHotstarUrl("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/123")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoId != actualVideoId {
		t.Error("Expected", expectedVideoId, "but got", actualVideoId)
	}
}

func TestIsValidHotstarUrl_InvalidUrl2(t *testing.T) {
	expectedVideoId := ""

	isValid, actualVideoId := utils.IsValidHotstarUrl("http://www.google.com")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoId != actualVideoId {
		t.Error("Expected", expectedVideoId, "but got", actualVideoId)
	}
}
