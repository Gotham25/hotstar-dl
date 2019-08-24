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

func downloadDashFile(filepath string, url string) error {

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

func getSegmentURL(playbackURL, streamID string) string {
	return strings.Replace(playbackURL, "master.mpd", streamID, -1)
}

func raiseFileDownloadError(err error) {
	fmt.Println("Error in downloading file. Error:", err)
	os.Exit(-1)
}

//DownloadDashFilesBatch downloads the dash chunks for the given video format
func DownloadDashFilesBatch(currentDirectoryPath, videoID string, vFormatCode string, format map[string]string) ([]string, string) {
	var dashFiles []string
	tempFolder := fmt.Sprintf("temp_%s_%s", videoID, vFormatCode)
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

	initSegmentURL := getSegmentURL(format["PLAYBACK-URL"], format["INIT-URL"])
	initSegmentURLValues := strings.Split(format["INIT-URL"], "/")
	initFilePath := filepath.Join(tempFolder, initSegmentURLValues[len(initSegmentURLValues)-1])
	dashFiles = append(dashFiles, initFilePath)
	initFileErr := downloadDashFile(initFilePath, initSegmentURL)
	if initFileErr != nil {
		raiseFileDownloadError(initFileErr)
	}
	for _, segmentNum := range MakeRange(1, totalSegments) {
		streamURL := strings.Replace(format["STREAM-URL"], "$Number$", fmt.Sprintf("%d", segmentNum), -1)
		streamURLValues := strings.Split(streamURL, "/")
		segmentURL := getSegmentURL(format["PLAYBACK-URL"], streamURL)
		segmentFilePath := filepath.Join(tempFolder, streamURLValues[len(streamURLValues)-1])
		dashFiles = append(dashFiles, segmentFilePath)
		segmentFileErr := downloadDashFile(segmentFilePath, segmentURL)
		if segmentFileErr != nil {
			raiseFileDownloadError(initFileErr)
		}
		bar.Increment()
	}

	bar.Finish()

	return dashFiles, tempDir
}
