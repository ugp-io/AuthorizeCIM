package AuthorizeCIM

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var api_endpoint string = "https://apitest.authorize.net/xml/v1/request.api"
var apiName *string
var apiKey *string
var testMode string
var showLogs bool = true
var connected bool = false

func SetAPIInfo(name string, key string, mode string) {
	apiKey = &key
	apiName = &name
	if mode == "live" {
		showLogs = false
		testMode = "liveMode"
		api_endpoint = "https://api.authorize.net/xml/v1/request.api"
	} else {
		showLogs = false
		testMode = "testMode"
		api_endpoint = "https://apitest.authorize.net/xml/v1/request.api"
	}
}

func IsConnected() (bool, error) {
	info, err := GetMerchantDetails()
	if err != nil {
		return false, err
	}
	if info.Ok() {
		return true, err
	}
	return false, err
}

func GetAuthentication() MerchantAuthentication {
	auth := MerchantAuthentication{
		Name:           apiName,
		TransactionKey: apiKey,
	}
	return auth
}

func SendRequest(input []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", api_endpoint, bytes.NewBuffer(input))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	if showLogs {
		fmt.Println(string(body))
	}
	return body, err
}

func (r AVS) Text() string {
	var response string
	switch r.avsResultCode {
	case "E":
		response = "AVS Error"
	case "R":
		response = "Retry, System Is Unavailable"
	case "G":
		response = "Non U.S. Card Issuing Bank"
	case "U":
		response = "Address Information For This Cardholder Is Unavailable"
	case "S":
		response = "AVS Not Supported by Card Issuing Bank"
	case "N":
		response = "Street Address: No Match - First 5 Digits of ZIP: No Match"
	case "A":
		response = "Street Address: Match - First 5 Digits of ZIP: No Match"
	case "B":
		response = "Address not provided for AVS check or street address match, postal code could not be verified"
	case "P":
		response = "AVS not applicable for this transaction"
	case "Z":
		response = "Street Address: No Match - First 5 Digits of ZIP: Match"
	case "W":
		response = "Street Address: No Match - All 9 Digits of ZIP: Match"
	case "X":
		response = "Street Address: Match - All 9 Digits of ZIP: Match"
	case "Y":
		response = "Street Address: Match - First 5 Digits of ZIP: Match"
	}
	return response
}

func (r CVV) Text() string {
	var response string
	switch r.cvvResultCode {
	case "N":
		response = "CVV Does Not Match"
	case "S":
		response = "CVV Should be on the card, but is not indicated"
	case "U":
		response = "The issuer is not certified for CVV processing or has not provided an encryption key"
	case "P":
		response = "CVV Is not processed"
	case "M":
		response = "CVV Matched"
	}
	return response
}
