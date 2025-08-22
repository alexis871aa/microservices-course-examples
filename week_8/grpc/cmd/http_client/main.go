package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

const (
	baseURL = "http://localhost:8080/api/v1/ufo"
	// Base URL для REST API через Envoy proxy

	httpTimeout = 30 * time.Second
	// Таймаут для HTTP запросов
)

// HTTPClient структура для работы с REST API
type HTTPClient struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient создает новый HTTP клиент
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: httpTimeout,
		},
		baseURL: baseURL,
	}
}

// SightingInfo структура для отправки данных о наблюдении (соответствует proto)
type SightingInfo struct {
	ObservedAt      string `json:"observed_at"`                // RFC3339 формат времени
	Location        string `json:"location"`                   // Место наблюдения
	Description     string `json:"description"`                // Описание объекта
	Color           string `json:"color,omitempty"`            // Цвет объекта (опционально)
	Sound           *bool  `json:"sound,omitempty"`            // Признак наличия звука (опционально)
	DurationSeconds *int32 `json:"duration_seconds,omitempty"` // Продолжительность в секундах (опционально)
}

// SightingUpdateInfo структура для обновления данных о наблюдении
type SightingUpdateInfo struct {
	ObservedAt      string `json:"observed_at,omitempty"`      // Время наблюдения
	Location        string `json:"location,omitempty"`         // Место наблюдения
	Description     string `json:"description,omitempty"`      // Описание объекта
	Color           string `json:"color,omitempty"`            // Цвет объекта
	Sound           *bool  `json:"sound,omitempty"`            // Признак наличия звука
	DurationSeconds *int32 `json:"duration_seconds,omitempty"` // Продолжительность в секундах
}

// Sighting полная структура наблюдения НЛО
type Sighting struct {
	UUID      string        `json:"uuid"`       // Уникальный идентификатор
	Info      *SightingInfo `json:"info"`       // Информация о наблюдении
	CreatedAt string        `json:"created_at"` // Время создания записи
	UpdatedAt string        `json:"updated_at"` // Время последнего обновления
	DeletedAt string        `json:"deleted_at"` // Время удаления (для мягкого удаления)
}

// CreateResponse ответ на создание наблюдения
type CreateResponse struct {
	UUID string `json:"uuid"` // UUID созданного наблюдения
}

// GetResponse ответ на получение наблюдения
type GetResponse struct {
	Sighting *Sighting `json:"sighting"` // Данные наблюдения
}

// createSighting создает новое наблюдение НЛО через HTTP API
func (c *HTTPClient) createSighting(ctx context.Context) (string, error) {
	// Генерируем случайные данные с помощью gofakeit
	observedAt := gofakeit.DateRange(
		time.Now().AddDate(-3, 0, 0), // за последние 3 года
		time.Now(),
	).Format(time.RFC3339)

	location := gofakeit.City() + ", " + gofakeit.StreetName()
	description := gofakeit.Sentence(gofakeit.Number(5, 15))

	// Создаем базовую информацию о наблюдении
	info := &SightingInfo{
		ObservedAt:  observedAt,
		Location:    location,
		Description: description,
	}

	// Иногда добавляем дополнительные поля (с вероятностью 70%)
	if gofakeit.Bool() {
		info.Color = gofakeit.Color()
	}

	if gofakeit.Bool() {
		sound := gofakeit.Bool()
		info.Sound = &sound
	}

	if gofakeit.Bool() {
		duration := gofakeit.Int32()
		info.DurationSeconds = &duration
	}

	// Оборачиваем в CreateRequest согласно proto
	createRequest := map[string]interface{}{
		"info": info,
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Отправляем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %w", err)
		}

		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Читаем ответ
	var createResp CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return createResp.UUID, nil
}

// getSighting получает информацию о наблюдении по UUID через HTTP API
func (c *HTTPClient) getSighting(ctx context.Context, uuid string) (*Sighting, error) {
	// Создаем URL с UUID
	url := fmt.Sprintf("%s/%s", c.baseURL, uuid)

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "application/json")

	// Отправляем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Читаем ответ
	var getResp GetResponse
	if err := json.NewDecoder(resp.Body).Decode(&getResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return getResp.Sighting, nil
}

