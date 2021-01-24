package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"text/template"
)

type (
	Client struct {
		LoginName string
		Password  string
	}

	FaultInfo struct {
		FaultCode   string `xml:"Fault>faultcode"`
		FaultString string `xml:"Fault>faultstring"`
	}

	VServerInfo struct {
		IPs      []string `xml:"getVServerInformationResponse>return>ips"`
		Nickname string   `xml:"getVServerInformationResponse>return>vServerNickname"`
		Status   string   `xml:"getVServerInformationResponse>return>status"`
	}

	VServerInformationResponseBody struct {
		XMLName xml.Name `xml:"Body"`
		VServerInfo
		FaultInfo
	}

	VServerInformationResponse struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    VServerInformationResponseBody
	}
)

const netcupWSUrl = "https://www.servercontrolpanel.de:443/SCP/WSEndUser"
const vServerRequestTpl = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:end="http://enduser.service.web.vcp.netcup.de/">
    <soap:Header/>
    <soap:Body>
      <end:{{.Operation}}>
        <loginName>{{.LoginName}}</loginName>
        <password>{{.Password}}</password>
        <vservername>{{.VServerName}}</vservername>
      </end:{{.Operation}}>
    </soap:Body>
  </soap:Envelope>`

func (self *Client) getVServerRequestBody(operation string, vServerName string) (*bytes.Buffer, error) {
	tpl, err := template.New("vServerRequest").Parse(vServerRequestTpl)
	if err != nil {
		return nil, err
	}
	tplData := struct {
		Operation   string
		LoginName   string
		Password    string
		VServerName string
	}{Operation: operation, LoginName: self.LoginName, Password: self.Password, VServerName: vServerName}
	requestBody := bytes.Buffer{}
	err = tpl.Execute(&requestBody, tplData)
	if err != nil {
		return nil, err
	}
	return &requestBody, nil
}

func (self *Client) sendRequest(requestBody *bytes.Buffer, responseData interface{}) error {
	resp, err := http.Post(netcupWSUrl, "text/xml", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = xml.NewDecoder(resp.Body).Decode(responseData)
	if err != nil {
		return err
	}
	return nil
}

func (self *Client) GetVServerInformation(vServerName string) (*VServerInfo, error) {
	requestBody, err := self.getVServerRequestBody("getVServerInformation", vServerName)
	if err != nil {
		return nil, err
	}
	r := new(VServerInformationResponse)
	err = self.sendRequest(requestBody, r)
	if err != nil {
		return nil, err
	}

	if len(r.Body.FaultCode) > 0 {
		return nil, fmt.Errorf("SOAP error: %s [%s]", r.Body.FaultString, r.Body.FaultCode)
	}

	return &r.Body.VServerInfo, nil
}
