package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type NotificationInterface interface {
	Send(data map[string]interface{}) error
}

type FirebaseCloudMessaging struct {
	client *http.Client
}

func NewFirebaseCloudMessaging() *FirebaseCloudMessaging {
	return &FirebaseCloudMessaging{client: &http.Client{
		Timeout: 5 * time.Second,
	}}
}

type fcmSendRequest struct {
	To   string                 `json:"to"`
	Data map[string]interface{} `json:"data"`
}

func (f *FirebaseCloudMessaging) Send(data map[string]interface{}) error {
	body := &fcmSendRequest{
		To:   data["to"].(string),
		Data: make(map[string]interface{}),
	}

	for k, v := range data {
		if strings.Contains(k, "data.") {
			key := strings.Split(k, "data.")[1]
			body.Data[key] = v
		}
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonBody))

	res, err := f.client.Do(req)

	if err != nil {
		log.Printf("[go-ciba][firebasecloudmessaging] an error occured %s\n", err.Error())
		return err
	}

	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Printf("[go-ciba][firebasecloudmessaging] non OK status code received %d\n", res.StatusCode)
		log.Printf("[go-ciba][firebasecloudmessaging] received response body %s\n", string(resBody))
		return errors.New(fmt.Sprintf("failed to send notification message %s", string(resBody)))
	}
	return nil
}

type ClientAppNotification struct {
	client *http.Client
}

type PushTokenRequest struct {
	AuthReqId string `json:"auth_req_id"`
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	IdToken string `json:"id_token"`
}

func (c *ClientAppNotification) Send(data map[string]interface{}) error {
	body := &PushTokenRequest{
		AuthReqId:   data["auth_req_id"].(string),
		AccessToken: data["access_token"].(string),
		TokenType:   data["token_type"].(string),
		ExpiresIn:   data["expires_in"].(int),
		IdToken:     data["id_token"].(string),
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, data["client_notification_endpoint"].(string), bytes.NewBuffer(jsonBody))
	res, err := c.client.Do(req)

	if err != nil {
		log.Printf("[go-ciba][pushtoken] an error occured %s\n", err.Error())
		return err
	}

	defer res.Body.Close()
	_, _ = ioutil.ReadAll(res.Body)
	return nil
}
