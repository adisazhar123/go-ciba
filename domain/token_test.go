package domain

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

func newIdTokenManager() TokenManager {
	return TokenManager{e: newEncryptionUtil()}
}

func newEncryptionUtil() util.EncryptionInterface {
	return &util.GoJoseEncryption{}
}

func TestIdTokenManager_CreateIdToken(t *testing.T) {
	mgr := newIdTokenManager()
	file, _ := os.Open("../test_data/key.pem")
	defer file.Close()

	privateKey, _ := ioutil.ReadAll(file)
	alg := "RS256"
	kId := "3d585c88-b2ac-4a07-824f-649cec260aa5"
	token := "2b7d56d2-7f1e-407f-896d-71bebfbdc0d4"
	now := 1596391533

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNkNTg1Yzg4LWIyYWMtNGEwNy04MjRmLTY0OWNlYzI2MGFhNSIsInR5cCI6Imp3dCJ9.eyJhdF9oYXNoIjoiWkRnek1HSm1OVEExTnpBeE9UZzJPREl4TWpVM1pXSXdaREE0Wm1RMllnPT0iLCJhdWQiOiJlNjczOGMxMS0zNjYzLTQzMDEtYmExMi0yMzEyZWI5YTEzNDQiLCJhdXRoX3JlcV9pZCI6ImYxNjlhZmVlLThhMmEtNDFkNS1iY2RhLTQ1MGY5OTk0MmYwZSIsImF1dGhfdGltZSI6MTU5NjM5MTUzMywiZXhwIjoxNTk2Mzk1MTMzLCJpYXQiOjE1OTYzOTE1MzMsImlzcyI6ImNpYmEtc2VydmVyLmV4YW1wbGUuY29tIiwibm9uY2UiOiI3Mjc1NDk5OS01ZmU1LTQ5NmItOTE4Mi05MDZmMTVhN2I4NmUiLCJzdWIiOiIzZTQzZjA4Yy05MzRhLTRmNmMtODBiNi02ZjIwNWFiZDk1NGIifQ.X0j8rDm5jJbP0WiYYZDdns4q8xPc6DWLNui7leiQZjUlqzdJWU-WG7la7lQhwbVwcahJbnTAOEIY2kgKuEsPyfrthM5OE9wdbpORVxrg4qvge6HHPRWRcVy50W_sE1_AbtQPk4u69O6z5n7dufkHRHR5iyJfM4dZjuvpL1pMJ1T_qMhMfusaXb7ZalVqGdPRwGior5sfYwJ2jN3Bk7AkOi22LL3BclC5CQLkyPKxAS6O3HfhwoKJ-vr1YYfB7zR41hv_ZmNzLQzApVCHG0yg7qJ62O7NTDVQwIbHWPD_yaGP1t5bSu-Yp8PdOKgc3WQRKuEa531GX6pjN-rCy5aPZGPoUCNETepPVNvOWxIm_Jb0S9dN_BLrjO8TisNQPxdQJB_I2gWkPIDDfQ5jnPQpCtemzx6kH9-ntOfmyxRzXwD0czsDS8XXGFAohUnE9cm0aTx8CZ3eTd1wPliX-f5WmSqLc6wnXsOg9IkIovMqQHE05aKoy9UnSOH5nM8qjbY2lks-iqfSzNWp_9OQQBl0q0esH5mfQolTvBb924W__U9H5bCcA8Rd9EJrbfAWY93jpVwvrFUMqp0eAdHXnYuyCuz7APdTEWORRGosd1KkXisJJnHMnga5gnF2hWxaJiXIEM4kiA-AnxiCYV6iOHhaAguATmfLCC_Foh4-0p79BsI"

	claims := make(map[string]interface{})
	claims["auth_req_id"] = "f169afee-8a2a-41d5-bcda-450f99942f0e"
	claims["aud"] = "e6738c11-3663-4301-ba12-2312eb9a1344"
	claims["auth_time"] = now
	claims["iat"] = now
	claims["exp"] = now + 3600
	claims["iss"] = "ciba-server.example.com"
	claims["sub"] = "3e43f08c-934a-4f6c-80b6-6f205abd954b"
	claims["nonce"] = "72754999-5fe5-496b-9182-906f15a7b86e"

	et := mgr.CreateIdToken(claims, string(privateKey), alg, kId, token)

	assert.NotNil(t, et)
	assert.Equal(t, idToken, et.value)
}