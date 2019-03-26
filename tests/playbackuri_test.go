package tests

import (
	"fmt"
	"github.com/Gotham25/hotstar-dl/utils"
	"io/ioutil"
	"log"
	"testing"
)

func TestGetPlaybackUri_ValidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/validTestPageContents.txt")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackUri := "https://api.hotstar.com/h/v1/play?contentId=1100003795"

	actualPlaybackUri, _, err := utils.GetPlaybackUri(fmt.Sprintf("%s", testPageContents), "http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/1100003795", "1100003795")

	if err == nil && expectedPlaybackUri != actualPlaybackUri {
		t.Error("Expected", expectedPlaybackUri, " but got", actualPlaybackUri)
	}
}

func TestGetPlaybackUri_InvalidPageContents(t *testing.T) {
	testPageContents, err := ioutil.ReadFile("resources/invalidTestPageContents.txt")
	if err != nil {
		log.Fatal(err)
	}

	expectedPlaybackUriError := "Invalid appState JSON. Cannot retrieve playbackUri"

	_, _, actualPlaybackUriError := utils.GetPlaybackUri(fmt.Sprintf("%s", testPageContents), "http://www.hotstar.com/tv/chinnathambi/15301/chinnathambi-yearns-for-nandini/1100003795", "1100003795")

	if actualPlaybackUriError != nil && expectedPlaybackUriError != actualPlaybackUriError.Error() {
		t.Error("Expected", expectedPlaybackUriError, " but got", actualPlaybackUriError)
	}
}
