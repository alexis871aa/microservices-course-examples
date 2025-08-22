package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"
)

// Простая конфигурация в константах
const (
	baseURL     = "http://localhost:8080"
	sessionUUID = "dc31f7c0-b736-406f-b328-76f270764959" // демо-UUID
)

type createRequest struct {
	Info sightingInfo `json:"info"`
}

type updateRequest struct {
	UUID       string          `json:"uuid"`
	UpdateInfo sightingUpdInfo `json:"update_info"`
}

type sightingInfo struct {
	ObservedAt      string  `json:"observed_at"`
	Location        string  `json:"location"`
	Description     string  `json:"description"`
	Color           *string `json:"color,omitempty"`
	Sound           *bool   `json:"sound,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}

type sightingUpdInfo struct {
	ObservedAt      *string `json:"observed_at,omitempty"`
	Location        *string `json:"location,omitempty"`
	Description     *string `json:"description,omitempty"`
	Color           *string `json:"color,omitempty"`
	Sound           *bool   `json:"sound,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}

type createResponse struct {
	UUID string `json:"uuid"`
}

func main() {
	client := &http.Client{Timeout: 15 * time.Second}

	// 1) CREATE
	observedAt := time.Now().UTC().Format(time.RFC3339)
	createReq := createRequest{
		Info: sightingInfo{
			ObservedAt:      observedAt,
			Location:        "Area 51, Nevada",
			Description:     "Bright circular object moving at high speed",
			Color:           lo.ToPtr("silver"),
			Sound:           lo.ToPtr(true),
			DurationSeconds: lo.ToPtr(120),
		},
	}

	fmt.Println("=> CREATE")
	createRespBody := doJSONRequest(client, http.MethodPost, baseURL+"/api/v1/ufo", createReq)

	var cr createResponse
	err := json.Unmarshal(createRespBody, &cr)
	if err != nil {
		log.Printf("unmarshal create response: %v\n", err.Error())
		return
	}
	if strings.TrimSpace(cr.UUID) == "" {
		log.Printf("create: empty UUID in response: %s\n", string(createRespBody))
		return
	}

	// 2) GET
	fmt.Println("=> GET (after create)")
	_ = doJSONRequest(client, http.MethodGet, baseURL+"/api/v1/ufo/"+cr.UUID, nil)

	// 3) UPDATE (partial)
	fmt.Println("=> UPDATE (description)")
	upd := updateRequest{
		UUID: cr.UUID,
		UpdateInfo: sightingUpdInfo{
			Description: lo.ToPtr("Updated description: moved to the east"),
		},
	}
	_ = doJSONRequest(client, http.MethodPut, baseURL+"/api/v1/ufo/"+cr.UUID, upd)

	// 4) GET (after update)
	fmt.Println("=> GET (after update)")
	_ = doJSONRequest(client, http.MethodGet, baseURL+"/api/v1/ufo/"+cr.UUID, nil)

	// 5) DELETE (soft delete)
	fmt.Println("=> DELETE (soft)")
	_ = doJSONRequest(client, http.MethodDelete, baseURL+"/api/v1/ufo/"+cr.UUID, nil)

	// 6) GET (after delete) — запись существует, но с filled deleted_at
	fmt.Println("=> GET (after soft delete)")
	_ = doJSONRequest(client, http.MethodGet, baseURL+"/api/v1/ufo/"+cr.UUID, nil)
}

func doJSONRequest(client *http.Client, method string, url string, body any) []byte {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			log.Printf("marshal request: %v\n", err.Error())
			return nil
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		log.Printf("new request: %v\n", err.Error())
		return nil
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Авторизация всегда включена (cookie X-Session-Uuid)
	req.Header.Add("Cookie", fmt.Sprintf("X-Session-Uuid=%s", urlEscape(sessionUUID)))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("http do: %v\n", err.Error())
		return nil
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("HTTP %d\n%s\n", resp.StatusCode, string(b))
	return b
}

func urlEscape(s string) string {
	// Avoid importing net/url for a single function; simple replacement suffices for demo
	r := strings.NewReplacer(" ", "%20", ";", "%3B", ",", "%2C")
	return r.Replace(s)
}
