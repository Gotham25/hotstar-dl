package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var refreshTokenRetryCount = 0
var videoURLRetryCount = 0
var playbackRetryCount = 0
var playbackURIContentRetryCount = 0
var masterPlaybackURLRetryCount = 0
var tempVideoFormatsRetryCount = 0

func getRequestHeaders() map[string]string {
	return map[string]string{
		"Hotstarauth":     GenerateHotstarAuth(),
		"X-Country-Code":  "IN",
		"X-Platform-Code": "JIO",
		"X-HS-Platform":   "web",
		"X-HS-AppVersion": "6.72.2",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.122 Safari/537.36",
		"Origin":          "https://www.hotstar.com",
	}
}

func getRefreshTokenHeaders() map[string]string {
	return map[string]string{
		"hotstarauth":  GenerateHotstarAuth(),
		"userIdentity": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1bV9hY2Nlc3MiLCJleHAiOjE1OTAzMDE5OTYsImlhdCI6MTU4OTY5NzE5NiwiaXNzIjoiVFMiLCJzdWIiOiJ7XCJoSWRcIjpcIjU5YjkxMTBkMDM2ZjQ3M2U5OTBhNzFiNDAwMDM5MzRkXCIsXCJwSWRcIjpcImVjNmRmY2Y1ZTJhYzRhYWJhZjNjOTBlY2I0YTY5MTVlXCIsXCJuYW1lXCI6XCJHdWVzdCBVc2VyXCIsXCJpcFwiOlwiMjIzLjIyNi4yOS4yMjdcIixcImNvdW50cnlDb2RlXCI6XCJpblwiLFwiY3VzdG9tZXJUeXBlXCI6XCJudVwiLFwidHlwZVwiOlwiZGV2aWNlXCIsXCJpc0VtYWlsVmVyaWZpZWRcIjpmYWxzZSxcImlzUGhvbmVWZXJpZmllZFwiOmZhbHNlLFwiZGV2aWNlSWRcIjpcImExMTg1MTFhLTJmYjktNDhmOS04MGM5LWY1OTlkMjdlYTZmNlwiLFwicHJvZmlsZVwiOlwiQURVTFRcIixcInZlcnNpb25cIjpcInYyXCIsXCJzdWJzY3JpcHRpb25zXCI6e1wiaW5cIjp7fX0sXCJpc3N1ZWRBdFwiOjE1ODk2OTcxOTY2NzJ9IiwidmVyc2lvbiI6IjFfMCJ9.bkx7DodQSFohwmzqf8RmKOr3tuORgVFEh_qbtdqzeVA",
		"Origin":       "https://www.hotstar.com",
		"deviceId":     uuid.New().String(),
	}
}

func getAggregatedFormats(videoFormatsTemp, audioDashFormatsTemp, videoDashFormatsTemp map[string][]map[string]string) map[string]map[string]string {
	totalFormats := make(map[string]map[string]string)

	for fid, formatsList := range videoFormatsTemp {
		if len(formatsList) == 1 {
			totalFormats[fid] = formatsList[0]
		} else {
			for index, formats := range formatsList {
				totalFormats[fmt.Sprintf("%s-%d", fid, index)] = formats
			}
		}
	}

	for _, value := range videoDashFormatsTemp {
		baseFormat := "dash-video"
		if len(value) == 1 {
			totalFormats[fmt.Sprintf("%s-%s", baseFormat, value[0]["K-FORM-NUMBER"])] = value[0]
		} else {
			for index, value2 := range value {
				totalFormats[fmt.Sprintf("%s-%s-%d", baseFormat, value2["K-FORM-NUMBER"], index)] = value2
			}
		}
	}

	for _, value := range audioDashFormatsTemp {
		baseFormat := "dash-audio"
		if len(value) == 1 {
			totalFormats[fmt.Sprintf("%s-%s", baseFormat, value[0]["K-FORM-NUMBER"])] = value[0]
		} else {
			for index, value2 := range value {
				totalFormats[fmt.Sprintf("%s-%s-%d", baseFormat, value2["K-FORM-NUMBER"], index)] = value2
			}
		}
	}

	return totalFormats
}

