package client

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

const error_response = `<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
  <S:Body>
    <S:Fault xmlns:ns4="http://www.w3.org/2003/05/soap-envelope">
      <faultcode>S:Server</faultcode>
      <faultstring>action not allowed</faultstring>
    </S:Fault>
  </S:Body>
</S:Envelope>`

func soapMatcher(user string, password string, serverName string) gock.MatchFunc {
	return func(req *http.Request, ereq *gock.Request) (bool, error) {
		defer req.Body.Close()
		if req.Header.Get("Content-Type") != "text/xml" {
			return false, nil
		}

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return false, err
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
		bodyRegex := regexp.MustCompile(`(?s)<soap:Envelope xmlns:soap="http://schemas\.xmlsoap\.org/soap/envelope/" xmlns:end="http://enduser\.service\.web\.vcp\.netcup\.de/">.*<soap:Body>.*<end:getVServerInformation>.*<loginName>` +
			user + `</loginName>.*<password>` + password + `</password>.*<vservername>` + serverName + `</vservername>.*</end:getVServerInformation>.*</soap:Body>.*</soap:Envelope>`)
		return bodyRegex.Match(b), nil
	}
}

func mockEndpoint(user, password, serverName string) *gock.Request {
	return gock.New("https://www.servercontrolpanel.de:443").
		Post("/SCP/WSEndUser").
		AddMatcher(soapMatcher(user, password, serverName))
}

func TestClient_GetVServerInformation(t *testing.T) {
	Convey("it sends a SOAP request to retrieve server information", t, func() {
		defer gock.Off()
		mockEndpoint("user", "password", "server_name").
			Reply(200).
			Type("text/xml").
			Body(strings.NewReader(`<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns2="http://enduser.service.web.vcp.netcup.de/">
                  <S:Header/>
                    <S:Body>
                      <ns2:getVServerInformationResponse>
                        <return>
                          <ips>1.1.1.1</ips>
                          <ips>2.2.2.2</ips>
                          <status>offline</status>
                          <vServerNickname>nickname</vServerNickname>
                        </return>
                      </ns2:getVServerInformationResponse>
                    </S:Body>
                  </S:Envelope>`))

		client := Client{LoginName: "user", Password: "password"}
		info, err := client.GetVServerInformation("server_name")

		So(err, ShouldBeNil)
		So(info.IPs, ShouldResemble, []string{"1.1.1.1", "2.2.2.2"})
		So(info.Status, ShouldEqual, "offline")
		So(info.Nickname, ShouldEqual, "nickname")
	})

	Convey("it handles error responses", t, func() {
		defer gock.Off()
		mockEndpoint("user", "password", "server_name").
			Reply(500).
			Type("text/xml").
			Body(strings.NewReader(error_response))

		client := Client{LoginName: "user", Password: "password"}
		info, err := client.GetVServerInformation("server_name")

		So(info, ShouldBeNil)
		So(err.Error(), ShouldEqual, "SOAP error: action not allowed [S:Server]")
	})
}
