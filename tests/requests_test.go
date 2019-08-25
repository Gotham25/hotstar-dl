package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Gotham25/hotstar-dl/utils"
)

func TestMakeGetRequest_ValidURL(t *testing.T) {
	expectedPageContents := "google-site-verification: google80369ee5f1cdb5f7.html"

	pageContents, err := utils.MakeGetRequest("https://hotstardownload.herokuapp.com/google80369ee5f1cdb5f7.html", nil)
	actualPageContents := fmt.Sprintf("%s", pageContents)

	if err == nil && expectedPageContents != actualPageContents {
		t.Error("Expected", expectedPageContents, " but got", actualPageContents)
	}
}

//TODO: Rework this tests
/*func TestMakeGetRequest_InvalidURL(t *testing.T) {
	expectedError := "Get https://www.blah.com: net/http: TLS handshake timeout"

	pageContents, actualError := utils.MakeGetRequest("https://www.blah.com", nil)

	if pageContents == nil && actualError != nil && expectedError != actualError.Error() {
		t.Error("Expected", expectedError, " but got", actualError.Error())
	}
}*/

func TestMakeGetRequest_ValidURL_ReturnsInvalidStatusCode(t *testing.T) {
	var expectedPageContents string
	expectedError := "Invalid response code: 403"
	expectedPageContents += "<HTML><HEAD>\n"
	expectedPageContents += "<TITLE>Access Denied</TITLE>\n"
	expectedPageContents += "</HEAD><BODY>\n"
	expectedPageContents += "<H1>Access Denied</H1>\n"
	expectedPageContents += " \n"
	expectedPageContents += "You don't have permission to access \"http&#58;&#47;&#47;hssouthsp&#45;vh&#46;akamaihd&#46;net&#47;i&#47;videos&#47;vijay&#95;hd&#47;chinnathambi&#47;149&#47;master&#95;&#44;106&#44;180&#44;400&#44;800&#44;1300&#44;2000&#44;3000&#44;4500&#44;kbps&#46;mp4&#46;csmil&#47;master&#46;m3u8&#63;\" on this server.<P>\n"
	expectedPageContents += "</BODY>\n"
	expectedPageContents += "</HTML>\n"

	pageContentsBytes, actualError := utils.MakeGetRequest("https://hssouthsp-vh.akamaihd.net/i/videos/vijay_hd/chinnathambi/149/master_,106,180,400,800,1300,2000,3000,4500,kbps.mp4.csmil/master.m3u8?hdnea=st=1551575624~exp=1551577424~acl=/*~hmac=3d89f2aab02315ee100156209746e0e9f3bc70b0b52c17573300b5caa517cfd6", nil)
	pageContentsString := fmt.Sprintf("%s", pageContentsBytes)

	actualPageContents := fmt.Sprintf("%s%s", pageContentsBytes[:strings.Index(pageContentsString, "Reference")], pageContentsBytes[strings.Index(pageContentsString, "</BODY>"):len(pageContentsString)])

	if pageContentsBytes == nil || actualError == nil || expectedPageContents != actualPageContents || expectedError != actualError.Error() {
		t.Error("Expected", expectedError, " but got", actualError.Error())
	}
}
