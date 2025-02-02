package cron

import (
	"bytes"
	"encoding/json"
	"fiber-app/pkg/cache"
	"fiber-app/pkg/database"
	"fiber-app/pkg/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronJob   *cron.Cron
	cronMutex sync.Mutex
	isRunning bool
	entryID   cron.EntryID
)

type WebhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func init() {
	cronJob = cron.New(cron.WithSeconds())
	isRunning = false
}

func logCronOperation(operation string, messageIDs []uint, count int, status bool, description string) {
	messageIDStrings := make([]string, len(messageIDs))
	for i, id := range messageIDs {
		messageIDStrings[i] = fmt.Sprint(id)
	}

	cronLog := models.CronLog{
		Operation:     operation,
		MessageIDs:    strings.Join(messageIDStrings, ","),
		MessagesCount: count,
		Status:        status,
		Description:   description,
	}

	if err := database.DB.Create(&cronLog).Error; err != nil {
		log.Printf("Error creating cron log: %v", err)
	}
}

func updateInactiveMessages() {
	var messages []models.Message

	// Her çalışmada sadece 2 mesaj işle
	batchSize := 2

	// Sadece 2 tane işlenmemiş mesajı al
	result := database.DB.Where("status = ?", false).Order("created_at asc").Limit(batchSize).Find(&messages)
	if result.Error != nil {
		log.Printf("Error fetching inactive messages: %v", result.Error)
		return
	}

	if len(messages) == 0 {
		log.Println("No inactive messages found")
		return
	}

	log.Printf("Processing %d messages in this cycle", len(messages))

	for _, message := range messages {
		log.Printf("Processing message ID: %d", message.ID)
		webhookURL := os.Getenv("WEBHOOK_URL")
		if webhookURL == "" {
			webhookURL = "https://webhook.site/03c75f60-8d13-47f9-b11b-4181faad6ce0"
		}

		requestBody := WebhookRequest{
			To:      message.Phone,
			Content: message.Content,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalf("Error marshaling request: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		req.Header.Add("Content-Type", "application/json")
		authKey := os.Getenv("WEBHOOK_AUTH_KEY")
		if authKey == "" {
			authKey = "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo"
		}
		req.Header.Add("x-ins-auth-key", authKey)

		// Debug için request body'yi oku ve tekrar oluştur
		bodyBytes, _ := io.ReadAll(req.Body)
		fmt.Println("Request Body:", string(bodyBytes))
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Tekrar oluştur

		fmt.Println("*************************************************")
		fmt.Println(req)
		fmt.Println("*************************************************")

		// IP'den banlandığı için servis isteği iptal edildi
		/*
			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}
			defer resp.Body.Close()
		*/

		// Simüle edilmiş başarılı yanıt
		simulatedResponse := WebhookResponse{
			Message:   "Message sent successfully",
			MessageID: fmt.Sprintf("SIMULATED_MSG_%d_%d", message.ID, time.Now().Unix()),
		}
		simulatedResponseBytes, _ := json.Marshal(simulatedResponse)

		// HTTP Yanıtı Debug
		fmt.Println("Response Status: 200 OK (Simulated)")
		fmt.Println("Response Body:", string(simulatedResponseBytes))

		log.Printf("Response for message %d - Status: 200 OK, Body: %s", message.ID, string(simulatedResponseBytes))

		// Parse the simulated response
		var response WebhookResponse
		if err := json.Unmarshal(simulatedResponseBytes, &response); err != nil {
			log.Printf("Error decoding response for message %d: %v", message.ID, err)
			logCronOperation("WEBHOOK_RESPONSE", []uint{message.ID}, 1, false, fmt.Sprintf("Response decode failed: %v", err))
			continue
		}

		log.Printf("Response for message %d: %+v", message.ID, response)

		// Always update as successful
		message.Status = true
		message.MessageID = response.MessageID
		if err := database.DB.Save(&message).Error; err != nil {
			log.Printf("Error updating message %d: %v", message.ID, err)
			logCronOperation("DATABASE_UPDATE", []uint{message.ID}, 1, false, fmt.Sprintf("DB update failed: %v", err))
			continue
		}

		// Redis cache'e kaydet
		cacheData := cache.MessageCache{
			ID:        message.ID,
			MessageID: response.MessageID,
			Status:    true,
			Content:   message.Content,
			Phone:     message.Phone,
		}
		if err := cache.SetMessageCache(message.ID, cacheData); err != nil {
			log.Printf("Warning: Failed to cache message %d: %v", message.ID, err)
		}

		log.Printf("Successfully updated message %d with message_id %s", message.ID, response.MessageID)
		logCronOperation("MESSAGE_PROCESSED", []uint{message.ID}, 1, true, fmt.Sprintf("Message processed successfully with ID: %s", response.MessageID))
	}
}

func StartCron() error {
	cronMutex.Lock()
	defer cronMutex.Unlock()

	if isRunning {
		return nil
	}

	schedule := os.Getenv("CRON_SCHEDULE")
	if schedule == "" {
		schedule = "*/30 * * * * *" // fallback schedule
	}

	var err error
	entryID, err = cronJob.AddFunc(schedule, updateInactiveMessages)
	if err != nil {
		logCronOperation("START", nil, 0, false, fmt.Sprintf("Failed to start cron: %v", err))
		return err
	}

	cronJob.Start()
	isRunning = true
	logCronOperation("START", nil, 0, true, "Cron job started successfully")
	log.Println("Cron job started")
	return nil
}

func StopCron() {
	cronMutex.Lock()
	defer cronMutex.Unlock()

	if !isRunning {
		return
	}

	cronJob.Remove(entryID)
	isRunning = false
	logCronOperation("STOP", nil, 0, true, "Cron job stopped successfully")
	log.Println("Cron job stopped")
}

func IsCronRunning() bool {
	cronMutex.Lock()
	defer cronMutex.Unlock()
	return isRunning
}

func GetCronLogs(limit int) ([]models.CronLog, error) {
	var logs []models.CronLog
	result := database.DB.Order("created_at desc").Limit(limit).Find(&logs)
	return logs, result.Error
}
