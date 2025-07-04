package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
	// "time"
)

var (
	URL = "https://console.lightnode.com/api/ec/ecs/"
)

var Instances = map[string]InstanceParams{
	"ecs-x300003tygcz-bg-sofia-1": {
		InstanceUUID: "ecs-x300003tygcz",
		RegionCode:   "bg-sofia-1",
		ZoneCode:     "bg-sofia-1-a",
		ForceStop:    false,
	},
	// Add more instances here as needed
}

type InstanceParams struct {
	InstanceUUID string `json:"instanceUUID"`
	RegionCode   string `json:"regionCode"`
	ZoneCode     string `json:"zoneCode"`
	ForceStop    bool   `json:"forceStop"`
}

type instanceRequest struct {
	Timestamp       int64  `json:"timestamp"`
	BatchParamsJson string `json:"batchParamsJson"`
}

// TokenGenerator generates the Authorization token based on timestamp and token part
func TokenGenerator(timestamp int64, token string) string {
	// Format: <timestamp>:<token>:zhavzhCN9yBMWF%
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s:zhavzhCN9yBMWF%%", timestamp, token)))
}

func sendInstanceRequest(url, body string, token string, timestamp int64) error {
	log.Printf("Sending request to %s with body: %s", url, body)

	auth := TokenGenerator(timestamp, token)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Response from %s: %s", url, string(respBody))

	if resp.StatusCode != http.StatusOK {
		return errors.New("request failed: " + string(respBody))
	}

	return nil
}

func sendInstanceAction(url string, params InstanceParams) error {
	batchParams, err := json.Marshal([]InstanceParams{params})
	if err != nil {
		return err
	}
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	token := string(b)
	timestamp := time.Now().UnixMilli()
	requestBody := instanceRequest{
		Timestamp:       timestamp,
		BatchParamsJson: string(batchParams),
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	log.Printf("Action: %s, Params: %+v", url, params)
	return sendInstanceRequest(url, string(bodyBytes), token, timestamp)
}

func StartInstance(params InstanceParams) error {
	return sendInstanceAction(URL+"start.do", params)
}

func StopInstance(params InstanceParams) error {
	return sendInstanceAction(URL+"stop.do", params)
}
