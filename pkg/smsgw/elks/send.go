package elks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/sakjur/telegraf/pkg/smsgw"

	"golang.org/x/text/currency"
)

type elksSmsResponse struct {
	Id   string `json:"id"`
	Cost int    `json:"cost"`
}

type elksProfileResponse struct {
	Currency string `json:"currency"`
}

const (
	smsUrl         = "https://api.46elks.com/a1/sms"
	profileUrl     = "https://api.46elks.com/a1/me"
	elksCostFactor = 10000
)

func Send(msg smsgw.Message) (*smsgw.ApiResponse, error) {
	smsResp, err := sendSms(msg)

	if err != nil {
		return nil, err
	}

	curr, err := getCurrency()

	if err != nil {
		return nil, err
	}

	cost := curr.Amount(float64(smsResp.Cost) / elksCostFactor)

	resp := smsgw.ApiResponse{RemoteId: smsResp.Id, Cost: &cost}

	return &resp, nil
}

func sendSms(msg smsgw.Message) (*elksSmsResponse, error) {
	elksUrl, err := url.Parse(smsUrl)

	if err != nil {
		return nil, err
	}

	err = authenticateUrl(elksUrl)

	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm(elksUrl.String(), encodeMessageForSending(msg))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	apiResponse := elksSmsResponse{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func getCurrency() (*currency.Unit, error) {
	elksUrl, err := url.Parse(profileUrl)

	if err != nil {
		return nil, err
	}

	err = authenticateUrl(elksUrl)

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(elksUrl.String())

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	profileResponse := elksProfileResponse{}
	err = json.Unmarshal(body, &profileResponse)

	if err != nil {
		return nil, err
	}

	switch profileResponse.Currency {
	case "SEK":
		return &currency.SEK, nil
	case "EUR":
		return &currency.EUR, nil
	case "USD":
		return &currency.USD, nil
	default:
		return nil, errors.New("Unknown currency " + profileResponse.Currency)
	}
}

func authenticateUrl(baseUrl *url.URL) error {
	apiKey, err := getApiKey()

	if err != nil {
		return err
	}

	baseUrl.User = apiKey
	return nil
}

func encodeMessageForSending(msg smsgw.Message) url.Values {
	vals := url.Values{}
	vals.Add("from", msg.From)
	vals.Add("to", msg.To)
	vals.Add("message", msg.Message)
	return vals
}

func getApiKey() (*url.Userinfo, error) {
	username := os.Getenv("TG_ELKS_USER")
	password := os.Getenv("TG_ELKS_PASS")

	if username == "" || password == "" {
		return nil, errors.New("Environment variables for 46elks username " +
			"and secret must be set.")
	}

	apiKey := url.UserPassword(username, password)

	return apiKey, nil
}
