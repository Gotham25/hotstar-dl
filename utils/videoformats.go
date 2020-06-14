package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

//GetVideoFormats gets all available video formats for given video url.
func GetVideoFormats(videoURL string, videoID string, meta map[string]string) (map[string]map[string]string, map[string]string, error) {
	//TODO: show retry info upon debug level

	var videoMetadata = meta
	var playbackURI = videoURL
	var requestHeaders = getRequestHeaders()

	if meta == nil && !strings.Contains(videoURL, "api.hotstar.com") {
		videoURLContent, videoURLDownloadError := getVideoURL(videoURL, requestHeaders)
		if videoURLDownloadError != nil {
			return nil, nil, errors.Wrapf(videoURLDownloadError, "\nGetVideoFormats: Error occurred in retrieving videoURLContent\n")
		}

		var playbackURIError error
		playbackURI, videoMetadata, playbackURIError = getPlayback(videoURLContent, videoURL, videoID, uuid.New().String())
		if playbackURIError != nil {
			return nil, nil, errors.Wrapf(playbackURIError, "\nGetVideoFormats: Error occurred in retrieving playbackURI\n")
		}
	}

	if drmProtected, isDrmKeyAvailable := videoMetadata["drmProtected"]; isDrmKeyAvailable {
		if drmProtected == "true" {
			return nil, nil, errors.New("The content is DRM Protected")
		}
	}

	resultToken, err := getRefreshToken(videoURL)
	if err != nil {
		fmt.Printf("Error in retrieving JWT token : %s", err)
		os.Exit(-1)
	}

	requestHeaders["X-HS-UserToken"] = resultToken

	playbackURIContentBytes, playbackURIContentError := getPlaybackURIContent(playbackURI, requestHeaders)
	if playbackURIContentError != nil {
		return nil, nil, errors.Wrapf(playbackURIContentError, "\nGetVideoFormats: Error occurred in retrieving playbackURIContent\n")
	}

	masterPlaybackURLs, masterPlaybackURLError := getMasterPlaybackURL(playbackURIContentBytes)
	if masterPlaybackURLError != nil {
		return nil, nil, errors.Wrapf(masterPlaybackURLError, "\nGetVideoFormats: Error occurred in retrieving masterPlaybackURLs\n")
	}

	requestHeaders["Referer"] = videoURL
	requestHeaders["Origin"] = "https://www.hotstar.com"
	requestHeaders["Host"] = "hses4.hotstar.com"

	videoFormatsTemp, audioDashFormatsTemp, videoDashFormatsTemp, videoFormatsError := getTempVideoFormats(masterPlaybackURLs, requestHeaders)

	if videoFormatsError != nil {
		return nil, nil, errors.Wrapf(videoFormatsError, "\nGetVideoFormats: Error occurred in retrieving videoFormats\n")
	}

	return getAggregatedFormats(videoFormatsTemp, audioDashFormatsTemp, videoDashFormatsTemp), videoMetadata, nil
}

//ListVideoFormats lists video formats (or) title (or) description of the video for given video url.
func ListVideoFormats(videoURL string, videoID string, metadata map[string]string, titleFlag bool, descriptionFlag bool) {
	videoFormats, videoMetadata, err := GetVideoFormats(videoURL, videoID, metadata)

	if err != nil {
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
		fmt.Printf("Error: %s", err)
		os.Exit(-1)
	}

	if titleFlag || descriptionFlag {
		if titleFlag {
			fmt.Println(videoMetadata["title"])
		}
		if descriptionFlag {
			fmt.Println(videoMetadata["synopsis"])
		}
		if metadata == nil {
			os.Exit(0)
		} else {
			return
		}
	}

	i := 0
	videoFormatsSortedKeys := make([]string, len(videoFormats))
	for formateID := range videoFormats {
		videoFormatsSortedKeys[i] = formateID
		i++
	}

	sort.Strings(videoFormatsSortedKeys)

	//NewWriter(io.Writer, minWidth, tabWidth, padding, padchar, flags)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //tabwriter.Debug
	fmt.Fprintln(tw, "format code\textension\tresolution\tbandwidth\tcodec & frame rate\t")

	for _, formateID := range videoFormatsSortedKeys {

		formatInfo := videoFormats[formateID]

		if mimeType, isMimeTypePresent := formatInfo["MIME-TYPE"]; isMimeTypePresent {
			if mimeType == "video/mp4" {
				fmt.Fprintf(tw, "%s\t%s\tmp4\t%s\t%s\t%s fps\t%s\n", formateID, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"], formatInfo["FRAME-RATE"], formatInfo["STREAM"])
			} else if mimeType == "audio/mp4" {
				fmt.Fprintf(tw, "%s\tm4a\t%s\t%s\t%s\t%s\n", formateID, formatInfo["STREAM"], formatInfo["K-FORM"], formatInfo["CODECS"], formatInfo["SAMPLING-RATE"])
			} else {
				//Handle undefined mime types for dash formats
			}
		} else {
			if frameRate, isFrameRatePresent := formatInfo["FRAME-RATE"]; isFrameRatePresent {
				fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s  %s fps\n", formateID, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"], frameRate)
			} else {
				fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s\n", formateID, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"])
			}
		}
	}
	tw.Flush()

	if metadata == nil {
		os.Exit(0)
	}
}

func isPathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func getFfmpegArgs(videoURL string, videoMetadata map[string]string, streamURL string, dashFiles []string, metadataFlag bool, outputFileName string, isDashFile bool) []string {

	ffmpegArgs := make([]string, 0)
	if !isDashFile {
		ffmpegArgs = append(ffmpegArgs, "-headers")
		ffmpegArgs = append(ffmpegArgs, fmt.Sprintf("Referer: %s", videoURL))
	}
	ffmpegArgs = append(ffmpegArgs, "-i")
	if isDashFile {
		input := "concat:"
		for index, filePath := range dashFiles {
			if index != 0 {
				input = fmt.Sprintf("%s|", input)
			}
			input = fmt.Sprintf("%s%s", input, filePath)
		}
		ffmpegArgs = append(ffmpegArgs, input)
	} else {
		ffmpegArgs = append(ffmpegArgs, streamURL)
	}

	if metadataFlag {
		for metaDataName, metaDataValue := range videoMetadata {
			ffmpegArgs = append(ffmpegArgs, "-metadata")
			metaData := fmt.Sprintf("%s=\"%s\"", metaDataName, metaDataValue)
			ffmpegArgs = append(ffmpegArgs, metaData)
		}
	} else {
		fmt.Println("Skipping adding metadata for video file")
	}

	ffmpegArgs = append(ffmpegArgs, "-c")
	ffmpegArgs = append(ffmpegArgs, "copy")
	ffmpegArgs = append(ffmpegArgs, "-y")
	ffmpegArgs = append(ffmpegArgs, outputFileName)

	return ffmpegArgs
}

