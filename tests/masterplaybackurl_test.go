package tests

import (
	"github.com/Gotham25/hotstar-dl/utils"
	"testing"
)

func TestGetMasterPlaybackUrl_ValidStatusCode(t *testing.T) {
	playbackUriContent := `{"body":{"results":{"item":{"playbackUrl":"https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/master.m3u8?hdnea=st=1552212988~exp=1552214788~acl=/*~hmac=32b39cb15ec4c73425126c64044b38983ea45eb8ed0f84fcc9d2ed5eb39a9c5f"},"responseType":"ITEM"}},"statusCode":"OK","statusCodeValue":200}`
	expectedMasterPlaybackUrl := "https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/master.m3u8?hdnea=st=1552212988~exp=1552214788~acl=/*~hmac=32b39cb15ec4c73425126c64044b38983ea45eb8ed0f84fcc9d2ed5eb39a9c5f"

	actualMasterPlaybackUrl, err := utils.GetMasterPlaybackUrl([]byte(playbackUriContent))

	if err == nil && expectedMasterPlaybackUrl != actualMasterPlaybackUrl {
		t.Error("Expected", expectedMasterPlaybackUrl, " but got", actualMasterPlaybackUrl)
	}
}

func TestGetMasterPlaybackUrl_InvalidStatusCode(t *testing.T) {
	playbackUriContent := `{"body":{"results":{"item":{"errorMessage":"Access Denied"},"responseType":"ITEM"}},"statusCode":"ERR","statusCodeValue":401}`
	expectedError := "Invalid status code 401"

	_, actualError := utils.GetMasterPlaybackUrl([]byte(playbackUriContent))

	if expectedError != actualError.Error() {
		t.Error("Expected", expectedError, " but got", actualError)
	}
}
