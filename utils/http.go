package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type StandardResponse struct {
	Datetime string      `json:"datetime"`
	Data     interface{} `json:"data"`
}

// SendJSONResponse prend en charge l'envoi des réponses JSON.
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Construire la réponse standard avec les données et la date/heure actuelles
	response := StandardResponse{
		Datetime: time.Now().Format(time.RFC3339), // Format ISO 8601, ex. : 2006-01-02T15:04:05Z07:00
		Data:     data,
	}

	// Encoder et envoyer la réponse standardisée
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log l'erreur d'encodage et envoie un code d'erreur 500
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// HandleError vérifie s'il y a une erreur et gère la réponse appropriée.
func HandleError(w http.ResponseWriter, errorTitle string, err error, errorCode int) {
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, errorTitle+": "+err.Error(), errorCode)
	}
}

// check URL
func IsValidURL(url string) bool {
	// Define a regular expression pattern for a valid URL
	pattern := `^((http|https)://)?[^\s]{2,256}\.[a-z]{2,6}(/?[\w/#.:?+&=%@!~\$'\(\)\*\,\;\[\]-]*)?$`
	return url != "" && regexp.MustCompile(pattern).MatchString(url)
}

// fetchURL fetches the content of a given URL
func FetchURLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// CheckUrl tests the URL and returns the status code and any error encountered
func CheckUrl(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
