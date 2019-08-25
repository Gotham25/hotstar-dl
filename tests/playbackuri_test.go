package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestGetPlaybackURI_ValidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/validTestPageContents.html")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackURI := "https://api.hotstar.com/h/v2/play/in/contents/1100025368?desiredConfig=encryption:plain;ladder:phone,tv;package:hls,dash&client=mweb&clientVersion=6.18.0&deviceId=79272307-fa98-4b08-8f5c-5afdde2687ff&osName=Windows&osVersion=10"

	actualPlaybackURI, _, err := utils.GetPlaybackURI(fmt.Sprintf("%s", testPageContents), "https://www.hotstar.com/tv/naam-iruvar-namaku-iruvar/s-1446/mayan-confronts-aravind/1100025368", "1100025368", "79272307-fa98-4b08-8f5c-5afdde2687ff")

	if err == nil && expectedPlaybackURI != actualPlaybackURI {
		t.Error("Expected", expectedPlaybackURI, " but got", actualPlaybackURI)
	}
}

func TestGetPlaybackURI_InvalidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/invalidTestPageContents.html")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackURIError := "Invalid appState JSON. Cannot retrieve playbackUri"

	_, _, actualPlaybackURIError := utils.GetPlaybackURI(fmt.Sprintf("%s", testPageContents), "https://www.hotstar.com/tv/naam-iruvar-namaku-iruvar/s-1446/mayan-confronts-aravind/1100025368", "1100025368", "79272307-fa98-4b08-8f5c-5afdde2687ff")

	if actualPlaybackURIError != nil && expectedPlaybackURIError != actualPlaybackURIError.Error() {
		t.Error("Expected", expectedPlaybackURIError, " but got", actualPlaybackURIError)
	}
}
