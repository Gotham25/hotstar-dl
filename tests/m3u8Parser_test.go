package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func getExpectedFormats1(kForm string, programID string, closedCaptions string, bandwidth string, codecs string, resolution string, streamURL string, bestResolution string, leastResolution string) map[string]string {
	expectedFormats := make(map[string]string)
	expectedFormats["K-FORM"] = kForm
	expectedFormats["PROGRAM-ID"] = programID
	expectedFormats["CLOSED-CAPTIONS"] = closedCaptions
	expectedFormats["BANDWIDTH"] = bandwidth
	expectedFormats["CODECS"] = codecs
	expectedFormats["RESOLUTION"] = resolution
	expectedFormats["STREAM-URL"] = streamURL
	expectedFormats["BEST_RESOLUTION"] = bestResolution
	expectedFormats["LEAST_RESOLUTION"] = leastResolution
	return expectedFormats
}

func getExpectedFormats2(kForm string, averageBandwidth string, bandwidth string, codecs string, resolution string, frameRate string, streamURL string, bestResolution string, leastResolution string) map[string]string {
	expectedFormats := make(map[string]string)
	expectedFormats["K-FORM"] = kForm
	expectedFormats["AVERAGE-BANDWIDTH"] = averageBandwidth
	expectedFormats["BANDWIDTH"] = bandwidth
	expectedFormats["CODECS"] = codecs
	expectedFormats["RESOLUTION"] = resolution
	expectedFormats["FRAME-RATE"] = frameRate
	expectedFormats["STREAM-URL"] = streamURL
	expectedFormats["BEST_RESOLUTION"] = bestResolution
	expectedFormats["LEAST_RESOLUTION"] = leastResolution
	return expectedFormats
}

func getExpectedFormats3(kForm string, averageBandwidth string, bandwidth string, codecs string, resolution string, streamURL string, bestResolution string, leastResolution string) map[string]string {
	expectedFormats := make(map[string]string)
	expectedFormats["K-FORM"] = kForm
	expectedFormats["AVERAGE-BANDWIDTH"] = averageBandwidth
	expectedFormats["BANDWIDTH"] = bandwidth
	expectedFormats["CODECS"] = codecs
	expectedFormats["RESOLUTION"] = resolution
	expectedFormats["STREAM-URL"] = streamURL
	expectedFormats["BEST_RESOLUTION"] = bestResolution
	expectedFormats["LEAST_RESOLUTION"] = leastResolution
	return expectedFormats
}

func getExpectedVideoFormats1() map[string]map[string]string {
	expectedVideoFormats := make(map[string]map[string]string)
	expectedVideoFormats["hls-167"] = getExpectedFormats1("167k", "1", "NONE", "167271", "\"avc1.66.30, mp4a.40.2\"", "320x180", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_0_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "false", "true")
	expectedVideoFormats["hls-327"] = getExpectedFormats1("327k", "1", "NONE", "327344", "\"avc1.66.30, mp4a.40.2\"", "320x180", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_1_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "false", "false")
	expectedVideoFormats["hls-552"] = getExpectedFormats1("552k", "1", "NONE", "552127", "\"avc1.66.30, mp4a.40.2\"", "416x234", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_2_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "false", "false")
	expectedVideoFormats["hls-960"] = getExpectedFormats1("960k", "1", "NONE", "960823", "\"avc1.66.30, mp4a.40.2\"", "640x360", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_3_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "true", "false")
	expectedVideoFormats["hls-1472"] = getExpectedFormats1("1472k", "1", "NONE", "1472714", "\"avc1.66.30, mp4a.40.2\"", "720x404", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_4_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "false", "false")
	expectedVideoFormats["hls-2188"] = getExpectedFormats1("2188k", "1", "NONE", "2188953", "\"avc1.66.30, mp4a.40.2\"", "1280x720", "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/index_5_av.m3u8?null=0&id=AgCdLeTnMSxxookre1yyOZrUVjGsAjTrI2jaZKKjKzRKekEWQ81I2j3HSzMs2ZZcxJTgLWz%2f4cRk1A%3d%3d&hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", "false", "false")
	return expectedVideoFormats
}

