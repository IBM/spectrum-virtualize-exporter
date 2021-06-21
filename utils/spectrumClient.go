package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

type SpectrumClient struct {
	UserName   string
	Password   string
	AuthToken  string
	IpAddress  string
	ErrorCount float64
	Hostname   string
	VerifyCert bool
}

func (s *SpectrumClient) RetriveAuthToken() (authToken string, err error) {
	reqAuthURL := "https://" + s.IpAddress + ":7443/rest/auth"
	httpclient := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: !s.VerifyCert},
	},
		Timeout: 45 * time.Second,
	}
	log.Debug("Skip verifying the server cert: %v", !s.VerifyCert)
	req, _ := http.NewRequest("POST", reqAuthURL, nil)
	req.Header.Add("X-Auth-Username", s.UserName)
	req.Header.Add("X-Auth-Password", s.Password)
	// req.SetBasicAuth(s.UserName, s.Password)
	resp, err := httpclient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error doing http request URL[%s] Error: %v", reqAuthURL, err)
		return
	}
	defer resp.Body.Close()

	log.Debugf("Response Status Code: %v", resp.StatusCode)
	log.Debugf("Response Status: %v", resp.Status)
	log.Debugf("Response Body: %v", resp.Body)

	if resp.StatusCode != 200 {
		// we didnt get a good response code, so bailing out
		log.Errorln("Got a non 200 response code: ", resp.StatusCode)
		log.Debugln("response was: ", resp)
		s.ErrorCount++
		return "", fmt.Errorf("received non 200 error code: %v. the response was: %v", resp.Status, resp)
	}
	respbody, err := ioutil.ReadAll(resp.Body)
	body := string(respbody)
	authToken = gjson.Get(body, "token").String()
	log.Debugf("AuthToken is: %v", authToken)
	return authToken, err

}

func (s *SpectrumClient) CallSpectrumAPI(request string) (body string, err error) {
	httpclient := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: !s.VerifyCert},
	},
		Timeout: 45 * time.Second}
	log.Debugf("Skip verifying the server cert: %v", !s.VerifyCert)
	// New POST request
	req, _ := http.NewRequest("POST", request, nil)
	// header parameters
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", s.AuthToken)
	resp, err := httpclient.Do(req)
	if err != nil {
		log.Debugf("\n - Error connecting to Spectrum: %s", err)
		return "", fmt.Errorf("\nError connecting to : %v. the error was: %v", request, err)
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	body = string(respbody)
	if resp.StatusCode != 200 {
		log.Debugf("\nGot error code: %v when accessing URL: %s\n Body text is: %s\n", resp.StatusCode, request, respbody)
		return "", fmt.Errorf("\nGot error code: %v when accessing URL: %s\n Body text is: %s", resp.StatusCode, request, respbody)
	}
	return body, nil

}
