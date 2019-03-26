package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//ParseM3u8Content parses given m3u8Content content and returns map of map of string containing video format list.
func ParseM3u8Content(m3u8Content string, playbackUrl string, playbackUrlData string) map[string]map[string]string {

	var m3u8Info map[string]string
	var urlFormats = make(map[string]map[string]string)
	for _, line := range strings.Split(m3u8Content, "\n") {

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {

			if m3u8Info == nil {
				m3u8Info = make(map[string]string)
			}

			m3u8InfoCsv := strings.Replace(line, "#EXT-X-STREAM-INF:", "", -1)
			m3u8InfoRegex := regexp.MustCompile(`([\w\-]+)\=([\w\-]+|"[^"]*")`)

			for _, info := range m3u8InfoRegex.FindAllStringSubmatch(m3u8InfoCsv, -1) {
				m3u8Info[info[1]] = info[2]
			}
		} else if strings.Contains(line, ".m3u8") {

			if m3u8Info != nil {

				averageBandwidthOrBandwidth := func() int {
					var bw string
					if m3u8Info["AVERAGE-BANDWIDTH"] != "" {
						bw = m3u8Info["AVERAGE-BANDWIDTH"]
					} else {
						bw = m3u8Info["BANDWIDTH"]
					}
					var bwInt, _ = strconv.Atoi(bw)
					return bwInt
				}()

				kFactor := averageBandwidthOrBandwidth / 1000
				kForm := fmt.Sprintf("%dk", kFactor)

				m3u8Info["K-FORM"] = kForm

				streamUrl := line

				if !strings.HasPrefix(line, "http") {
					streamUrl = strings.Replace(playbackUrl, "master.m3u8", line, -1)
				}

				if !strings.Contains(streamUrl, "~acl=/*~hmac") {
					if !strings.Contains(streamUrl, "?") {
						streamUrl += "?"
					}
					streamUrl += ("&" + playbackUrlData)
				}

				re := regexp.MustCompile(`\r`)
				streamUrl = re.ReplaceAllString(streamUrl, "")

				m3u8Info["STREAM-URL"] = streamUrl

				urlFormats[fmt.Sprintf("hls-%d", kFactor)] = CopyMap(m3u8Info)

				//Reset m3u8InfoArray for next layer
				m3u8Info = nil
			}

		} else {
			//do nothing
		}
	}

	return urlFormats
}
