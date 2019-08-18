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

func DownloadDashFile(filepath string, url string) error {

	//Create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	//Get data
	resp, err := http.Get(url)
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

func getSegmentUrl(playbackUrl, streamId string) string {
	return strings.Replace(playbackUrl, "master.mpd", streamId, -1)
}

func raiseFileDownloadError(err error) {
	fmt.Println("Error in downloading file. Error:", err)
	os.Exit(-1)
}

func DownloadDashFilesBatch(currentDirectoryPath, videoId string, vFormatCode string, format map[string]string) ([]string, string) {
	var dashFiles []string
	tempFolder := fmt.Sprintf("temp_%s_%s", videoId, vFormatCode)
	tempDir := filepath.Join(currentDirectoryPath, tempFolder)

	if _, dirExistenceErr := os.Stat(tempDir); dirExistenceErr == nil {
		fmt.Println("Temp directory exists from previous run.")
		removeErr := os.RemoveAll(tempDir)
		if removeErr != nil {
			fmt.Println("Error in removing temp directory")
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

	initSegmentUrl := getSegmentUrl(format["PLAYBACK-URL"], format["INIT-URL"])
	initSegmentUrlValues := strings.Split(format["INIT-URL"], "/")
	initFilePath := filepath.Join(tempFolder, initSegmentUrlValues[len(initSegmentUrlValues)-1])
	dashFiles = append(dashFiles, initFilePath)
	initFileErr := DownloadDashFile(initFilePath, initSegmentUrl)
	if initFileErr != nil {
		raiseFileDownloadError(initFileErr)
	}
	for _, segmentNum := range makeRange(1, totalSegments) {
		streamUrl := strings.Replace(format["STREAM-URL"], "$Number$", fmt.Sprintf("%d", segmentNum), -1)
		streamUrlValues := strings.Split(streamUrl, "/")
		segmentUrl := getSegmentUrl(format["PLAYBACK-URL"], streamUrl)
		segmentFilePath := filepath.Join(tempFolder, streamUrlValues[len(streamUrlValues)-1])
		dashFiles = append(dashFiles, segmentFilePath)
		segmentFileErr := DownloadDashFile(segmentFilePath, segmentUrl)
		if segmentFileErr != nil {
			raiseFileDownloadError(initFileErr)
		}
		bar.Increment()
	}

	bar.Finish()

	return dashFiles, tempDir
}