func getRefreshToken(videoURL string) (string, error) {
	var result map[string]interface{}
	var refreshTokenHeaders = getRefreshTokenHeaders()
	refreshTokenHeaders["Referer"] = videoURL

	refreshTokenURL := "https://api.hotstar.com/in/aadhar/v2/web/in/user/refresh-token"
	refreshTokenURLContentBytes, err := MakeGetRequest(refreshTokenURL, refreshTokenHeaders)

	if err != nil {
		if refreshTokenRetryCount < 10 {

			//retry again for fetching JWT token
			refreshTokenRetryCount++

			//fmt.Printf("getRefreshToken: GET request to videoURL failed... Retrying count : #%d\n", refreshTokenRetryCount)
			return getRefreshToken(videoURL)
		}
	}

	json.Unmarshal(refreshTokenURLContentBytes, &result)

	description := result["description"].(map[string]interface{})
	userIdentity := description["userIdentity"].(string)

	return userIdentity, nil
}

func getVideoURL(videoURL string, requestHeaders map[string]string) (string, error) {
	videoURLContentBytes, err := MakeGetRequest(videoURL, requestHeaders)

	if err != nil {
		if videoURLRetryCount+1 < 10 {
			//retry again for fetching formats
			videoURLRetryCount++
			//fmt.Printf("GetVideoFormats: GET request to videoURL failed... Retrying count : #%d\n", videoURLRetryCount)
			return getVideoURL(videoURL, requestHeaders)
		}
		return "", err
	}

	return fmt.Sprintf("%s", videoURLContentBytes), nil
}

