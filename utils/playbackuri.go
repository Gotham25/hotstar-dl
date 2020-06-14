package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var playbackURIRetryCount = 0

//ACTOR a metadata constant for actor entry
const ACTOR = "actor"

//GENRE a metadata constant for genre entry
const GENRE = "genre"

//CONTENTID a metadata constant for content id entry
const CONTENTID = "contentId"

func getMetadata(metaValue interface{}, metaDataType string) string {
	switch metaDataType {
	case ACTOR:
		actors := ""
		for _, actor := range metaValue.([]interface{}) {
			if len(actors) != 0 {
				actors += ",\n"
			}
			actors += actor.(string)
		}
		return actors

	case CONTENTID:
		if strings.EqualFold(fmt.Sprintf("%T", metaValue), "string") {
			return metaValue.(string)
		}
		return fmt.Sprintf("%d", int(metaValue.(float64)))

	case GENRE:
		genres := ""

		if strings.EqualFold(fmt.Sprintf("%T", metaValue), "string") {
			return metaValue.(string)
		}

		for _, genre := range metaValue.([]interface{}) {
			if len(genres) != 0 {
				genres += ",\n"
			}
			genres += genre.(string)
		}

		return genres

	default:
		return "Invalid metadata type"

	}
}

//PopulateMetaDataMapWithMetadata poulates metadata map with required meta-data
func PopulateMetaDataMapWithMetadata(metaDataMap map[string]string, metadata map[string]interface{}) {
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
			actors := getMetadata(v1, k1)
			metaDataMap["artist"] = actors
			metaDataMap["album_artist"] = actors
		case "description":
			metaDataMap["comment"] = v1.(string)
			metaDataMap["synopsis"] = v1.(string)
		case "genre":
			metaDataMap[k1] = getMetadata(v1, k1)
		case "showName":
			metaDataMap["show"] = v1.(string)
		case "episodeNo":
			metaDataMap["episode_id"] = fmt.Sprintf("%d", int64(v1.(float64)))
		case "seasonNo":
			metaDataMap["season_number"] = fmt.Sprintf("%d", int64(v1.(float64)))
		case "contentId":
			metaDataMap["id"] = getMetadata(v1, k1)

		default:
			//do nothing
		}
	}
}

//GetPlaybackURI3 gets the playback uri v2 from videoID
func GetPlaybackURI3(videoID string, uuid string) string {
	var playbackURI3 strings.Builder
	playbackURI3.WriteString(fmt.Sprintf("https://api.hotstar.com/play/v1/playback/content/%s?", videoID))
	playbackURI3.WriteString(fmt.Sprintf("%s=%s&", "device-id", uuid))
	playbackURI3.WriteString(fmt.Sprintf("%s=%s&", "desired-config", "encryption:plain;ladder:phone,tv;package:hls,dash"))
	playbackURI3.WriteString(fmt.Sprintf("%s=%s&", "os-name", "Windows"))
	playbackURI3.WriteString(fmt.Sprintf("%s=%s", "os-version", "10"))
	return playbackURI3.String()
}

//GetPlaybackURI2 gets the playback uri v2 from videoID
func GetPlaybackURI2(videoID string, uuid string) string {
	var playbackURI2 strings.Builder
	playbackURI2.WriteString(fmt.Sprintf("https://api.hotstar.com/h/v2/play/in/contents/%s?", videoID))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "desiredConfig", "encryption:plain;ladder:phone,tv;package:hls,dash"))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "client", "mweb"))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "clientVersion", "6.18.0"))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "deviceId", uuid))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "osName", "Windows"))
	playbackURI2.WriteString(fmt.Sprintf("%s=%s", "osVersion", "10"))
	return playbackURI2.String()
}

//GetPlaybackURI gets the playback uri from the metadata in the given page contents.
func GetPlaybackURI(videoURLPageContents string, videoURL string, videoID string, uuid string) (string, map[string]string, error) {
	//TODO: show retry info upon debug level

	var metadata = make(map[string]interface{})
	appStateSearchRegex := *regexp.MustCompile(`<script>window.APP_STATE=(.+?)</script>`)
	appStateSearchMatch := appStateSearchRegex.FindAllStringSubmatch(videoURLPageContents, -1)

	if len(appStateSearchMatch) > 0 {

		var result map[string]interface{}
		json.Unmarshal([]byte(appStateSearchMatch[0][1]), &result)

		for resultKey, resultValue := range result {

			//videoID := helper.After(k, "/")
			if len(videoID) != 0 && strings.Contains(resultKey, videoID) {

				root, isCastOk := resultValue.(map[string]interface{})

				if !isCastOk && (playbackURIRetryCount+1 < 5) {
					playbackURIRetryCount++
					//fmt.Printf("GetPlaybackURI: cast to map[string]interface{} failed. retrying count : #%d\n", playbackURIRetryCount)
					return GetPlaybackURI(videoURLPageContents, videoURL, videoID, uuid)
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

	if playbackURI, ok := metadata["playbackUri"].(string); ok {

		metaDataMap := make(map[string]string)
		PopulateMetaDataMapWithMetadata(metaDataMap, metadata)
		metaDataMap["playbackUri"] = playbackURI
		/* var playbackURI2 strings.Builder

		playbackURI2.WriteString(fmt.Sprintf("https://api.hotstar.com/h/v2/play/in/contents/%s?", videoID))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "desiredConfig", "encryption:plain;ladder:phone,tv;package:hls,dash"))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "client", "mweb"))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "clientVersion", "6.18.0"))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "deviceId", uuid))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s&", "osName", "Windows"))
		playbackURI2.WriteString(fmt.Sprintf("%s=%s", "osVersion", "10")) */

		// return /* playbackURI2.String() */ GetPlaybackURI2(videoID, uuid), metaDataMap, nil

		return /* playbackURI2.String() */ GetPlaybackURI3(videoID, uuid), metaDataMap, nil

	}

	return "", nil, errors.New("Error msg : Cannot retrieve playbackUri")

}
