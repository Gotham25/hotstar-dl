package utils

import (
	"encoding/json"
	"fmt"
)

//GetMasterPlaybackUrl gets master playback url from playback uri page contents.
func GetMasterPlaybackUrl(playbackUriPageContents []byte) (string, error) {

	var masterPlaybackUrl string
	var result map[string]interface{}
	json.Unmarshal(playbackUriPageContents, &result)

	statusCode := int(result["statusCodeValue"].(float64))

	if statusCode == 200 {
		body := result["body"].(map[string]interface{})
		results := body["results"].(map[string]interface{})
		item := results["item"].(map[string]interface{})
		masterPlaybackUrl = item["playbackUrl"].(string)
		return masterPlaybackUrl, nil
	}

	return "", fmt.Errorf("Invalid status code %d", statusCode)
}