func getPlayback(videoURLContent, videoURL, videoID, uuid string) (string, map[string]string, error) {
	playbackURI, videoMetadata, err := GetPlaybackURI(videoURLContent, videoURL, videoID, uuid)

	if err != nil {
		if playbackRetryCount+1 < 10 {
			//retry again for fetching formats
			playbackRetryCount++
			//fmt.Printf("GetVideoFormats: Invalid APP_STATE json... retrying count : #%d\n", playbackRetryCount)
			return getPlayback(videoURLContent, videoURL, videoID, uuid)
		}
		return "", nil, err
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	return playbackURI, videoMetadata, nil
}

func getPlaybackURIContent(playbackURI string, requestHeaders map[string]string) ([]byte, error) {
	playbackURIContentBytes, err := MakeGetRequest(playbackURI, requestHeaders)

	if err != nil {
		if playbackURIContentRetryCount+1 < 10 {
			//retry again for fetching formats
			playbackURIContentRetryCount++
			//fmt.Printf("GetVideoFormats: GET request to playbackURI failed... Retrying count : #%d\n", playbackURIContentRetryCount)
			return getPlaybackURIContent(playbackURI, requestHeaders)
		}
		return nil, err
	}

	return playbackURIContentBytes, nil
}

func getMasterPlaybackURL(playbackURIContentBytes []byte) ([]string, error) {
	masterPlaybackURLs, err := GetMasterPlaybackURLs(playbackURIContentBytes)

	if err != nil {
		if masterPlaybackURLRetryCount+1 < 10 {
			//retry again for fetching formats
			masterPlaybackURLRetryCount++
			//fmt.Printf("GetVideoFormats: Retriving masterPlaybackURL failed... Retrying count : #%d\n", masterPlaybackURLRetryCount)
			return getMasterPlaybackURL(playbackURIContentBytes)
		}
		return nil, err
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}
	return masterPlaybackURLs, nil
}

func getTempVideoFormats(masterPlaybackURLs []string, requestHeaders map[string]string) (map[string][]map[string]string, map[string][]map[string]string, map[string][]map[string]string, error) {
	videoFormatsTemp := make(map[string][]map[string]string)
	videoDashFormatsTemp := make(map[string][]map[string]string)
	audioDashFormatsTemp := make(map[string][]map[string]string)

	for _, masterPlaybackURL := range masterPlaybackURLs {

		if masterPlaybackURL != "" {

			var queryParams string
			masterPlaybackURLQueryParam := strings.Split(masterPlaybackURL, "?")

			if len(masterPlaybackURLQueryParam) > 1 {
				queryParams = masterPlaybackURLQueryParam[1]
			}

			if strings.Contains(masterPlaybackURL, "m3u8") {

				masterPlaybackPageContentsM3u8Bytes, err := MakeGetRequest(masterPlaybackURL, requestHeaders)

				if err != nil {

					if tempVideoFormatsRetryCount+1 < 10 {
						//retry again for fetching formats
						tempVideoFormatsRetryCount++
						//fmt.Printf("GetVideoFormats: GET request to masterPlaybackURL failed... Retrying count : #%d\n", tempVideoFormatsRetryCount)
						return getTempVideoFormats(masterPlaybackURLs, requestHeaders)
					}

					return nil, nil, nil, err
				}

				for fid, formatsList := range ParseM3u8Content(fmt.Sprintf("%s", masterPlaybackPageContentsM3u8Bytes), masterPlaybackURL, queryParams) {
					videoFormatsTemp[fid] = append(videoFormatsTemp[fid], formatsList)
				}
			} else {

				masterPlaybackPageContentsMpdBytes, err := MakeGetRequest(masterPlaybackURL, requestHeaders)

				if err != nil {

					if tempVideoFormatsRetryCount+1 < 10 {
						//retry again for fetching formats
						tempVideoFormatsRetryCount++
						//fmt.Printf("GetVideoFormats: GET request to masterPlaybackURL failed... Retrying count : #%d\n", tempVideoFormatsRetryCount)
						return getTempVideoFormats(masterPlaybackURLs, requestHeaders)
					}

					return nil, nil, nil, err
				}

				dFormats := GetDashFormats(masterPlaybackPageContentsMpdBytes, masterPlaybackURL)

				for avType, formatsList := range dFormats {
					for formatCode, formatInfo := range formatsList {
						if avType == "video" {
							videoDashFormatsTemp[formatCode] = append(videoDashFormatsTemp[formatCode], formatInfo)
						} else {
							audioDashFormatsTemp[formatCode] = append(audioDashFormatsTemp[formatCode], formatInfo)
						}
					}
				}

			}

		}

	}

	return videoFormatsTemp, audioDashFormatsTemp, videoDashFormatsTemp, nil
}

func raiseError(errorMsg string) {
	fmt.Println(errorMsg)
	os.Exit(-1)
}

func raiseConversionError(varname, varvalue string) {
	raiseError(fmt.Sprintf("\nError in converting %s, %s to integer", varname, varvalue))
}

func isValidPlaylistBounds(playlistItemCount int, playlistStartRange, playlistEndRange string) (int, int, string, bool) {

	var validationMessage strings.Builder
	isValid := false
	startRange, startRangeError := strconv.Atoi(playlistStartRange)
	if startRangeError != nil {
		raiseConversionError("playlistStartRange", playlistStartRange)
	}
	endRange, endRangeError := strconv.Atoi(playlistEndRange)
	if endRangeError != nil {
		raiseConversionError("playlistEndRange", playlistEndRange)
	}

	validationMessage.WriteString("")

	if startRange < 1 {
		validationMessage.WriteString(fmt.Sprintf("\nInvalid start range %s provided. Should be >= 1", playlistStartRange))
	}

	if endRange > playlistItemCount {
		validationMessage.WriteString(fmt.Sprintf("\nInvalid end range %s provided. Should be <= %d", playlistEndRange, playlistItemCount))
	}

	if startRange > endRange {
		validationMessage.WriteString(fmt.Sprintf("\nInvalid start range %d provided. Should be <= %d", startRange, endRange))
	} else if endRange < startRange {
		validationMessage.WriteString(fmt.Sprintf("\nInvalid end range %d provided. Should be >= %d", endRange, startRange))
	} else {
		isValid = true
	}

	return startRange, endRange, validationMessage.String(), isValid
}
