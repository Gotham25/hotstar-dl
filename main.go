package main

import (
	"flag"
	"fmt"
	"github.com/Gotham25/hotstar-dl/utils"
	"os"
	"strings"
)

//Build version info vars injected by goreleaser
var version string
var commit string
var date string

//flag descriptions
var helpFlagDesc = "Prints this help and exit"
var listFormatsFlagDesc = "List available video formats for given url"
var formatFlagDesc = "Video format to download video in specified resolution"
var ffmpegPathFlagDesc = "Location of the ffmpeg binary(absolute path)"
var metadataFlagDesc = "Add metadata to the video file"
var outputFileNameFlagDesc = "Output file name"
var titleFlagDesc = "Prints video title and exit"
var descriptionFlagDesc = "Prints video description and exit"
var versionFlagDesc = "Prints version info and exits"

//flag declarations
var helpFlag = flag.Bool("help", false, helpFlagDesc)
var listFormatsFlag = flag.Bool("list", false, listFormatsFlagDesc)
var formatFlag = flag.String("format", "", formatFlagDesc)
var ffmpegPathFlag = flag.String("ffmpeg-location", "", ffmpegPathFlagDesc)
var metadataFlag = flag.Bool("add-metadata", false, metadataFlagDesc)
var outputFileNameFlag = flag.String("output", "", outputFileNameFlagDesc)
var titleFlag = flag.Bool("get-title", false, titleFlagDesc)
var descriptionFlag = flag.Bool("get-description", false, descriptionFlagDesc)
var versionFlag = flag.Bool("version", false, versionFlagDesc)

func init() {
	//shorthand notations
	flag.BoolVar(helpFlag, "h", false, helpFlagDesc)
	flag.BoolVar(listFormatsFlag, "l", false, listFormatsFlagDesc)
	flag.StringVar(formatFlag, "f", "", formatFlagDesc)
	flag.BoolVar(metadataFlag, "m", false, metadataFlagDesc)
	flag.StringVar(outputFileNameFlag, "o", "", outputFileNameFlagDesc)
	flag.BoolVar(titleFlag, "t", false, titleFlagDesc)
	flag.BoolVar(descriptionFlag, "i", false, descriptionFlagDesc)
	flag.BoolVar(versionFlag, "v", false, versionFlagDesc)

	//custom flag usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Fprintf(os.Stdout, "-h, --help\t\t%s\n", helpFlagDesc)
		fmt.Fprintf(os.Stdout, "-l, --list\t\t%s\n", listFormatsFlagDesc)
		fmt.Fprintf(os.Stdout, "-f, --format\t\t%s\n", formatFlagDesc)
		fmt.Fprintf(os.Stdout, "--ffmpeg-location\t%s\n", ffmpegPathFlagDesc)
		fmt.Fprintf(os.Stdout, "-m, --add-metadata\t%s\n", metadataFlagDesc)
		fmt.Fprintf(os.Stdout, "-t, --get-title\t\t%s\n", titleFlagDesc)
		fmt.Fprintf(os.Stdout, "-i, --get-description\t%s\n", descriptionFlagDesc)
		fmt.Fprintf(os.Stdout, "-o, --output\t\t%s\n", outputFileNameFlagDesc)
		fmt.Fprintf(os.Stdout, "-v, --version\t\t%s\n", versionFlagDesc)
		os.Exit(0)
		//flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()
	flagCount := len(flag.Args())
	if *helpFlag {
		flag.Usage()
	} else if *versionFlag {
		fmt.Printf("Version : %s\ngit commit SHA : %s \nBuilt on : %s\n", version, commit, date)
	} else if flagCount == 0 {
		fmt.Println("Must provide atleast one url at end")
		flag.Usage()
		os.Exit(-1)
	} else if flagCount > 1 {
		fmt.Println("Url must be provided at end before options")
		flag.Usage()
		os.Exit(-1)
	} else if videoUrl := flag.Args()[0]; videoUrl != "" {

		videoUrl = utils.GetParsedVideoUrl(videoUrl)

		isValidUrl, videoId := utils.IsValidHotstarUrl(videoUrl)
		if isValidUrl {
			if *listFormatsFlag || *titleFlag || *descriptionFlag {
				//list video formats

				utils.ListVideoFormats(videoUrl, videoId, *titleFlag, *descriptionFlag)

			} else if *formatFlag != "" {
				if !strings.HasPrefix(*formatFlag, "hls-") {
					fmt.Println("Invalid format specified")
					os.Exit(-1)
				} else {

					utils.DownloadVideo(videoUrl, videoId, *formatFlag, *ffmpegPathFlag, *outputFileNameFlag, *metadataFlag)

				}
			} else {
				//TODO: Check for other flags if associated with url if any
			}
		} else {
			fmt.Println("Invalid hotstar url. Please enter a valid one")
			os.Exit(-1)
		}

	} else {
		fmt.Println("Invalid args specified")
		flag.Usage()
	}

}
