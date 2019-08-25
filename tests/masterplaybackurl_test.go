package tests

import (
	"reflect"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestGetMasterPlaybackURL_ValidStatusCode(t *testing.T) {
	playbackURIContent := `{"body":{"results":{"contentId":"1100025368","requestedConfig":"encryption:plain;ladder:phone,tv;package:hls,dash","drmClass":"BEST_EFFORT","downloadDrmClass":"BEST_EFFORT","match":false,"playBackSets":[{"tagsCombination":"encryption:plain;package:hls","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/ac151726f6cf064074a53a40d32f16e7/master.m3u8?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"},{"tagsCombination":"encryption:plain;ladder:phone;package:dash","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"},{"tagsCombination":"encryption:plain;ladder:phone;package:hls","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/ac151726f6cf064074a53a40d32f16e7/master.m3u8?ladder=phone&hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"},{"tagsCombination":"encryption:plain;ladder:tv;package:dash","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"},{"tagsCombination":"encryption:plain;ladder:tv;package:hls","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/f54186d66b64cbe087441cba561e7dfb/master.m3u8?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"},{"tagsCombination":"encryption:plain;package:dash","playbackUrl":"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a","playbackCDNType":"INTERNAL"}]}},"statusCodeValue":200,"statusCode":"OK"}`
	expectedMasterPlaybackURLs := []string{
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/ac151726f6cf064074a53a40d32f16e7/master.m3u8?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?ladder=phone&hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/ac151726f6cf064074a53a40d32f16e7/master.m3u8?ladder=phone&hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/9c2049fc628eb8669170c2289b7d48e5/master.mpd?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/f54186d66b64cbe087441cba561e7dfb/master.m3u8?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
		"https://hses.akamaized.net/videos/vijay_hd/naam_iruvar_namakku_iruvar/c18c23262c/386/1100025368/1565374162079/69b5fa122ada150073875ff77a52bbee/master.mpd?hdnea=st=1566234025~exp=1566237625~acl=/*~hmac=a65fa84f64a81c40e05c2d5394a25eb80ad83572c69752168998a64a60d9f77a",
	}

	actualMasterPlaybackURLs, err := utils.GetMasterPlaybackURLs([]byte(playbackURIContent))

	if err == nil && !reflect.DeepEqual(expectedMasterPlaybackURLs, actualMasterPlaybackURLs) {
		t.Error("Expected", expectedMasterPlaybackURLs, " but got", actualMasterPlaybackURLs)
	}
}

func TestGetMasterPlaybackURL_InvalidStatusCode(t *testing.T) {
	playbackURIContent := `{"body":{"results":{"item":{"errorMessage":"Access Denied"},"responseType":"ITEM"}},"statusCode":"ERR","statusCodeValue":401}`
	expectedError := "Invalid status code 401"

	_, actualError := utils.GetMasterPlaybackURLs([]byte(playbackURIContent))

	if expectedError != actualError.Error() {
		t.Error("Expected", expectedError, " but got", actualError)
	}
}
