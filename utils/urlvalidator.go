package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
)

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}

//IsValidHotstarURL validates if the given video url is a valid Hotstar url or not.
func IsValidHotstarURL(videoURL string) (bool, string) {
	var urlRegex = regexp.MustCompile(`(https|http?://)?(www|uk\.)?hotstar\.com/(?:.+?[/-])+(?P<videoId>\d{10})`)
	if urlRegex.MatchString(videoURL) {
		match := reSubMatchMap(urlRegex, videoURL)
		return true, match["videoId"]
	}

	return false, ""
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
