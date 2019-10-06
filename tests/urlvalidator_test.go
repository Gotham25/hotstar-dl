package tests

import (
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestIsValidHotstarVideoURL_ValidURL1(t *testing.T) {
	expectedVideoID := "1100003795"

	isValid, actualVideoID, isPlaylistURL := utils.IsValidHotstarURL("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/1100003795")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}

	if isPlaylistURL != false {
		t.Error("Expected", !isPlaylistURL, "but got", isPlaylistURL)
	}
}

func TestIsValidHotstarVideoURL_ValidURL2(t *testing.T) {
	expectedVideoID := "1100020335"

	isValid, actualVideoID, isPlaylistURL := utils.IsValidHotstarURL("https://uk.hotstar.com/tv/vijay-television-awards/s-203/the-main-event/1100020335")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}

	if isPlaylistURL != false {
		t.Error("Expected", !isPlaylistURL, "but got", isPlaylistURL)
	}
}

func TestIsValidHotstarVideoURLOrPlaylist_InvalidURL1(t *testing.T) {
	expectedVideoID := ""

	isValid, actualVideoID, _ := utils.IsValidHotstarURL("http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/123")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}

func TestIsValidHotstarVideoURLOrPlaylist_InvalidURL2(t *testing.T) {
	expectedVideoID := ""

	isValid, actualVideoID, _ := utils.IsValidHotstarURL("http://www.google.com")

	if isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedVideoID != actualVideoID {
		t.Error("Expected", expectedVideoID, "but got", actualVideoID)
	}
}

func TestIsValidHotstarPlaylistURL_ValidURL(t *testing.T) {
	expectedPlaylistID := "1_2_2213"

	isValid, actualPlaylistID, isPlaylistURL := utils.IsValidHotstarURL("https://www.hotstar.com/tv/ayudha-ezhuthu/s-2213/list/episodes/t-1_2_2213")

	if !isValid {
		t.Error("Expected", !isValid, "but got", isValid)
	}

	if expectedPlaylistID != actualPlaylistID {
		t.Error("Expected", expectedPlaylistID, "but got", actualPlaylistID)
	}

	if isPlaylistURL != true {
		t.Error("Expected", !isPlaylistURL, "but got", isPlaylistURL)
	}
}
