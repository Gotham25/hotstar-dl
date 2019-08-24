package utils

import (
	"encoding/xml"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//SegmentTemplate struct contains info about segments
type SegmentTemplate struct {
	Duration       string `xml:"duration,attr"`
	Initialization string `xml:"initialization,attr"`
	Media          string `xml:"media,attr"`
	StartNumber    string `xml:"startNumber,attr"`
	Timescale      string `xml:"timescale,attr"`
}

//AudioChannelConfiguration struct contains config of channel
type AudioChannelConfiguration struct {
	SchemeIDURI string `xml:"schemeIdUri,attr"`
	Value       string `xml:"value,attr"`
}

//Representation struct contains
type Representation struct {
	Bandwidth                 string `xml:"bandwidth,attr"`
	Codecs                    string `xml:"codecs,attr"`
	FrameRate                 string `xml:"frameRate,attr"`
	Height                    string `xml:"height,attr"`
	ID                        string `xml:"id,attr"`
	ScanType                  string `xml:"scanType,attr"`
	Width                     string `xml:"width,attr"`
	AudioSamplingRate         string `xml:"audioSamplingRate,attr"`
	AudioChannelConfiguration string `xml:"AudioChannelConfiguration"`
}

//AdaptationSet struct contains
type AdaptationSet struct {
	MaxHeight        string           `xml:"maxHeight,attr"`
	MaxWidth         string           `xml:"maxWidth,attr"`
	MimeType         string           `xml:"mimeType,attr"`
	SegmentAlignment string           `xml:"segmentAlignment,attr"`
	StartWithSAP     string           `xml:"startWithSAP,attr"`
	SegTemplate      SegmentTemplate  `xml:"SegmentTemplate"`
	Representations  []Representation `xml:"Representation"`
}

//MPD struct contains
type MPD struct {
	MediaPresentationDuration string          `xml:"mediaPresentationDuration,attr"`
	MinBufferTime             string          `xml:"minBufferTime,attr"`
	Profiles                  string          `xml:"profiles,attr"`
	Xmlns                     string          `xml:"xmlns,attr"`
	Period                    []AdaptationSet `xml:"Period>AdaptationSet"`
}

func getURL(text string, old string, new string) string {
	return strings.Replace(text, old, new, -1)
}

func getParsedTimeUnit(unit string) (float64, error) {

	if unit == "" {
		return 0.0, nil
	}

	return strconv.ParseFloat(unit, 64)
}

//GetDashFormats gives the dash formats for any given dash URL
func GetDashFormats(data []byte, masterPlaybackURL string) map[string]map[string]map[string]string {
	var mpd MPD
	var totalSeconds float64
	var totalSegments int
	var audioOrVideo = make(map[string]map[string]map[string]string)
	xml.Unmarshal(data, &mpd)
	mediaPresentationDurationRegex := regexp.MustCompile(`PT((\d+)H)?((\d+)M)?((\d+)\.(\d+)S)?`)
	matches := mediaPresentationDurationRegex.FindAllStringSubmatch(mpd.MediaPresentationDuration, -1)

	if len(matches) > 0 {

		hours, _ := getParsedTimeUnit(matches[0][2])
		minutes, _ := getParsedTimeUnit(matches[0][4])
		seconds, _ := getParsedTimeUnit(matches[0][6])
		milliSeconds, _ := getParsedTimeUnit(matches[0][7])

		totalSeconds = (hours * 60 * 60) + (minutes * 60) + seconds + milliSeconds/1000

		for _, adaptationSet := range mpd.Period {
			switch adaptationSet.MimeType {
			case "video/mp4":
				duration, _ := strconv.ParseFloat(adaptationSet.SegTemplate.Duration, 64)
				timeScale, _ := strconv.ParseFloat(adaptationSet.SegTemplate.Timescale, 64)
				segmentScale := duration / timeScale
				totalSegments = int(math.Ceil(totalSeconds / segmentScale))
				audioOrVideo["video"] = make(map[string]map[string]string)
				var initializationURL = adaptationSet.SegTemplate.Initialization
				var mediaURL = adaptationSet.SegTemplate.Media
				for _, representation := range adaptationSet.Representations {
					var format = make(map[string]string)
					bandwidth, _ := strconv.Atoi(representation.Bandwidth)
					format["BANDWIDTH"] = fmt.Sprintf("%d", bandwidth)
					format["K-FORM"] = fmt.Sprintf("DASH video %dk", bandwidth/1000)
					format["K-FORM-NUMBER"] = fmt.Sprintf("%d", bandwidth/1000)
					format["CODECS"] = fmt.Sprintf("mp4_dash container, %s", representation.Codecs)
					format["RESOLUTION"] = fmt.Sprintf("%sx%s", representation.Width, representation.Height)
					format["FRAME-RATE"] = fmt.Sprintf("%s", representation.FrameRate)
					format["MIME-TYPE"] = adaptationSet.MimeType
					format["STREAM"] = "video only"
					format["TOTAL-SEGMENTS"] = fmt.Sprintf("%d", totalSegments)
					format["INIT-URL"] = getURL(initializationURL, "$RepresentationID$", representation.ID)
					format["STREAM-URL"] = getURL(mediaURL, "$RepresentationID$", representation.ID)
					format["PLAYBACK-URL"] = masterPlaybackURL
					audioOrVideo["video"][fmt.Sprintf("%dk", bandwidth/1000)] = format
				}
			case "audio/mp4":
				duration, _ := strconv.ParseFloat(adaptationSet.SegTemplate.Duration, 64)
				timeScale, _ := strconv.ParseFloat(adaptationSet.SegTemplate.Timescale, 64)
				segmentScale := duration / timeScale
				totalSegments = int(math.Ceil(totalSeconds / segmentScale))
				audioOrVideo["audio"] = make(map[string]map[string]string)
				var initializationURL = adaptationSet.SegTemplate.Initialization
				var mediaURL = adaptationSet.SegTemplate.Media
				for _, representation := range adaptationSet.Representations {
					var format = make(map[string]string)
					bandwidth, _ := strconv.Atoi(representation.Bandwidth)
					format["BANDWIDTH"] = fmt.Sprintf("%d", bandwidth)
					format["K-FORM"] = fmt.Sprintf("DASH audio %dk", bandwidth/1000)
					format["K-FORM-NUMBER"] = fmt.Sprintf("%d", bandwidth/1000)
					format["CODECS"] = fmt.Sprintf("m4a_dash container, %s", representation.Codecs)
					format["MIME-TYPE"] = adaptationSet.MimeType
					format["STREAM"] = "audio only"
					format["TOTAL-SEGMENTS"] = fmt.Sprintf("%d", totalSegments)
					format["SAMPLING-RATE"] = fmt.Sprintf("(%s Hz)", representation.AudioSamplingRate)
					format["INIT-URL"] = getURL(initializationURL, "$RepresentationID$", representation.ID)
					format["STREAM-URL"] = getURL(mediaURL, "$RepresentationID$", representation.ID)
					format["PLAYBACK-URL"] = masterPlaybackURL
					audioOrVideo["audio"][fmt.Sprintf("%dk", bandwidth/1000)] = format
				}
			default:
				fmt.Println("Unsupported format")
			}
		}

	}

	return audioOrVideo
}
