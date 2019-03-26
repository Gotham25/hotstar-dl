package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var playbackUriRetryCount = 0

func populateMetaDataMapWithMetadata(metaDataMap map[string]string, metadata map[string]interface{}) {
	for k1, v1 := range metadata {
		switch k1 {
		case "title":
			title := v1.(string)
			metaDataMap[k1] = title
			metaDataMap["album"] = title
		case "broadcastDate":
			metaDataMap["date"] = GetDateStr(v1.(float64))
		case "channelName":
			metaDataMap["copyright"] = v1.(string)
		case "drmProtected":
			metaDataMap["drmProtected"] = fmt.Sprintf("%v", v1)
		case "actors":
			actors := ""
			for _, actor := range v1.([]interface{}) {
				if len(actors) != 0 {
					actors += ",\n"
				}
				actors += actor.(string)
			}
			metaDataMap["artist"] = actors
			metaDataMap["album_artist"] = actors
		case "description":
			metaDataMap["comment"] = v1.(string)
			metaDataMap["synopsis"] = v1.(string)
		case "genre":
			metaDataMap[k1] = v1.(string)
		case "showName":
			metaDataMap["show"] = v1.(string)
		case "episodeNo":
			metaDataMap["episode_id"] = fmt.Sprintf("%d", int64(v1.(float64)))
		case "seasonNo":
			metaDataMap["season_number"] = fmt.Sprintf("%d", int64(v1.(float64)))
		default:
			//do nothing
		}
	}
}

//GetPlaybackUri gets the playback uri from the metadata in the given page contents.
func GetPlaybackUri(videoUrlPageContents string, videoUrl string, videoId string) (string, map[string]string, error) {
	//TODO: show retry info upon debug level

	var metadata = make(map[string]interface{})
	appStateSearchRegex := *regexp.MustCompile(`<script>window.APP_STATE=(.+?)</script>`)
	appStateSearchMatch := appStateSearchRegex.FindAllStringSubmatch(videoUrlPageContents, -1)

	if len(appStateSearchMatch) > 0 {

		var result map[string]interface{}
		json.Unmarshal([]byte(appStateSearchMatch[0][1]), &result)

		for k := range result {

			//videoId := helper.After(k, "/")
			if len(videoId) != 0 && strings.Contains(videoUrl, videoId) {

				root, isCastOk := result[k].(map[string]interface{})

				if !isCastOk && (playbackUriRetryCount+1 < 5) {
					playbackUriRetryCount++
					//fmt.Printf("GetPlaybackUri: cast to map[string]interface{} failed. retrying count : #%d\n", playbackUriRetryCount)
					return GetPlaybackUri(videoUrlPageContents, videoUrl, videoId)
				}

				initialState := root["initialState"].(map[string]interface{})
				contentData := initialState["contentData"].(map[string]interface{})
				content := contentData["content"].(map[string]interface{})

				for contentKey, contentValue := range content {
					metadata[contentKey] = contentValue
				}
				//quit looping as we have got the required metdata
				break
			}

		}

	} else {
		return "", nil, errors.New("Invalid appState JSON. Cannot retrieve playbackUri")
	}

	if playbackUri, ok := metadata["playbackUri"].(string); ok {

		metaDataMap := make(map[string]string)
		populateMetaDataMapWithMetadata(metaDataMap, metadata)
		return playbackUri, metaDataMap, nil

	}

	return "", nil, errors.New("Error msg : Cannot retrieve playbackUri")

}
