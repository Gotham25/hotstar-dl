package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/google/uuid"
)

var videoFormatsRetryCount = 0

//GetVideoFormats gets all available video formats for given video url.
func GetVideoFormats(videoUrl string, videoId string) (map[string]map[string]string, map[string]string, error) {
	//TODO: show retry info upon debug level

	var requestHeaders = map[string]string{
		"Hotstarauth":     GenerateHotstarAuth(),
		"X-Country-Code":  "IN",
		"X-Platform-Code": "JIO",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.122 Safari/537.36",
	}

	videoUrlContentBytes, err := MakeGetRequest(videoUrl, requestHeaders)
	totalFormats := make(map[string]map[string]string)
	videoFormatsTemp := make(map[string][]map[string]string)
	videoDashFormatsTemp := make(map[string][]map[string]string)
	audioDashFormatsTemp := make(map[string][]map[string]string)

	if err != nil {
		if videoFormatsRetryCount+1 < 10 {
			//retry again for fetching formats
			videoFormatsRetryCount++
			//fmt.Printf("GetVideoFormats: GET request to videoUrl failed... Retrying count : #%d\n", videoFormatsRetryCount)
			return GetVideoFormats(videoUrl, videoId)
		}
		return nil, nil, err
	}

	videoUrlContent := fmt.Sprintf("%s", videoUrlContentBytes)

	playbackUri, videoMetadata, err := GetPlaybackUri(videoUrlContent, videoUrl, videoId, uuid.New().String())

	if err != nil {
		if videoFormatsRetryCount+1 < 10 {
			//retry again for fetching formats
			videoFormatsRetryCount++
			//fmt.Printf("GetVideoFormats: Invalid APP_STATE json... retrying count : #%d\n", videoFormatsRetryCount)
			return GetVideoFormats(videoUrl, videoId)
		}
		return nil, nil, err
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	if drmProtected, isDrmKeyAvailable := videoMetadata["drmProtected"]; isDrmKeyAvailable {
		if drmProtected == "true" {
			return nil, nil, fmt.Errorf("the content is DRM Protected")
		}
	}

	playbackUriContentBytes, err := MakeGetRequest(playbackUri, requestHeaders)

	if err != nil {
		if videoFormatsRetryCount+1 < 10 {
			//retry again for fetching formats
			videoFormatsRetryCount++
			//fmt.Printf("GetVideoFormats: GET request to playbackUri failed... Retrying count : #%d\n", videoFormatsRetryCount)
			return GetVideoFormats(videoUrl, videoId)
		}
		log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	masterPlaybackUrls, err := GetMasterPlaybackUrls(playbackUriContentBytes)

	if err != nil {
		if videoFormatsRetryCount+1 < 10 {
			//retry again for fetching formats
			videoFormatsRetryCount++
			//fmt.Printf("GetVideoFormats: Retriving masterPlaybackUrl failed... Retrying count : #%d\n", videoFormatsRetryCount)
			return GetVideoFormats(videoUrl, videoId)
		}
		return nil, nil, err
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	for _, masterPlaybackUrl := range masterPlaybackUrls {

		if masterPlaybackUrl != "" {

			var queryParams string
			masterPlaybackUrlQueryParam := strings.Split(masterPlaybackUrl, "?")

			if len(masterPlaybackUrlQueryParam) > 1 {
				queryParams = masterPlaybackUrlQueryParam[1]
			}

			if strings.Contains(masterPlaybackUrl, "m3u8") {

				masterPlaybackPageContentsM3u8Bytes, err := MakeGetRequest(masterPlaybackUrl, requestHeaders)

				if err != nil {

					if videoFormatsRetryCount+1 < 10 {
						//retry again for fetching formats
						videoFormatsRetryCount++
						//fmt.Printf("GetVideoFormats: GET request to masterPlaybackUrl failed... Retrying count : #%d\n", videoFormatsRetryCount)
						return GetVideoFormats(videoUrl, videoId)
					}

					return nil, nil, err
				}

				for fid, formatsList := range ParseM3u8Content(fmt.Sprintf("%s", masterPlaybackPageContentsM3u8Bytes), masterPlaybackUrl, queryParams) {
					videoFormatsTemp[fid] = append(videoFormatsTemp[fid], formatsList)
				}
			} else {

				masterPlaybackPageContentsMpdBytes, err := MakeGetRequest(masterPlaybackUrl, requestHeaders)

				if err != nil {

					if videoFormatsRetryCount+1 < 10 {
						//retry again for fetching formats
						videoFormatsRetryCount++
						//fmt.Printf("GetVideoFormats: GET request to masterPlaybackUrl failed... Retrying count : #%d\n", videoFormatsRetryCount)
						return GetVideoFormats(videoUrl, videoId)
					}

					return nil, nil, err
				}

				dFormats := GetDashFormats(masterPlaybackPageContentsMpdBytes, masterPlaybackUrl)

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

	return totalFormats, videoMetadata, nil
}

//ListVideoFormats lists video formats (or) title (or) description of the video for given video url.
func ListVideoFormats(videoUrl string, videoId string, titleFlag bool, descriptionFlag bool) {
	videoFormats, videoMetadata, err := GetVideoFormats(videoUrl, videoId)

	if err != nil {
		log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	if titleFlag || descriptionFlag {
		if titleFlag {
			fmt.Println(videoMetadata["title"])
		}
		if descriptionFlag {
			fmt.Println(videoMetadata["synopsis"])
		}
		os.Exit(0)
	}

	i := 0
	videoFormatsSortedKeys := make([]string, len(videoFormats))
	for formateId := range videoFormats {
		videoFormatsSortedKeys[i] = formateId
		i++
	}

	sort.Strings(videoFormatsSortedKeys)

	//NewWriter(io.Writer, minWidth, tabWidth, padding, padchar, flags)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //tabwriter.Debug
	fmt.Fprintln(tw, "format code\textension\tresolution\tbandwidth\tcodec & frame rate\t")

	for _, formateId := range videoFormatsSortedKeys {

		formatInfo := videoFormats[formateId]

		if mimeType, isMimeTypePresent := formatInfo["MIME-TYPE"]; isMimeTypePresent {
			if mimeType == "video/mp4" {
				fmt.Fprintf(tw, "%s\t%s\tmp4\t%s\t%s\t%s fps\t%s\n", formateId, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"], formatInfo["FRAME-RATE"], formatInfo["STREAM"])
			} else if mimeType == "audio/mp4" {
				fmt.Fprintf(tw, "%s\tm4a\t%s\t%s\t%s\t%s\n", formateId, formatInfo["STREAM"], formatInfo["K-FORM"], formatInfo["CODECS"], formatInfo["SAMPLING-RATE"])
			} else {
				//Handle undefined mime types for dash formats
			}
		} else {
			if frameRate, isFrameRatePresent := formatInfo["FRAME-RATE"]; isFrameRatePresent {
				fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s  %s fps\n", formateId, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"], frameRate)
			} else {
				fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s\n", formateId, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"])
			}
		}
	}
	tw.Flush()
	os.Exit(0)
}

func isPathExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func getFfmpegArgs(videoMetadata map[string]string, streamUrl string, dashFiles []string, metadataFlag bool, outputFileName string, isDashFile bool) []string {
	ffmpegArgs := make([]string, 0)
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
		ffmpegArgs = append(ffmpegArgs, streamUrl)
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

func runFfmpegCommand(ffmpegPath string, videoMetadata map[string]string, streamUrl string, dashFiles []string, metadataFlag bool, outputFileName string, isDashFile bool) {

	var stdoutBuf, stderrBuf bytes.Buffer

	ffmpegArgs := getFfmpegArgs(videoMetadata, streamUrl, dashFiles, metadataFlag, outputFileName, isDashFile)

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

//DownloadAudioOrVideo downloads the video for given video format and video url. It also adds metadata to it if needed. FFMPEG path and Output video file name can be customized.
func DownloadAudioOrVideo(videoUrl string, videoId string, vFormat string, userFfmpegPath string, outputFileName string, metadataFlag bool, isDashAV bool) {

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

	videoFormats, videoMetadata, err := GetVideoFormats(videoUrl, videoId)

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
		format := videoFormats[vFormat]
		if outputFileName == "" {
			outputFileName = fmt.Sprintf("%s_%s__DASH_AV.mp4", strings.Replace(videoMetadata["title"], " ", "_", -1), videoId)
		}
		outputFilePath := filepath.Join(currentDirectoryPath, outputFileName)

		if isPathExists(outputFilePath) {
			fmt.Printf("File %s already present in %s", outputFileName, currentDirectoryPath)
			os.Exit(0)
		}
		dashFiles, tempDashFileDir := DownloadDashFilesBatch(currentDirectoryPath, videoId, vFormat, format)
		runFfmpegCommand(ffmpegPath, videoMetadata, "", dashFiles, metadataFlag, outputFileName, true)
		removeErr := os.RemoveAll(tempDashFileDir)
		if removeErr != nil {
			fmt.Println("Error in removing temp directory")
			os.Exit(-1)
		} else {
			fmt.Printf("\nTemp directory %s removed\n", tempDashFileDir)
			os.Exit(0)
		}
	} else {
		if videoFormat, isValidFormat := videoFormats[vFormat]; isValidFormat {

			if streamUrl, isStreamUrlAvailable := videoFormat["STREAM-URL"]; isStreamUrlAvailable {

				if outputFileName == "" {
					outputFileName = fmt.Sprintf("%s-%s.mp4", strings.Replace(videoMetadata["title"], " ", "_", -1), videoId)
				}

				outputFilePath := filepath.Join(currentDirectoryPath, outputFileName)

				if isPathExists(outputFilePath) {
					fmt.Printf("File %s already present in %s", outputFileName, currentDirectoryPath)
					os.Exit(0)
				}

				runFfmpegCommand(ffmpegPath, videoMetadata, streamUrl, nil, metadataFlag, outputFileName, false)
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

}
