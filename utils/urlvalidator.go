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

//IsValidHotstarUrl validates if the given video url is a valid Hotstar url or not.
func IsValidHotstarUrl(videoUrl string) (bool, string) {
	var urlRegex = regexp.MustCompile(`((https|http)?://)?(www\.)?hotstar\.com/(?:.+?[/-])+(?P<videoId>\d{10})`)
	if urlRegex.MatchString(videoUrl) {
		match := reSubMatchMap(urlRegex, videoUrl)
		return true, match["videoId"]
	}

	return false, ""
}

//GetParsedVideoUrl parses given video url for proper url scheme.
func GetParsedVideoUrl(videoUrl string) string {
	parsedUrl, err := url.Parse(videoUrl)

	if err != nil {
		log.Fatal(err)
	}

	switch parsedUrl.Scheme {
	case "":
		//fmt.Println("Replacing empty url scheme with https")
		parsedUrl.Scheme = "https"
	case "https":
		//do nothing
	case "http":
		//fmt.Println("Replacing http url scheme with https")
		parsedUrl.Scheme = "https"
	default:
		fmt.Println("Invalid url scheme please enter valid one")
		os.Exit(-1)
	}

	videoUrl = fmt.Sprintf("%v", parsedUrl)

	fmt.Println("Parsed video url is", parsedUrl)

	return videoUrl
}
