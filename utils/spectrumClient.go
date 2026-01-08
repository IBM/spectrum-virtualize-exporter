// Copyright 2021-2024 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

var logger = log.With("component", "spectrum_exporter")
var ExtraLabelNames, ExtraLabelValues []string

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
			logger.Debugln("return existing token updated in 28s")
			return *s.ColCounter, 1
		}

		if needVerify {
			updatePassedMins := time.Since(s.AuthTokenCache.UpdateTime).Minutes()
			if updatePassedMins < 118 {
				/* 			logger.Debugln("Verify existing token")
				   			_, err := s.CallSpectrumAPI("lssystem", false)
				   			if err == nil {
				   				logger.Debugln("Existing token verified successfully")
				   				retVal = 1
				   			} else {
				   				logger.Debugln("Existing token validation failed")
				   			} */
				retVal = 1
			} else {
				logger.Debugf("it's been %.0f minutes since the token update", updatePassedMins)
			}
			if retVal == 1 {
				if s.Hostname == "" {
					s.Hostname = s.AuthTokenCache.Hostname
				}
				logger.Debugf("return cached token updated in %.0f minutes", updatePassedMins)
				return *s.ColCounter, retVal
			}
		}
	}
	// Start to renew auth token
	lc := 0
	for lc = 0; lc < 3; lc++ {
		logger.Debugln("getting authToken for ", s.IpAddress)
		authtoken, err := s.retrieveAuthToken()
		if err != nil {
			logger.Errorf("failed to request auth token for %s, the error is: %v", s.IpAddress, err)
			s.ColCounter.AuthTokenRenewFailureCount++
			return *s.ColCounter, retVal
		}
		logger.Debugln("got new authToken for ", s.IpAddress)

		s.AuthTokenCache.Token = authtoken
		if !s.AuthTokenCache.UpdateTime.IsZero() {
			s.ColCounter.AuthTokenRenewIntervalSeconds = int(time.Since(s.AuthTokenCache.UpdateTime).Seconds())
		}
		s.AuthTokenCache.UpdateTime = time.Now()

		//test to make sure that current auth token is good
		if needVerify {
			logger.Debugln("verify new auth token for ", s.IpAddress)
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
				logger.Infof("token verification failed for %s, re-requesting authtoken....", s.IpAddress)
				lc++
			} else { //auth token verification succeeded
				logger.Debugln("new auth token verified successfully for ", s.IpAddress)
				break
			}
		} else {
			break
		}
	}
	if lc > 2 {
		s.ColCounter.AuthTokenRenewFailureCount++
		logger.Errorf("failed getting auth token for %s, please check network or username and password", s.IpAddress)
		retVal = 0
	} else {
		logger.Debugln("generated new auth token for ", s.IpAddress)
		s.ColCounter.AuthTokenRenewSuccessCount++
		retVal = 1
	}
	return *s.ColCounter, retVal
}

func (s *SpectrumClient) retrieveAuthToken() (authToken string, err error) {
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
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: !s.VerifyCert, MinVersion: tls.VersionTLS12},
	},
		Timeout: 45 * time.Second,
	}
	req, _ := http.NewRequest("POST", requestURL, nil)
	req.Header.Add("X-Auth-Username", s.UserName)
	req.Header.Add("X-Auth-Password", s.Password)
	// req.SetBasicAuth(s.UserName, s.Password)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error connecting to : %s. the error is: %s", requestURL, err.Error())
	}
	defer resp.Body.Close()

	logger.Debugf("response's status code: %v", resp.StatusCode)
	logger.Debugf("response's status: %v", resp.Status)

	respbody, err := io.ReadAll(resp.Body)
	body := string(respbody)
	logger.Debugf("Response Body: %s", body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http status code is %v when accessing URL: %s. Body text is: %s", resp.StatusCode, requestURL, body)
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
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: !s.VerifyCert, MinVersion: tls.VersionTLS12},
	},
		Timeout: 45 * time.Second}
	// New POST request
	req, _ := http.NewRequest("POST", requestURL, nil)
	// header parameters
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", s.AuthTokenCache.Token)
	logger.Debugf("Request %s using token: %s", requestURL, s.AuthTokenCache.Token)
	//var resp *http.Response
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Debugf("error connecting to Spectrum: %s", err.Error())
		return "", fmt.Errorf("error connecting to : %s. the error is: %s", requestURL, err.Error())
	}
	if autoRenewToken && (resp.StatusCode == 401 || resp.StatusCode == 403) {
		logger.Infoln("token is invalid, start to auto renew auth token")
		_, success := s.RenewAuthToken(false)
		if success == 0 {
			return "", fmt.Errorf("failed to auto renew auth token for %s", s.IpAddress)
		}
		logger.Infoln("auto renewed token and retry rest cmd")
		req, _ := http.NewRequest("POST", requestURL, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Auth-Token", s.AuthTokenCache.Token)
		logger.Debugf("Re-request %s using token: %s", requestURL, s.AuthTokenCache.Token)
		resp, err = httpClient.Do(req)
		if err != nil {
			logger.Debugf("error connecting to Spectrum: %s", err.Error())
			return "", fmt.Errorf("error connecting to : %s. the error is: %s", requestURL, err.Error())
		}
	}
	defer resp.Body.Close()
	respbody, err := io.ReadAll(resp.Body)
	body = string(respbody)
	if resp.StatusCode != 200 {
		logger.Debugf("http status code is %v when accessing URL: %s. Body text is: %s", resp.StatusCode, requestURL, body)
		return "", fmt.Errorf("http status code is %v when accessing URL: %s. Body text is: %s", resp.StatusCode, requestURL, body)
	}
	return body, nil
}

func SpectrumLogger() *log.Logger {
	return &logger
}