func getExpectedVideoFormats2() map[string]map[string]string {
	expectedVideoFormats := make(map[string]map[string]string)
	expectedVideoFormats["hls-141"] = getExpectedFormats2("141k", "141703", "157168", "\"avc1.42c015,mp4a.40.2\"", "320x180", "15", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer1_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "true")
	expectedVideoFormats["hls-280"] = getExpectedFormats2("280k", "280690", "304560", "\"avc1.42c015,mp4a.40.2\"", "320x180", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer2_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	expectedVideoFormats["hls-505"] = getExpectedFormats2("505k", "505575", "555477", "\"avc1.66.30,mp4a.40.2\"", "416x234", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer3_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	expectedVideoFormats["hls-914"] = getExpectedFormats2("914k", "914365", "1014698", "\"avc1.66.30,mp4a.40.2\"", "640x360", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer4_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "true", "false")
	expectedVideoFormats["hls-1425"] = getExpectedFormats2("1425k", "1425351", "1588474", "\"avc1.66.30,mp4a.40.2\"", "720x404", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer5_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	expectedVideoFormats["hls-2140"] = getExpectedFormats2("2140k", "2140799", "2380832", "\"avc1.42c01f,mp4a.40.2\"", "1280x720", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer6_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	expectedVideoFormats["hls-3297"] = getExpectedFormats2("3297k", "3297345", "3656474", "\"avc1.640029,mp4a.40.2\"", "1600x900", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer7_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	expectedVideoFormats["hls-4830"] = getExpectedFormats2("4830k", "4830306", "5360256", "\"avc1.640032,mp4a.40.2\"", "1920x1080", "25", "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master_Layer8_.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669", "false", "false")
	return expectedVideoFormats
}

func getExpectedVideoFormats3() map[string]map[string]string {
	expectedVideoFormats := make(map[string]map[string]string)
	expectedVideoFormats["hls-178"] = getExpectedFormats3("178k", "178039", "236504", "\"avc1.42C00C,mp4a.40.2\"", "320x180", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-1/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "false", "true")
	expectedVideoFormats["hls-234"] = getExpectedFormats3("234k", "234185", "324488", "\"avc1.42C015,mp4a.40.2\"", "426x240", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-2/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "false", "false")
	expectedVideoFormats["hls-361"] = getExpectedFormats3("361k", "361956", "499704", "\"avc1.4D401E,mp4a.40.2\"", "640x360", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-3/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "true", "false")
	expectedVideoFormats["hls-576"] = getExpectedFormats3("576k", "576455", "877584", "\"avc1.4D401F,mp4a.40.2\"", "854x480", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-4/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "false", "false")
	expectedVideoFormats["hls-1003"] = getExpectedFormats3("1003k", "1003957", "1608904", "\"avc1.4D401F,mp4a.40.2\"", "1280x720", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-5/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "false", "false")
	expectedVideoFormats["hls-1987"] = getExpectedFormats3("1987k", "1987031", "3017776", "\"avc1.640028,mp4a.40.2\"", "1920x1080", "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/media-6/index.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7", "false", "false")
	return expectedVideoFormats
}

func TestParseM3u8Content1(t *testing.T) {

	playbackURL := "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/master.m3u8?hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6"
	playbackURLData := "hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6"

	expectedVideoFormats := getExpectedVideoFormats1()

	m3u8Content, err := ioutil.ReadFile("resources/m3u8Content1.m3u8")
	if err != nil {
		log.Fatal(err)
	}

	actualVideoFormats := utils.ParseM3u8Content(fmt.Sprintf("%s", m3u8Content), playbackURL, playbackURLData)

	if !reflect.DeepEqual(expectedVideoFormats, actualVideoFormats) {
		t.Error("Expected \n", expectedVideoFormats, "\n\n\nbut got \n", actualVideoFormats)
	}

}

func TestParseM3u8Content2(t *testing.T) {
	playbackURL := "https://hsdesinova.akamaized.net/video/vijay_hd/chinnathambi/92df3509e0/337/master.m3u8?hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669"
	playbackURLData := "hdnea=st=1551575720~exp=1551577520~acl=/*~hmac=75f2905ca5d5f79a674205e3e0e25b622ff9d08f77dbc2d50374d70ddb706669"

	expectedVideoFormats := getExpectedVideoFormats2()

	m3u8Content, err := ioutil.ReadFile("resources/m3u8Content2.m3u8")
	if err != nil {
		log.Fatal(err)
	}

	actualVideoFormats := utils.ParseM3u8Content(fmt.Sprintf("%s", m3u8Content), playbackURL, playbackURLData)

	if !reflect.DeepEqual(expectedVideoFormats, actualVideoFormats) {
		t.Error("Expected \n", expectedVideoFormats, "\n\n\nbut got \n", actualVideoFormats)
	}

}

func TestParseM3u8Content3(t *testing.T) {
	playbackURL := "https://hses.akamaized.net/videos/vijay_hd/chinnathambi/0b3c2675ea/362/1100017417/phone/master.m3u8?hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7"
	playbackURLData := "hdnea=st=1551575749~exp=1551577549~acl=/*~hmac=45b40d19a096f5a9e1d0eb68c2c9577ae349443dde273a9ce393f17686badcb7"

	expectedVideoFormats := getExpectedVideoFormats3()

	m3u8Content, err := ioutil.ReadFile("resources/m3u8Content3.m3u8")
	if err != nil {
		log.Fatal(err)
	}

	actualVideoFormats := utils.ParseM3u8Content(fmt.Sprintf("%s", m3u8Content), playbackURL, playbackURLData)

	if !reflect.DeepEqual(expectedVideoFormats, actualVideoFormats) {
		t.Error("Expected \n", expectedVideoFormats, "\n\n\nbut got \n", actualVideoFormats)
	}

}
