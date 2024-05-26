package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchURLInfoHandler(t *testing.T) {

	tests := []struct {
		name               string
		requestBody        map[string]string
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Valid URL",
			requestBody: map[string]string{
				"url": "https://www.google.com/",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid URL",
			requestBody: map[string]string{
				"url": "invalid-url",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Empty URL",
			requestBody: map[string]string{
				"url": "",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/urlinfo", bytes.NewReader(body))
			w := httptest.NewRecorder()

			FetchURLInfoHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			// Vérification uniquement du statut HTTP
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
		})
	}

}
