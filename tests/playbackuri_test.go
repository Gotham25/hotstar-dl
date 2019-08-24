package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestGetPlaybackUri_ValidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/validTestPageContents.html")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackUri := "https://api.hotstar.com/h/v2/play/in/contents/1100025368?desiredConfig=encryption:plain;ladder:phone,tv;package:hls,dash&client=mweb&clientVersion=6.18.0&deviceId=79272307-fa98-4b08-8f5c-5afdde2687ff&osName=Windows&osVersion=10"

	actualPlaybackUri, _, err := utils.GetPlaybackUri(fmt.Sprintf("%s", testPageContents), "https://www.hotstar.com/tv/naam-iruvar-namaku-iruvar/s-1446/mayan-confronts-aravind/1100025368", "1100025368", "79272307-fa98-4b08-8f5c-5afdde2687ff")

	if err == nil && expectedPlaybackUri != actualPlaybackUri {
		t.Error("Expected", expectedPlaybackUri, " but got", actualPlaybackUri)
	}
}

func TestGetPlaybackUri_InvalidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/invalidTestPageContents.html")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackUriError := "Invalid appState JSON. Cannot retrieve playbackUri"

	_, _, actualPlaybackUriError := utils.GetPlaybackUri(fmt.Sprintf("%s", testPageContents), "https://www.hotstar.com/tv/naam-iruvar-namaku-iruvar/s-1446/mayan-confronts-aravind/1100025368", "1100025368", "79272307-fa98-4b08-8f5c-5afdde2687ff")

	if actualPlaybackUriError != nil && expectedPlaybackUriError != actualPlaybackUriError.Error() {
		t.Error("Expected", expectedPlaybackUriError, " but got", actualPlaybackUriError)
	}
}
