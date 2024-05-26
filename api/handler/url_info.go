package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sykell-back/service"
	"sykell-back/utils"
)

// fetchURLInfoHandler handles requests to fetch information about a URL
func FetchURLInfoHandler(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into a struct containing the URL
	var request struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		utils.HandleError(w, "Error while processing the URL", err, http.StatusBadRequest)
		return
	}

	// Check URL validity
	if !utils.IsValidURL(request.URL) {
		utils.HandleError(w, "Error", errors.New("URL is not valid"), http.StatusBadRequest)
		return
	}

	// Fetch information about the URL
	info, err := service.UrlProcess(request.URL)
	if err != nil {
		utils.HandleError(w, "Error while processing the URL", err, http.StatusInternalServerError)
		return
	}

	// Send result
	utils.SendJSONResponse(w, info, http.StatusOK)
}
