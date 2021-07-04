package transport

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

const (
	firebaseToken           = "firebase_device_token"
	authReqId               = "AD1755F0-A083-4EC4-9ED7-59697110CC1C"
	firebaseServerKey       = "firebase_server_key"
	accessToken             = "d3373451-fef9-4c76-87f0-6537e364e366"
	tokenType               = "bearer"
	expiresIn               = 99999
	idToken                 = "f0c076bb-f903-4062-a0e0-03cb33a749e3"
	clientNotificationToken = "cb46d947-16fd-4089-8dd1-292d64a68e91"
	endpoint                = "https://auth.go-ciba.com"
	userId                  = "123"
)

func TestFirebaseCloudMessaging_Send(t *testing.T) {
	defer gock.Off()
	body := &fcmSendRequest{
		To:   fmt.Sprintf("/topics/ciba_consent/%s", userId),
		Data: make(map[string]interface{}),
	}
	body.Data["auth_req_id"] = authReqId
	jsonBody, _ := json.Marshal(body)
	gock.New("https://fcm.googleapis.com").
		Post("/fcm/send").
		MatchHeader("Authorization", "key="+firebaseServerKey).
		JSON(jsonBody).
		Reply(200)
	client := NewFirebaseCloudMessaging(firebaseServerKey)

	err := client.Send(map[string]interface{}{
		"to":               userId,
		"data.auth_req_id": authReqId,
	})

	assert.NoError(t, err)
}

func TestFirebaseCloudMessaging_Send_ShouldReturnError(t *testing.T) {
	defer gock.Off()
	body := &fcmSendRequest{
		To:   fmt.Sprintf("/topics/ciba_consent/%s", userId),
		Data: make(map[string]interface{}),
	}
	body.Data["auth_req_id"] = authReqId
	jsonBody, _ := json.Marshal(body)
	gock.New("https://fcm.googleapis.com").
		Post("/fcm/send").
		MatchHeader("Authorization", "key="+firebaseServerKey).
		JSON(jsonBody).
		Reply(400).
		JSON(map[string]string{"message": "validation error"})
	client := NewFirebaseCloudMessaging(firebaseServerKey)

	err := client.Send(map[string]interface{}{
		"to":               userId,
		"data.auth_req_id": authReqId,
	})

	assert.Error(t, err)
}

func TestClientAppNotification_Send_SuccessfulPush(t *testing.T) {
	defer gock.Off()
	requestBody := map[string]interface{}{
		"token_method":              domain.ModePush,
		"success":                   true,
		"auth_req_id":               authReqId,
		"access_token":              accessToken,
		"token_type":                tokenType,
		"expires_in":                expiresIn,
		"id_token":                  idToken,
		"client_notification_token": clientNotificationToken,
		"endpoint":                  endpoint,
	}
	jsonBody, _ := json.Marshal(map[string]interface{}{
		"auth_req_id":  authReqId,
		"access_token": accessToken,
		"token_type":   tokenType,
		"expires_in":   expiresIn,
		"id_token":     idToken,
	})
	gock.New(endpoint).
		Post("").
		MatchHeader("Authorization", "Bearer "+clientNotificationToken).
		JSON(jsonBody).
		Reply(200)
	client := NewClientAppNotificationClient()

	err := client.Send(requestBody)

	assert.NoError(t, err)
}

func TestClientAppNotification_Send_PushOidcError(t *testing.T) {
	defer gock.Off()
	requestBody := map[string]interface{}{
		"oidc_error":                *util.ErrAccessDenied,
		"client_notification_token": clientNotificationToken,
		"endpoint":                  endpoint,
		"success":                   false,
		"token_method":              domain.ModePush,
	}
	jsonBody, _ := json.Marshal(util.ErrAccessDenied)
	gock.New(endpoint).
		Post("").
		MatchHeader("Authorization", "Bearer "+clientNotificationToken).
		JSON(jsonBody).
		Reply(200)
	client := NewClientAppNotificationClient()

	err := client.Send(requestBody)

	assert.NoError(t, err)
}

func TestClientAppNotification_Send_SuccessfulPing(t *testing.T) {
	defer gock.Off()
	requestBody := map[string]interface{}{
		"client_notification_token": clientNotificationToken,
		"endpoint":                  endpoint,
		"success":                   true,
		"token_method":              domain.ModePing,
		"auth_req_id":               authReqId,
	}
	jsonBody, _ := json.Marshal(map[string]interface{}{
		"auth_req_id": authReqId,
	})
	gock.New(endpoint).
		Post("").
		MatchHeader("Authorization", "Bearer "+clientNotificationToken).
		JSON(jsonBody).
		Reply(200)
	client := NewClientAppNotificationClient()

	err := client.Send(requestBody)

	assert.NoError(t, err)
}
