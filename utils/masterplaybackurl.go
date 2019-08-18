package utils

import (
	"encoding/json"
	"fmt"
)

//GetMasterPlaybackUrl gets master playback url from playback uri page contents.
func GetMasterPlaybackUrls(playbackUriPageContents []byte) ([]string, error) {

	//var masterPlaybackUrl string
	var result map[string]interface{}
	json.Unmarshal(playbackUriPageContents, &result)

	statusCode := int(result["statusCodeValue"].(float64))

	if statusCode == 200 {
		body := result["body"].(map[string]interface{})
		results := body["results"].(map[string]interface{})
		//item := results["item"].(map[string]interface{})
		//masterPlaybackUrl = item["playbackUrl"].(string)
		//return masterPlaybackUrl, nil
		playbackSets := results["playBackSets"].([]interface{})
		masterPlaybackUrls := make([]string, len(playbackSets))
		for _, v := range playbackSets {
			playbackSet := v.(map[string]interface{})
			//fmt.Println("playbackSet: ", playbackSet["playbackUrl"].(string))
			masterPlaybackUrls = append(masterPlaybackUrls, playbackSet["playbackUrl"].(string))
		}
		//fmt.Println("\nmasterPlaybackUrls: ", masterPlaybackUrls, "\n")
		//fmt.Println("\nplaybackSets size : ", len(playbackSets), "\n")
		//playbackSet1 := playbackSets[0].(map[string]interface{})
		//return playbackSet1["playbackUrl"].(string), nil
		return masterPlaybackUrls, nil
	}

	return make([]string, 0), fmt.Errorf("Invalid status code %d", statusCode)
}