func runFfmpegCommand(videoURL string, ffmpegPath string, videoMetadata map[string]string, streamURL string, dashFiles []string, metadataFlag bool, outputFileName string, isDashFile bool) {

	var stdoutBuf, stderrBuf bytes.Buffer

	ffmpegArgs := getFfmpegArgs(videoURL, videoMetadata, streamURL, dashFiles, metadataFlag, outputFileName, isDashFile)

	ffmpegCmd := exec.Command(ffmpegPath, ffmpegArgs...)

	if isDashFile {
		fmt.Println("\nStarting ffmpeg to merge downloaded DASH audio/video...")
	} else {
		fmt.Println("Starting ffmpeg to download video...")
	}

	stdoutIn, _ := ffmpegCmd.StdoutPipe()
	stderrIn, _ := ffmpegCmd.StderrPipe()

	var errStdout, errStderr error

	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := ffmpegCmd.Start()

	if err != nil {
		log.Fatalf("ffmpegCmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = ffmpegCmd.Wait()
	if err != nil {
		log.Fatalf("ffmpegCmd.Run() failed with %s\n", err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

}

func getBestOrLeastResolutionFormat(videoFormats map[string]map[string]string, bestOrLeast string) string {

	for formatCode, formatInfo := range videoFormats {
		if strings.HasPrefix(formatCode, "hls-") {
			//video format
			if isBestOrLeast, bestOrLeastAvailable := formatInfo[bestOrLeast]; bestOrLeastAvailable {
				if isBestOrLeast == "true" {
					return formatCode
				}
			}
		}
	}

	return ""
}

func downloadDashAudioOrVideo(videoURL string, videoFormats map[string]map[string]string, vFormat string, outputFileName string, videoID string, videoMetadata map[string]string, currentDirectoryPath string, ffmpegPath string, metadataFlag bool) {
	format := videoFormats[vFormat]
	if outputFileName == "" {
		outputFileName = fmt.Sprintf("%s_%s__DASH_AV.mp4", videoID, strings.Replace(videoMetadata["title"], " ", "_", -1))
	}
	outputFilePath := filepath.Join(currentDirectoryPath, outputFileName)

	if isPathExists(outputFilePath) {
		fmt.Printf("File %s already present in %s", outputFileName, currentDirectoryPath)
		os.Exit(0)
	}

	requestHeaders := getRequestHeaders()

	requestHeaders["Referer"] = videoURL
	requestHeaders["Origin"] = "https://www.hotstar.com"
	requestHeaders["Host"] = "hses4.hotstar.com"

	dashFiles, tempDashFileDir := DownloadDashFilesBatch(currentDirectoryPath, videoID, vFormat, format, requestHeaders)
	runFfmpegCommand(videoURL, ffmpegPath, videoMetadata, "", dashFiles, metadataFlag, outputFileName, true)
	removeErr := os.RemoveAll(tempDashFileDir)
	if removeErr != nil {
		fmt.Println("Error in removing temp directory")
		os.Exit(-1)
	} else {
		fmt.Printf("\nTemp directory %s removed\n", tempDashFileDir)
		os.Exit(0)
	}
}

func downloadVideo(videoURL string, vFormat string, videoFormats map[string]map[string]string, outputFileName string, videoID string, videoMetadata map[string]string, currentDirectoryPath string, ffmpegPath string, metadataFlag bool) {
	//Check if vFormat fallback is enabled by empty value passed
	if len(strings.TrimSpace(vFormat)) == 0 {
		fmt.Println("Missing format flag falling back to best formats for video")
		vFormat = getBestOrLeastResolutionFormat(videoFormats, "BEST_RESOLUTION")
		if len(strings.TrimSpace(vFormat)) != 0 {
			fmt.Println("Best format for video is, ", videoFormats[vFormat]["RESOLUTION"])
		} else {
			vFormat = getBestOrLeastResolutionFormat(videoFormats, "LEAST_RESOLUTION")
			fmt.Println("Best formats for the video isn't available falling back to least resolution, ", videoFormats[vFormat]["RESOLUTION"])

		}
	}

	if videoFormat, isValidFormat := videoFormats[vFormat]; isValidFormat {

		if streamURL, isStreamURLAvailable := videoFormat["STREAM-URL"]; isStreamURLAvailable {

			if outputFileName == "" {
				outputFileName = fmt.Sprintf("%s-%s.mp4", videoID, strings.Replace(videoMetadata["title"], " ", "_", -1))
			}

			outputFilePath := filepath.Join(currentDirectoryPath, outputFileName)

			if isPathExists(outputFilePath) {
				fmt.Printf("File %s already present in %s", outputFileName, currentDirectoryPath)
				os.Exit(0)
			}

			runFfmpegCommand(videoURL, ffmpegPath, videoMetadata, streamURL, nil, metadataFlag, outputFileName, false)
			os.Exit(0)
		} else {
			fmt.Println("The STREAM-URL is not available. Please try again")
			os.Exit(-3)
		}

	} else {
		fmt.Printf("The specified video format %s is not available. Specify existing format from the list", vFormat)
		os.Exit(-4)
	}
}

//DownloadAudioOrVideo downloads the video for given video format and video url. It also adds metadata to it if needed. FFMPEG path and Output video file name can be customized.
func DownloadAudioOrVideo(videoURL string, videoID string, vFormat string, userFfmpegPath string, outputFileName string, metadataFlag bool, isDashAV bool) {

	var ffmpegPath string

	if len(strings.TrimSpace(userFfmpegPath)) != 0 {
		ffmpegPath = userFfmpegPath
	} else {
		path, err := exec.LookPath("ffmpeg")
		if err != nil {
			log.Fatal("Error in finding command ffmpeg. Please install one and try again. ", err)
		}
		ffmpegPath = path
	}

	currentDirectoryPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	videoFormats, videoMetadata, err := GetVideoFormats(videoURL, videoID, nil)

	if err != nil {
		log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	if drmProtected, isDrmKeyAvailable := videoMetadata["drmProtected"]; isDrmKeyAvailable {
		if drmProtected == "true" {
			fmt.Println("Error: The video is DRM Protected")
			os.Exit(-1)
		}
	}

	if err := os.Chmod(ffmpegPath, 0555); err != nil {
		log.Fatal(err)
	}

	if isDashAV {
		downloadDashAudioOrVideo(videoURL, videoFormats, vFormat, outputFileName, videoID, videoMetadata, currentDirectoryPath, ffmpegPath, metadataFlag)
		fmt.Println("Downloaded DASH audio/video successfully...")
	} else {
		downloadVideo(videoURL, vFormat, videoFormats, outputFileName, videoID, videoMetadata, currentDirectoryPath, ffmpegPath, metadataFlag)
		fmt.Println("Downloaded video successfully...")
	}

}

//ListOrDownloadPlaylistVideoFormats lists video formats (or) title (or) description of each video url in the list.
func ListOrDownloadPlaylistVideoFormats(playlistID string, titleFlag bool, descriptionFlag bool, playlistStartRange string, playlistEndRange string, isDownloadSwitch bool, vFormat string, userFfmpegPath string, outputFileName string, metadataFlag bool, isDashAV bool) {
	var result map[string]interface{}
	playlistURI := fmt.Sprintf("https://api.hotstar.com/o/v1/tray/find?uqId=%s&tas=10000", playlistID)

	playlistURIContentBytes, err := MakeGetRequest(playlistURI, getRequestHeaders())

	if err != nil {
		panic(err)
	}

	json.Unmarshal(playlistURIContentBytes, &result)

	statusCode, isStatusCastOk := result["statusCodeValue"].(float64)

	if isStatusCastOk {
		if statusCode != 200 {
			panic(errors.New("Invalid status code returned"))
		} else {
			body, isBodyCastOk := result["body"].(map[string]interface{})

			if isBodyCastOk {
				results := body["results"].(map[string]interface{})
				assets := results["assets"].(map[string]interface{})
				items := assets["items"].([]interface{})
				playlistItemCount := len(items)

				if strings.Compare(playlistStartRange, "") == 0 {
					fmt.Println("Start range not specified falling back to upper bound, 1")
					playlistStartRange = "1"
				}

				if strings.Compare(playlistEndRange, "") == 0 {
					fmt.Println("End range not specified falling back to lower bound,", playlistItemCount)
					playlistEndRange = fmt.Sprintf("%d", playlistItemCount)
				}

				fmt.Printf("\nCollected %d video id(s) from playlist\n", playlistItemCount)

				startRange, endRange, boundValidationErrors, isValidBounds := isValidPlaylistBounds(playlistItemCount, playlistStartRange, playlistEndRange)

				if !isValidBounds {
					fmt.Println(boundValidationErrors)
					os.Exit(-1)
				}

				for itemIndex := (playlistItemCount - 1) - (startRange - 1); itemIndex >= (playlistItemCount-1)-(endRange-1); itemIndex-- {
					itemValue := items[itemIndex]
					itemInfo := itemValue.(map[string]interface{})

					var metadata = make(map[string]interface{})
					for contentKey, contentValue := range itemInfo {
						metadata[contentKey] = contentValue
					}

					metaDataMap := make(map[string]string)
					PopulateMetaDataMapWithMetadata(metaDataMap, metadata)

					fmt.Printf("\nFor video id, %s\n", metaDataMap["id"])

					if !isDownloadSwitch {
						ListVideoFormats(GetPlaybackURI2(metaDataMap["id"], uuid.New().String()), metaDataMap["id"], metaDataMap, titleFlag, descriptionFlag)
					} else {
						DownloadAudioOrVideo(GetPlaybackURI2(metaDataMap["id"], uuid.New().String()), metaDataMap["id"], vFormat, userFfmpegPath, outputFileName, metadataFlag, isDashAV)
					}
				}
			} else {
				fmt.Println("result body to map[string]interface{} cast unsuccessful")
				os.Exit(-1)
			}
		}
	} else {
		fmt.Println("result statusCodeValue to map[string]interface{} cast unsuccessful")
		os.Exit(-1)
	}

}
