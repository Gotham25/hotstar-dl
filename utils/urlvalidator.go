package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
)

//IsValidHotstarURL validates if the given video url is a valid Hotstar url or not.
func IsValidHotstarURL(videoOrPlaylistURL string) (bool, string, bool) {
	var videoURLRegex = regexp.MustCompile(`(https?://)?(www|uk\.)?hotstar\.com/(?:.+?[/-])+(?P<videoId>\d{10})`)
	var playlistURLRegex = regexp.MustCompile(`(https?://)?(www|uk\.)?hotstar\.com/tv/[^/]+/s-\w+/list/[^/]+/t-(?P<playlistId>\w+)`)

	if videoURLRegex.MatchString(videoOrPlaylistURL) {
		match := ReSubMatchMap(videoURLRegex, videoOrPlaylistURL)
		return true, match["videoId"], false
	} else if playlistURLRegex.MatchString(videoOrPlaylistURL) {
		match := ReSubMatchMap(playlistURLRegex, videoOrPlaylistURL)
		return true, match["playlistId"], true
	}

	return false, "", false
}

//GetParsedVideoURL parses given video url for proper url scheme.
func GetParsedVideoURL(videoURL string) string {
	parsedURL, err := url.Parse(videoURL)

	if err != nil {
		log.Fatal(err)
	}

	switch parsedURL.Scheme {
	case "":
		//fmt.Println("Replacing empty url scheme with https")
		parsedURL.Scheme = "https"
	case "https":
		//do nothing
	case "http":
		//fmt.Println("Replacing http url scheme with https")
		parsedURL.Scheme = "https"
	default:
		fmt.Println("Invalid url scheme please enter valid one")
		os.Exit(-1)
	}

	videoURL = fmt.Sprintf("%v", parsedURL)

	fmt.Println("Parsed video url is", parsedURL)

	return videoURL
}
