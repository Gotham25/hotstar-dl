package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

var videoFormatsRetryCount = 0

//GetVideoFormats gets all available video formats for given video url.
func GetVideoFormats(videoUrl string, videoId string) (map[string]map[string]string, map[string]string, error) {
	//TODO: show retry info upon debug level

	var requestHeaders = map[string]string{
		"Hotstarauth":     GenerateHotstarAuth(),
		"X-Country-Code":  "IN",
		"X-Platform-Code": "JIO",
	}

	videoUrlContentBytes, err := MakeGetRequest(videoUrl, requestHeaders)

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

	playbackUri, videoMetadata, err := GetPlaybackUri(videoUrlContent, videoUrl, videoId)

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
			return nil, nil, fmt.Errorf("The content is DRM Protected.")
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

	masterPlaybackUrl, err := GetMasterPlaybackUrl(playbackUriContentBytes)

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

	var queryParams string
	masterPlaybackUrlQueryParam := strings.Split(masterPlaybackUrl, "?")

	if len(masterPlaybackUrlQueryParam) > 1 {
		queryParams = masterPlaybackUrlQueryParam[1]
	}

	masterPlaybackPageContentsBytes, err := MakeGetRequest(masterPlaybackUrl, requestHeaders)

	if err != nil {
		if videoFormatsRetryCount+1 < 10 {
			//retry again for fetching formats
			videoFormatsRetryCount++
			//fmt.Printf("GetVideoFormats: GET request to masterPlaybackUrl failed... Retrying count : #%d\n", videoFormatsRetryCount)
			return GetVideoFormats(videoUrl, videoId)
		}
		return nil, nil, err
		//log.Fatal(fmt.Errorf("Error occurred : %s", err))
	}

	//fmt.Printf("\nmasterPlaybackPageContentsBytes : \n%s\n", masterPlaybackPageContentsBytes)

	//return fmt.Sprintf("%s", masterPlaybackPageContentsBytes)

	return ParseM3u8Content(fmt.Sprintf("%s", masterPlaybackPageContentsBytes), masterPlaybackUrl, queryParams), videoMetadata, nil

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

	//NewWriter(io.Writer, minWidth, tabWidth, padding, padchar, flags)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0) //tabwriter.Debug
	fmt.Fprintln(tw, "format code\textension\tresolution\tbandwidth\tcodec & frame rate\t")

	for formateId, formatInfo := range videoFormats {
		if frameRate, isFrameRatePresent := formatInfo["FRAME-RATE"]; isFrameRatePresent {
			fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s  %s fps\n", formateId, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"], frameRate)
		} else {
			fmt.Fprintf(tw, "%s\tmp4\t%s\t%s\t%s\n", formateId, formatInfo["RESOLUTION"], formatInfo["K-FORM"], formatInfo["CODECS"])
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

func getFfmpegArgs(videoMetadata map[string]string, streamUrl string, metadataFlag bool, outputFileName string) []string {
	ffmpegArgs := make([]string, 0)
	ffmpegArgs = append(ffmpegArgs, "-i")
	ffmpegArgs = append(ffmpegArgs, streamUrl)

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

func runFfmpegCommand(ffmpegPath string, videoMetadata map[string]string, streamUrl string, metadataFlag bool, outputFileName string) {

	var stdoutBuf, stderrBuf bytes.Buffer

	ffmpegArgs := getFfmpegArgs(videoMetadata, streamUrl, metadataFlag, outputFileName)

	ffmpegCmd := exec.Command(ffmpegPath, ffmpegArgs...)

	fmt.Println("Starting ffmpeg to download video...")

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

	os.Exit(0)

}

//DownloadVideo downloads the video for given video format and video url. It also adds metadata to it if needed. FFMPEG path and Output video file name can be customized.
func DownloadVideo(videoUrl string, videoId string, vFormat string, userFfmpegPath string, outputFileName string, metadataFlag bool) {

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

	if videoFormat, isValidFormat := videoFormats[vFormat]; isValidFormat {

		if streamUrl, isStreamUrlAvailable := videoFormat["STREAM-URL"]; isStreamUrlAvailable {

			if outputFileName == "" {
				outputFileName = fmt.Sprintf("%s.mp4", strings.Replace(videoMetadata["title"], " ", "_", -1))
			}

			outputFilePath := filepath.Join(currentDirectoryPath, outputFileName)

			if isPathExists(outputFilePath) {
				fmt.Printf("File %s already present in %s", outputFileName, currentDirectoryPath)
				os.Exit(0)
			}

			if err := os.Chmod(ffmpegPath, 0555); err != nil {
				log.Fatal(err)
			}

			runFfmpegCommand(ffmpegPath, videoMetadata, streamUrl, metadataFlag, outputFileName)

		} else {
			fmt.Println("The STREAM-URL is not available. Please try again")
			os.Exit(-3)
		}

	} else {
		fmt.Printf("The specified video format %s is not available. Specify existing format from the list", vFormat)
		os.Exit(-4)
	}
}
