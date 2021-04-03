package transport

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestFirebaseCloudMessaging_Send(t *testing.T) {
	defer gock.Off()
	var token, authReqId, serverKey = "firebase_device_token", "AD1755F0-A083-4EC4-9ED7-59697110CC1C", "firebase_server_key"
	body := &fcmSendRequest{
		To:   token,
		Data: make(map[string]interface{}),
	}
	body.Data["auth_req_id"] = authReqId
	jsonBody, _ := json.Marshal(body)

	gock.New("https://fcm.googleapis.com").
		Post("/fcm/send").
		MatchHeader("Authorization", "key=" + serverKey).
		JSON(jsonBody).
		Reply(200)


	client := NewFirebaseCloudMessaging(serverKey)

	err := client.Send(map[string]interface{}{
		"to":               token,
		"data.auth_req_id": authReqId,
	})

	assert.NoError(t, err)
}

func TestClientAppNotification_Send_ShouldReturnError(t *testing.T) {
	defer gock.Off()
	var token, authReqId, serverKey = "firebase_device_token", "AD1755F0-A083-4EC4-9ED7-59697110CC1C", "firebase_server_key"
	body := &fcmSendRequest{
		To:   token,
		Data: make(map[string]interface{}),
	}
	body.Data["auth_req_id"] = authReqId
	jsonBody, _ := json.Marshal(body)

	gock.New("https://fcm.googleapis.com").
		Post("/fcm/send").
		MatchHeader("Authorization", "key=" + serverKey).
		JSON(jsonBody).
		Reply(400).
		JSON(map[string]string{"message": "validation error"})


	client := NewFirebaseCloudMessaging(serverKey)

	err := client.Send(map[string]interface{}{
		"to":               token,
		"data.auth_req_id": authReqId,
	})

	assert.Error(t, err)
}