// updateSighting обновляет наблюдение НЛО через HTTP API
func (c *HTTPClient) updateSighting(ctx context.Context, uuid string) error {
	// Генерируем рандомные данные для обновления
	updateInfo := &SightingUpdateInfo{}

	// Обновляем часть полей случайным образом
	if gofakeit.Bool() {
		observedAt := gofakeit.DateRange(
			time.Now().AddDate(-3, 0, 0),
			time.Now(),
		).Format(time.RFC3339)
		updateInfo.ObservedAt = observedAt
	}

	if gofakeit.Bool() {
		location := gofakeit.City() + ", " + gofakeit.StreetName()
		updateInfo.Location = location
	}

	if gofakeit.Bool() {
		description := gofakeit.Sentence(gofakeit.Number(5, 15))
		updateInfo.Description = description
	}

	if gofakeit.Bool() {
		updateInfo.Color = gofakeit.Color()
	}

	if gofakeit.Bool() {
		sound := gofakeit.Bool()
		updateInfo.Sound = &sound
	}

	if gofakeit.Bool() {
		duration := gofakeit.Int32()
		updateInfo.DurationSeconds = &duration
	}

	// Оборачиваем в UpdateRequest согласно proto
	updateRequest := map[string]interface{}{
		"update_info": updateInfo,
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем URL с UUID
	url := fmt.Sprintf("%s/%s", c.baseURL, uuid)

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Отправляем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// deleteSighting удаляет наблюдение НЛО через HTTP API
func (c *HTTPClient) deleteSighting(ctx context.Context, uuid string) error {
	// Создаем URL с UUID
	url := fmt.Sprintf("%s/%s", c.baseURL, uuid)

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Accept", "application/json")

	// Отправляем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Создаем HTTP клиент
	client := NewHTTPClient()

	log.Println("=== Тестирование UFO Sightings REST API через Envoy ===")
	log.Println()

	// 1. Создаем несколько наблюдений
	log.Println("🛸 Создание наблюдения НЛО")
	log.Println("===========================")
	uuid, err := client.createSighting(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при создании наблюдения: %v\n", err)
		return
	}

	// Выводим информацию о созданном наблюдении
	log.Printf("✅ Создано наблюдение НЛО: UUID=%s\n", uuid)

	// 2. Получаем информацию о наблюдении
	log.Println()
	log.Println("🔍 Получение информации о наблюдении")
	log.Println("==================================")
	sighting, err := client.getSighting(ctx, uuid)
	if err != nil {
		log.Printf("❌ Ошибка при получении наблюдения: %v\n", err)
		return
	}

	// Выводим информацию о полученном наблюдении
	log.Printf("✅ Получено наблюдение НЛО: UUID=%s", uuid)
	sightingJSON, err := json.MarshalIndent(sighting, "", "  ")
	if err != nil {
		log.Printf("❌ Ошибка при маршалинге наблюдения: %v", err)
		return
	}
	log.Printf("📄 Данные:\n%s\n", string(sightingJSON))

	// 3. Обновляем наблюдение
	log.Println("✏️ Обновление наблюдения")
	log.Println("=======================")

	err = client.updateSighting(ctx, uuid)
	if err != nil {
		log.Printf("❌ Ошибка при обновлении наблюдения: %v", err)
		return
	}

	log.Printf("✅ Наблюдение НЛО обновлено: UUID=%s\n", uuid)

	// 4. Проверяем обновленное наблюдение
	log.Println()
	log.Println("🔍 Проверка обновленного наблюдения")
	log.Println("=================================")
	updatedSighting, err := client.getSighting(ctx, uuid)
	if err != nil {
		log.Printf("❌ Ошибка при получении обновленного наблюдения: %v", err)
		return
	}

	// Выводим информацию об обновленном наблюдении
	log.Printf("✅ Получено обновленное наблюдение НЛО: UUID=%s", uuid)
	updatedSightingJSON, err := json.MarshalIndent(updatedSighting, "", "  ")
	if err != nil {
		log.Printf("❌ Ошибка при маршалинге обновленного наблюдения: %v", err)
		return
	}
	log.Printf("📄 Обновленные данные:\n%s\n", string(updatedSightingJSON))

	// 5. Удаляем наблюдение
	log.Println("🗑️ Удаление наблюдения")
	log.Println("======================")
	err = client.deleteSighting(ctx, uuid)
	if err != nil {
		log.Printf("❌ Ошибка при удалении наблюдения: %v", err)
		return
	}

	log.Printf("✅ Наблюдение НЛО удалено: UUID=%s\n", uuid)

	// 6. Проверяем удаленное наблюдение
	log.Println()
	log.Println("🔍 Проверка удаленного наблюдения")
	log.Println("=================================")
	deletedSighting, err := client.getSighting(ctx, uuid)
	if err != nil {
		log.Printf("❌ Ошибка при получении удаленного наблюдения: %v", err)
		return
	}

	// Выводим информацию об удаленном наблюдении
	log.Printf("✅ Получено удаленное наблюдение НЛО: UUID=%s", uuid)
	deletedSightingJSON, err := json.MarshalIndent(deletedSighting, "", "  ")
	if err != nil {
		log.Printf("❌ Ошибка при маршалинге удаленного наблюдения: %v", err)
		return
	}
	log.Printf("📄 Данные удаленного наблюдения:\n%s\n", string(deletedSightingJSON))

	log.Println("🎉 Тестирование REST API завершено!")
	log.Println()
	log.Println("📊 Статистика:")
	log.Printf("   • Envoy Proxy URL: %s", baseURL)
	log.Printf("   • Envoy Admin UI:  http://localhost:8081")
	log.Printf("   • Архитектура:     HTTP Client → Envoy → gRPC Server")
}
