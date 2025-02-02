package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nlopes/slack/slackevents"

	"github.com/kyma-incubator/hack-showcase/slack-connector/internal/handlers/mocks"
	slack "github.com/kyma-incubator/hack-showcase/slack-connector/internal/slack/mocks"

	"github.com/kyma-incubator/hack-showcase/slack-connector/internal/apperrors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type toJSON struct {
	TestJSON string `json:TestJSON"`
}

//createRequest creates an HTTP request for test purposes
func createRequest(t *testing.T) *http.Request {

	payload := toJSON{TestJSON: "test"}
	toSend, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(toSend))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestWebhookHandler(t *testing.T) {
	t.Run("Should respond with 403 status code when given a bad secret", func(t *testing.T) {
		// given

		payload := toJSON{TestJSON: "test"}
		toSend, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(toSend))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		mockValidator := &slack.Validator{}
		mockSender := &mocks.Sender{}

		mockValidator.On("GetToken").Return("test")
		mockValidator.On("ValidatePayload", req, []byte("test")).Return(nil, apperrors.AuthenticationFailed("fail"))

		// when
		wh := NewWebHookHandler(mockValidator, mockSender)

		handler := http.HandlerFunc(wh.HandleWebhook)
		handler.ServeHTTP(rr, req)

		// then
		mockValidator.AssertExpectations(t)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)

	})

	t.Run("Should respond with 400 status code when given wrong payload ", func(t *testing.T) {

		// given
		req := createRequest(t)
		rr := httptest.NewRecorder()

		mockValidator := &slack.Validator{}
		mockSender := &mocks.Sender{}
		mockPayload, err := json.Marshal(toJSON{TestJSON: "test"})
		require.NoError(t, err)

		mockValidator.On("GetToken").Return("test")
		mockValidator.On("ValidatePayload", req, []byte("test")).Return(mockPayload, nil)
		mockValidator.On("ParseWebHook", mockPayload).Return(nil, apperrors.WrongInput("fail"))

		wh := NewWebHookHandler(mockValidator, mockSender)

		// when
		handler := http.HandlerFunc(wh.HandleWebhook)
		handler.ServeHTTP(rr, req)

		// then
		mockValidator.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Should respond with 200 status code, when given a payload with a known event", func(t *testing.T) {

		// given
		req := createRequest(t)
		rr := httptest.NewRecorder()

		mockValidator := &slack.Validator{}
		mockSender := &mocks.Sender{}
		mockPayload, err := json.Marshal(toJSON{TestJSON: "test"})
		require.NoError(t, err)
		rawPayload := json.RawMessage(mockPayload)
		mockSender.On("SendToKyma", "", "v1", "", os.Getenv("SLACK_CONNECTOR_NAME")+"-app", rawPayload).Return(nil)

		mockValidator.On("GetToken").Return("test")
		mockValidator.On("ValidatePayload", req, []byte("test")).Return(mockPayload, nil)
		event := slackevents.EventsAPIEvent{}
		mockValidator.On("ParseWebHook", mockPayload).Return(event, nil)

		wh := NewWebHookHandler(mockValidator, mockSender)

		// when
		handler := http.HandlerFunc(wh.HandleWebhook)
		handler.ServeHTTP(rr, req)

		// then
		mockValidator.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should respond with 400 status code, when given a payload with an unknown event", func(t *testing.T) {

		// given
		req := createRequest(t)
		rr := httptest.NewRecorder()

		mockValidator := &slack.Validator{}
		mockSender := &mocks.Sender{}

		mockPayload, err := json.Marshal(toJSON{TestJSON: "test"})
		require.NoError(t, err)
		mockValidator.On("GetToken").Return("test")
		mockValidator.On("ValidatePayload", req, []byte("test")).Return(mockPayload, nil)
		mockValidator.On("ParseWebHook", mockPayload).Return(nil, apperrors.NotFound("Unknown event"))

		wh := NewWebHookHandler(mockValidator, mockSender)

		// when
		handler := http.HandlerFunc(wh.HandleWebhook)
		handler.ServeHTTP(rr, req)

		// then
		mockValidator.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

}
