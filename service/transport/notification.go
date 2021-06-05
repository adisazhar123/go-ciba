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

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/util"
)

type NotificationInterface interface {
	Send(data map[string]interface{}) error
}

type FirebaseCloudMessaging struct {
	client    *http.Client
	serverKey string
	baseUrl string
}

func NewFirebaseCloudMessaging(serverKey string) *FirebaseCloudMessaging {
	return &FirebaseCloudMessaging{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		serverKey: serverKey,
		baseUrl: "https://fcm.googleapis.com/fcm/send",
	}
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
	req, _ := http.NewRequest(http.MethodPost, f.baseUrl, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", f.serverKey))
	req.Header.Add("Content-Type", "application/json")

	fmt.Println(string(jsonBody))

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

func NewClientAppNotificationClient() *ClientAppNotification {
	return &ClientAppNotification{client: &http.Client{
		Timeout: 5 * time.Second,
	}}
}

type TokenCallbackRequest struct {
	AuthReqId               string `json:"auth_req_id,omitempty"`
	AccessToken             string `json:"access_token,omitempty"`
	TokenType               string `json:"token_type,omitempty"`
	ExpiresIn               int    `json:"expires_in,omitempty"`
	IdToken                 string `json:"id_token,omitempty"`
	endpoint                string
	clientNotificationToken string

	util.OidcError
}

func (c *ClientAppNotification) buildRequest(data map[string]interface{}) *TokenCallbackRequest {
	var body *TokenCallbackRequest = nil

	tokenMethod := data["token_method"].(string)

	if tokenMethod == domain.ModePush {
		success := data["success"].(bool)
		if success {
			body = &TokenCallbackRequest{
				AuthReqId:               data["auth_req_id"].(string),
				AccessToken:             data["access_token"].(string),
				TokenType:               data["token_type"].(string),
				ExpiresIn:               data["expires_in"].(int),
				IdToken:                 data["id_token"].(string),
				clientNotificationToken: data["client_notification_token"].(string),
				endpoint:                data["endpoint"].(string),
			}
		} else {
			body = &TokenCallbackRequest{
				OidcError:               data["oidc_error"].(util.OidcError),
				clientNotificationToken: data["client_notification_token"].(string),
				endpoint:                data["endpoint"].(string),
			}
		}
	} else if tokenMethod == domain.ModePing {
		body = &TokenCallbackRequest{
			AuthReqId:               data["auth_req_id"].(string),
			clientNotificationToken: data["client_notification_token"].(string),
			endpoint:                data["endpoint"].(string),
		}
	}

	return body
}

func (c *ClientAppNotification) Send(data map[string]interface{}) error {
	body := c.buildRequest(data)
	jsonBody, _ := json.Marshal(body)

	fmt.Println(string(jsonBody))

	req, _ := http.NewRequest(http.MethodPost, body.endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", body.clientNotificationToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)

	if err != nil {
		log.Printf("[go-ciba][pushtoken] an error occured %s\n", err.Error())
		return err
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		log.Printf("[go-ciba][client-app-notification] non OK status code received %d\n", res.StatusCode)
		log.Printf("[go-ciba][client-app-notification] received response body %s\n", string(resBody))
		return errors.New(fmt.Sprintf("failed to send authentication result %s", string(resBody)))
	}

	return nil
}
