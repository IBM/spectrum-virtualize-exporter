package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

type SpectrumClient struct {
	UserName       string
	Password       string
	IpAddress      string
	ErrorCount     float64
	Hostname       string
	VerifyCert     bool
	AuthTokenCache *AuthToken
	AuthTokenMutex *sync.Mutex
	ColCounter     *Counter //shared cross all SpectrumClients of a target
}

type AuthToken struct {
	Token      string
	Hostname   string
	UpdateTime time.Time
}

type Counter struct {
	AuthTokenRenewIntervalSeconds int
	AuthTokenRenewSuccessCount    int
	AuthTokenRenewFailureCount    int
}

func (s *SpectrumClient) RenewAuthToken(needVerify bool) (Counter, int) {
	defer s.AuthTokenMutex.Unlock()
	s.AuthTokenMutex.Lock()

	// A single session lasts a maximum of two active hours or thirty inactive minutes, whichever occurs first.
	retVal := 0 // 0: failed, 1: success
	if s.AuthTokenCache.Token != "" {
		if time.Since(s.AuthTokenCache.UpdateTime).Seconds() < 28 {
			log.Debugln("Return existing token updated in 28s")
			return *s.ColCounter, 1
		}

		if needVerify {
			updatePassedMins := time.Since(s.AuthTokenCache.UpdateTime).Minutes()
			if updatePassedMins < 118 {
				/* 			log.Debugln("Verify existing token")
				   			_, err := s.CallSpectrumAPI("lssystem", false)
				   			if err == nil {
				   				log.Debugln("Existing token verified successfully")
				   				retVal = 1
				   			} else {
				   				log.Debugln("Existing token validation failed")
				   			} */
				retVal = 1
			} else {
				log.Debugf("It's been %.0f minutes since the token update", updatePassedMins)
			}
			if retVal == 1 {
				if s.Hostname == "" {
					s.Hostname = s.AuthTokenCache.Hostname
				}
				log.Debugf("Return cached token updated in %.0f minutes", updatePassedMins)
				return *s.ColCounter, retVal
			}
		}
	}
	// Start to renew auth token
	lc := 0
	for lc = 0; lc < 3; lc++ {
		log.Debugln("Getting authToken for ", s.IpAddress)
		authtoken, err := s.retriveAuthToken()
		if err != nil {
			log.Errorf("Failed to request auth token for %s, the error is: %v.", s.IpAddress, err)
			s.ColCounter.AuthTokenRenewFailureCount++
			return *s.ColCounter, retVal
		}
		log.Debugln("Got new authToken for ", s.IpAddress)

		s.AuthTokenCache.Token = authtoken
		if !s.AuthTokenCache.UpdateTime.IsZero() {
			s.ColCounter.AuthTokenRenewIntervalSeconds = int(time.Since(s.AuthTokenCache.UpdateTime).Seconds())
		}
		s.AuthTokenCache.UpdateTime = time.Now()

		//test to make sure that current auth token is good
		if needVerify {
			log.Debugln("Verify new auth token for ", s.IpAddress)
			i := 0
			for i < 2 {
				systemMetrics, err := s.CallSpectrumAPI("lssystem", false)
				if err != nil {
					if i == 0 {
						time.Sleep(2 * time.Second)
					}
					i++
				} else {
					//We have a valid auth token, we can break out of this loop
					if s.Hostname == "" {
						s.Hostname = gjson.Get(systemMetrics, "name").String()
						s.AuthTokenCache.Hostname = s.Hostname
					}
					break
				}
			}
			if i > 1 { //auth token verification failed
				s.AuthTokenCache.Token = ""
				log.Infof("\nToken verification failed for %s, re-requesting authtoken....", s.IpAddress)
				lc++
			} else { //auth token verification succeeded
				log.Debugln("New auth token verified successfully for ", s.IpAddress)
				break
			}
		} else {
			break
		}
	}
	if lc > 2 {
		s.ColCounter.AuthTokenRenewFailureCount++
		log.Errorf("Failed getting auth token for %s, please check network or username and password", s.IpAddress)
		retVal = 0
	} else {
		log.Debugln("Generated new auth token for ", s.IpAddress)
		s.ColCounter.AuthTokenRenewSuccessCount++
		retVal = 1
	}
	return *s.ColCounter, retVal
}

func (s *SpectrumClient) retriveAuthToken() (authToken string, err error) {
	requestURL := "https://" + s.IpAddress + ":7443/rest/auth"
	httpClient := &http.Client{Transport: &http.Transport{
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
	log.Debugf("Skip verifying the server cert: %v", !s.VerifyCert)
	req, _ := http.NewRequest("POST", requestURL, nil)
	req.Header.Add("X-Auth-Username", s.UserName)
	req.Header.Add("X-Auth-Password", s.Password)
	// req.SetBasicAuth(s.UserName, s.Password)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("\nError connecting to : %s. the error is: %s", requestURL, err.Error())
	}
	defer resp.Body.Close()

	log.Debugf("Response Status Code: %v", resp.StatusCode)
	log.Debugf("Response Status: %v", resp.Status)

	respbody, err := ioutil.ReadAll(resp.Body)
	body := string(respbody)
	log.Debugf("Response Body: %s", body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("\nhttp status code is %v when accessing URL: %s\n Body text is: %s", resp.StatusCode, requestURL, body)
	}
	authToken = gjson.Get(body, "token").String()
	return authToken, nil
}

func (s *SpectrumClient) CallSpectrumAPI(restCmd string, autoRenewToken bool) (body string, err error) {
	requestURL := "https://" + s.IpAddress + ":7443/rest/" + restCmd
	httpClient := &http.Client{Transport: &http.Transport{
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
	req, _ := http.NewRequest("POST", requestURL, nil)
	// header parameters
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", s.AuthTokenCache.Token)
	log.Debugf("Request %s using token: %s", requestURL, s.AuthTokenCache.Token)
	//var resp *http.Response
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Debugf("\n - Error connecting to Spectrum: %s", err.Error())
		return "", fmt.Errorf("\nError connecting to : %s. the error is: %s", requestURL, err.Error())
	}
	if autoRenewToken && resp.StatusCode == 403 {
		log.Infoln("Token is invalid, start to auto renew auth token.")
		_, success := s.RenewAuthToken(false)
		if success == 0 {
			return "", fmt.Errorf("\nFailed to auto renew auth token for %s", s.IpAddress)
		}
		log.Infoln("Auto renewed token and retry rest cmd.")
		req, _ := http.NewRequest("POST", requestURL, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Auth-Token", s.AuthTokenCache.Token)
		log.Debugf("Re-request %s using token: %s", requestURL, s.AuthTokenCache.Token)
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Debugf("\n - Error connecting to Spectrum: %s", err.Error())
			return "", fmt.Errorf("\nError connecting to : %s. the error is: %s", requestURL, err.Error())
		}
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	body = string(respbody)
	if resp.StatusCode != 200 {
		log.Debugf("\nhttp status code is %v when accessing URL: %s\n Body text is: %s\n", resp.StatusCode, requestURL, body)
		return "", fmt.Errorf("\nhttp status code is %v when accessing URL: %s\n Body text is: %s", resp.StatusCode, requestURL, body)
	}
	return body, nil
}
