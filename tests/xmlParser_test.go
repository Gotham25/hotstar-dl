package tests

import (
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func getExpectedDashFormats3() map[string]map[string]map[string]string {

	return map[string]map[string]map[string]string{

		"audio": {
			"33k": {
				"MIME-TYPE":      "audio/mp4",
				"STREAM":         "audio only",
				"STREAM-URL":     "audio/und/mp4a/1/seg-$Number$.m4s",
				"BANDWIDTH":      "33763",
				"K-FORM-NUMBER":  "33",
				"TOTAL-SEGMENTS": "317",
				"SAMPLING-RATE":  "(48000 Hz)",
				"INIT-URL":       "audio/und/mp4a/1/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM":         "DASH audio 33k",
				"CODECS":         "m4a_dash container, mp4a.40.2",
			},
			"65k": {
				"STREAM":         "audio only",
				"TOTAL-SEGMENTS": "317",
				"STREAM-URL":     "audio/und/mp4a/2/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"INIT-URL":       "audio/und/mp4a/2/init.mp4",
				"BANDWIDTH":      "65654",
				"K-FORM":         "DASH audio 65k",
				"K-FORM-NUMBER":  "65",
				"CODECS":         "m4a_dash container, mp4a.40.2",
				"MIME-TYPE":      "audio/mp4",
				"SAMPLING-RATE":  "(48000 Hz)",
			},
		},

		"video": {
			"96k": {
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/1/init.mp4",
				"K-FORM-NUMBER":  "96",
				"RESOLUTION":     "320x180",
				"CODECS":         "mp4_dash container, avc1.42C00C",
				"FRAME-RATE":     "25",
				"STREAM-URL":     "video/avc1/1/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "96191",
				"K-FORM":         "DASH video 96k",
			},
			"163k": {
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"INIT-URL":       "video/avc1/2/init.mp4",
				"STREAM-URL":     "video/avc1/2/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM-NUMBER":  "163",
				"CODECS":         "mp4_dash container, avc1.42C015",
				"RESOLUTION":     "426x240",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"BANDWIDTH":      "163093",
				"K-FORM":         "DASH video 163k",
			},
			"242k": {
				"BANDWIDTH":      "242217",
				"CODECS":         "mp4_dash container, avc1.4D401E",
				"FRAME-RATE":     "25",
				"TOTAL-SEGMENTS": "317",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM":         "DASH video 242k",
				"K-FORM-NUMBER":  "242",
				"RESOLUTION":     "640x360",
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"INIT-URL":       "video/avc1/3/init.mp4",
				"STREAM-URL":     "video/avc1/3/seg-$Number$.m4s",
			},
			"475k": {
				"K-FORM":         "DASH video 475k",
				"K-FORM-NUMBER":  "475",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/4/init.mp4",
				"STREAM-URL":     "video/avc1/4/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "475884",
				"RESOLUTION":     "854x480",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"CODECS":         "mp4_dash container, avc1.4D401F",
			},
			"822k": {
				"BANDWIDTH":      "822677",
				"K-FORM-NUMBER":  "822",
				"CODECS":         "mp4_dash container, avc1.640028",
				"MIME-TYPE":      "video/mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM":         "DASH video 822k",
				"RESOLUTION":     "1280x720",
				"FRAME-RATE":     "25",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/5/init.mp4",
				"STREAM-URL":     "video/avc1/5/seg-$Number$.m4s",
			},
			"1847k": {
				"K-FORM-NUMBER":  "1847",
				"RESOLUTION":     "1920x1080",
				"MIME-TYPE":      "video/mp4",
				"INIT-URL":       "video/avc1/6/init.mp4",
				"STREAM-URL":     "video/avc1/6/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "1847477",
				"CODECS":         "mp4_dash container, avc1.640028",
				"FRAME-RATE":     "25",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"K-FORM":         "DASH video 1847k",
			},
		},
	}

}

func getExpectedDashFormats2() map[string]map[string]map[string]string {

	return map[string]map[string]map[string]string{

		"audio": {
			"33k": {
				"BANDWIDTH":      "33763",
				"K-FORM-NUMBER":  "33",
				"CODECS":         "m4a_dash container, mp4a.40.2",
				"MIME-TYPE":      "audio/mp4",
				"STREAM":         "audio only",
				"SAMPLING-RATE":  "(48000 Hz)",
				"INIT-URL":       "audio/und/mp4a/1/init.mp4",
				"K-FORM":         "DASH audio 33k",
				"TOTAL-SEGMENTS": "317",
				"STREAM-URL":     "audio/und/mp4a/1/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
			},
			"65k": {
				"K-FORM":         "DASH audio 65k",
				"MIME-TYPE":      "audio/mp4",
				"STREAM":         "audio only",
				"TOTAL-SEGMENTS": "317",
				"BANDWIDTH":      "65654",
				"K-FORM-NUMBER":  "65",
				"CODECS":         "m4a_dash container, mp4a.40.2",
				"SAMPLING-RATE":  "(48000 Hz)",
				"INIT-URL":       "audio/und/mp4a/2/init.mp4",
				"STREAM-URL":     "audio/und/mp4a/2/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
			},
		},

		"video": {
			"171k": {
				"MIME-TYPE":      "video/mp4",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/1/init.mp4",
				"STREAM-URL":     "video/avc1/1/seg-$Number$.m4s",
				"BANDWIDTH":      "171790",
				"K-FORM":         "DASH video 171k",
				"K-FORM-NUMBER":  "171",
				"RESOLUTION":     "426x240",
				"CODECS":         "mp4_dash container, avc1.42C015",
				"FRAME-RATE":     "25",
				"STREAM":         "video only",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
			},
			"311k": {
				"INIT-URL":       "video/avc1/2/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM-NUMBER":  "311",
				"STREAM":         "video only",
				"CODECS":         "mp4_dash container, avc1.4D401E",
				"RESOLUTION":     "640x360",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"TOTAL-SEGMENTS": "317",
				"STREAM-URL":     "video/avc1/2/seg-$Number$.m4s",
				"BANDWIDTH":      "311810",
				"K-FORM":         "DASH video 311k",
			},
			"566k": {
				"CODECS":         "mp4_dash container, avc1.4D401F",
				"RESOLUTION":     "854x480",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"STREAM-URL":     "video/avc1/3/seg-$Number$.m4s",
				"BANDWIDTH":      "566395",
				"K-FORM":         "DASH video 566k",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/3/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM-NUMBER":  "566",
				"STREAM":         "video only",
			},
			"1074k": {
				"K-FORM":         "DASH video 1074k",
				"RESOLUTION":     "1280x720",
				"STREAM-URL":     "video/avc1/4/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"INIT-URL":       "video/avc1/4/init.mp4",
				"BANDWIDTH":      "1074338",
				"K-FORM-NUMBER":  "1074",
				"CODECS":         "mp4_dash container, avc1.640028",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
			},
			"2408k": {
				"K-FORM-NUMBER":  "2408",
				"RESOLUTION":     "1920x1080",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/5/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "2408914",
				"K-FORM":         "DASH video 2408k",
				"CODECS":         "mp4_dash container, avc1.640028",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"STREAM-URL":     "video/avc1/5/seg-$Number$.m4s",
			},
		},
	}

}

func getExpectedDashFormats1() map[string]map[string]map[string]string {
	return map[string]map[string]map[string]string{
		"audio": {
			"33k": {
				"BANDWIDTH":      "33763",
				"K-FORM":         "DASH audio 33k",
				"TOTAL-SEGMENTS": "317",
				"SAMPLING-RATE":  "(48000 Hz)",
				"STREAM-URL":     "audio/und/mp4a/1/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM-NUMBER":  "33",
				"CODECS":         "m4a_dash container, mp4a.40.2",
				"MIME-TYPE":      "audio/mp4",
				"STREAM":         "audio only",
				"INIT-URL":       "audio/und/mp4a/1/init.mp4",
			},
			"65k": {
				"BANDWIDTH":      "65654",
				"K-FORM-NUMBER":  "65",
				"MIME-TYPE":      "audio/mp4",
				"STREAM":         "audio only",
				"INIT-URL":       "audio/und/mp4a/2/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM":         "DASH audio 65k",
				"CODECS":         "m4a_dash container, mp4a.40.2",
				"TOTAL-SEGMENTS": "317",
				"SAMPLING-RATE":  "(48000 Hz)",
				"STREAM-URL":     "audio/und/mp4a/2/seg-$Number$.m4s",
			},
		},

		"video": {
			"96k": {
				"BANDWIDTH":      "96191",
				"K-FORM-NUMBER":  "96",
				"CODECS":         "mp4_dash container, avc1.42C00C",
				"STREAM":         "video only",
				"INIT-URL":       "video/avc1/1/init.mp4",
				"K-FORM":         "DASH video 96k",
				"RESOLUTION":     "320x180",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"TOTAL-SEGMENTS": "317",
				"STREAM-URL":     "video/avc1/1/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
			},
			"163k": {
				"BANDWIDTH":      "163093",
				"K-FORM-NUMBER":  "163",
				"RESOLUTION":     "426x240",
				"TOTAL-SEGMENTS": "317",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"STREAM-URL":     "video/avc1/2/seg-$Number$.m4s",
				"K-FORM":         "DASH video 163k",
				"CODECS":         "mp4_dash container, avc1.42C015",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"INIT-URL":       "video/avc1/2/init.mp4",
			},
			"242k": {
				"STREAM":         "video only",
				"STREAM-URL":     "video/avc1/3/seg-$Number$.m4s",
				"K-FORM":         "DASH video 242k",
				"RESOLUTION":     "640x360",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"INIT-URL":       "video/avc1/3/init.mp4",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "242217",
				"K-FORM-NUMBER":  "242",
				"CODECS":         "mp4_dash container, avc1.4D401E",
				"TOTAL-SEGMENTS": "317",
			},
			"475k": {
				"MIME-TYPE":      "video/mp4",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"BANDWIDTH":      "475884",
				"K-FORM":         "DASH video 475k",
				"K-FORM-NUMBER":  "475",
				"CODECS":         "mp4_dash container, avc1.4D401F",
				"RESOLUTION":     "854x480",
				"FRAME-RATE":     "25",
				"INIT-URL":       "video/avc1/4/init.mp4",
				"STREAM-URL":     "video/avc1/4/seg-$Number$.m4s",
			},
			"822k": {
				"K-FORM":         "DASH video 822k",
				"K-FORM-NUMBER":  "822",
				"CODECS":         "mp4_dash container, avc1.640028",
				"FRAME-RATE":     "25",
				"MIME-TYPE":      "video/mp4",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/5/init.mp4",
				"BANDWIDTH":      "822677",
				"RESOLUTION":     "1280x720",
				"STREAM":         "video only",
				"STREAM-URL":     "video/avc1/5/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
			},
			"1847k": {
				"STREAM-URL":     "video/avc1/6/seg-$Number$.m4s",
				"PLAYBACK-URL":   "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d",
				"K-FORM-NUMBER":  "1847",
				"CODECS":         "mp4_dash container, avc1.640028",
				"MIME-TYPE":      "video/mp4",
				"FRAME-RATE":     "25",
				"STREAM":         "video only",
				"TOTAL-SEGMENTS": "317",
				"INIT-URL":       "video/avc1/6/init.mp4",
				"BANDWIDTH":      "1847477",
				"K-FORM":         "DASH video 1847k",
				"RESOLUTION":     "1920x1080",
			},
		},
	}
}

func TestGetDashFormats1(t *testing.T) {
	masterPlaybackURL := "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d"
	mpdContent, err := ioutil.ReadFile("resources/mpdContent1.xml")

	if err != nil {
		log.Fatal(err)
	}

	expectedDashFormats := getExpectedDashFormats1()
	actualDashFormats := utils.GetDashFormats(mpdContent, masterPlaybackURL)

	if !reflect.DeepEqual(expectedDashFormats, actualDashFormats) {
		t.Error("Expected \n", expectedDashFormats, "\n\n\nbut got \n", actualDashFormats)
	}
}

func TestGetDashFormats2(t *testing.T) {
	masterPlaybackURL := "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d"
	mpdContent, err := ioutil.ReadFile("resources/mpdContent2.xml")

	if err != nil {
		log.Fatal(err)
	}

	expectedDashFormats := getExpectedDashFormats2()
	actualDashFormats := utils.GetDashFormats(mpdContent, masterPlaybackURL)

	if !reflect.DeepEqual(expectedDashFormats, actualDashFormats) {
		t.Error("Expected \n", expectedDashFormats, "\n\n\nbut got \n", actualDashFormats)
	}
}

func TestGetDashFormats3(t *testing.T) {
	masterPlaybackURL := "https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566668622~exp=1566672222~acl=/*~hmac=c880125b7a047e7406d97ca30d11a95504773a5ce5bf053ceef748118ce2986d"
	mpdContent, err := ioutil.ReadFile("resources/mpdContent3.xml")

	if err != nil {
		log.Fatal(err)
	}

	expectedDashFormats := getExpectedDashFormats3()
	actualDashFormats := utils.GetDashFormats(mpdContent, masterPlaybackURL)

	if !reflect.DeepEqual(expectedDashFormats, actualDashFormats) {
		t.Error("Expected \n", expectedDashFormats, "\n\n\nbut got \n", actualDashFormats)
	}
}
