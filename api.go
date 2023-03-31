package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func extractDomainFromUri(uri string) string {
	// extract domain from uri
	parsedURL, err := url.Parse(uri)
	ret := ""
	if err == nil {
		ret = parsedURL.Hostname()
	}
	return ret
}

func callApi(method string, entity string, name string, id string, data interface{}) (string, error) {
	// create Jwt token
	token, err := GenerateJWTAccessToken(extractDomainFromUri(shieldooUri))
	if err != nil {
		return "", err
	}
	// call REST API
	httpClient := &http.Client{}
	myurl := shieldooUri + "/cliapi/" + entity
	if id != "" {
		myurl += "/" + url.QueryEscape(id)
	}
	if name != "" {
		// url encode name
		myurl += "?name=" + url.QueryEscape(name)
	}
	var buff *bytes.Buffer = nil
	// convert data to json if it is not nil
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		buff = bytes.NewBuffer(jsonData)
	} else {
		buff = &bytes.Buffer{}
	}
	req, err := http.NewRequest(method, myurl, buff)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AuthToken", token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return string(body), errors.New(resp.Status)
	}
	return strings.TrimSpace(string(body)), nil
}
