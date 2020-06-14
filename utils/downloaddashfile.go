package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func downloadDashFile(filepath string, url string, requestHeaders map[string]string) error {

	//Create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	//Creating a custom request
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)

	//Adding custom headers
	for requestHeaderKey, requestHeaderValue := range requestHeaders {
		request.Header.Add(requestHeaderKey, requestHeaderValue)
	}

	//Get data
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getSegmentURL(playbackURL, streamID string) string {
	return strings.Replace(playbackURL, "master.mpd", streamID, -1)
}

func raiseFileDownloadError(err error) {
	fmt.Println("Error in downloading file. Error:", err)
	os.Exit(-1)
}

//DownloadDashFilesBatch downloads the dash chunks for the given video format
func DownloadDashFilesBatch(currentDirectoryPath, videoID string, vFormatCode string, format map[string]string, requestHeaders map[string]string) ([]string, string) {
	var dashFiles []string
	tempFolder := fmt.Sprintf("temp_%s_%s", videoID, vFormatCode)
	tempDir := filepath.Join(currentDirectoryPath, tempFolder)

	if _, dirExistenceErr := os.Stat(tempDir); dirExistenceErr == nil {
		fmt.Printf("\nTemp %s directory exists from previous run.\n", tempFolder)
		removeErr := os.RemoveAll(tempDir)
		if removeErr != nil {
			fmt.Printf("\nError in removing temp directory %s\n", tempFolder)
			os.Exit(-1)
		} else {
			fmt.Printf("\nTemp directory %s removed\n", tempFolder)
		}
	}

	dirCreationErr := os.Mkdir(tempDir, os.ModeDir)

	if dirCreationErr != nil || os.IsNotExist(dirCreationErr) {
		fmt.Print("\nError in creating temp directory\n")
		os.Exit(-1)
	}

	fmt.Printf("\nTemp directory %s created\n", tempFolder)

	totalSegments, _ := strconv.Atoi(format["TOTAL-SEGMENTS"])

	bar := pb.StartNew(totalSegments)

	fmt.Printf("\nDownloading DASH chunks to above directory\n")

	dashFiles = make([]string, 0)

	initSegmentURL := getSegmentURL(format["PLAYBACK-URL"], format["INIT-URL"])
	initSegmentURLValues := strings.Split(format["INIT-URL"], "/")
	initFilePath := filepath.Join(tempFolder, initSegmentURLValues[len(initSegmentURLValues)-1])
	dashFiles = append(dashFiles, initFilePath)
	initFileErr := downloadDashFile(initFilePath, initSegmentURL, requestHeaders)
	if initFileErr != nil {
		raiseFileDownloadError(initFileErr)
	}
	for _, segmentNum := range MakeRange(1, totalSegments) {
		streamURL := strings.Replace(format["STREAM-URL"], "$Number$", fmt.Sprintf("%d", segmentNum), -1)
		streamURLValues := strings.Split(streamURL, "/")
		segmentURL := getSegmentURL(format["PLAYBACK-URL"], streamURL)
		segmentFilePath := filepath.Join(tempFolder, streamURLValues[len(streamURLValues)-1])
		dashFiles = append(dashFiles, segmentFilePath)
		segmentFileErr := downloadDashFile(segmentFilePath, segmentURL, requestHeaders)
		if segmentFileErr != nil {
			raiseFileDownloadError(initFileErr)
		}
		bar.Increment()
	}

	bar.Finish()

	return dashFiles, tempDir
}
