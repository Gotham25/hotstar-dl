package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

//GetMasterPlaybackURLs gets master playback urls from playback uri page contents.
func GetMasterPlaybackURLs(playbackURIPageContents []byte) ([]string, error) {
	var result map[string]interface{}
	json.Unmarshal(playbackURIPageContents, &result)

	message := result["message"].(string)

	if strings.Contains(message, "success") {
		data := result["data"].(map[string]interface{})
		playbackSets := data["playBackSets"].([]interface{})
		masterPlaybackUrls := make([]string, 0)
		for _, v := range playbackSets {
			playbackSet := v.(map[string]interface{})
			masterPlaybackUrls = append(masterPlaybackUrls, playbackSet["playbackUrl"].(string))
		}
		return masterPlaybackUrls, nil
	}

	return make([]string, 0), fmt.Errorf("Error: %s", message)
}

/*
//GetMasterPlaybackURLs gets master playback urls from playback uri page contents.
func GetMasterPlaybackURLs(playbackURIPageContents []byte) ([]string, error) {

	var result map[string]interface{}
	json.Unmarshal(playbackURIPageContents, &result)

	statusCode := int(result["statusCodeValue"].(float64))

	if statusCode == 200 {
		body := result["body"].(map[string]interface{})
		results := body["results"].(map[string]interface{})
		playbackSets := results["playBackSets"].([]interface{})
		masterPlaybackUrls := make([]string, 0)
		for _, v := range playbackSets {
			playbackSet := v.(map[string]interface{})
			masterPlaybackUrls = append(masterPlaybackUrls, playbackSet["playbackUrl"].(string))
		}
		return masterPlaybackUrls, nil
	}

	return make([]string, 0), fmt.Errorf("Invalid status code %d", statusCode)
}
*/
