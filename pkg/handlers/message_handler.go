package handlers

import (
	"fiber-app/pkg/cache"
	"fiber-app/pkg/database"
	"fiber-app/pkg/models"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

type CreateMessageRequest struct {
	Content string `json:"content" example:"Hello, your order is being prepared."`
	Phone   string `json:"phone" example:"+905551234567"`
}

type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"Content field required"`
	Code    string `json:"code,omitempty" example:"CONTENT_REQUIRED"`
}

type MessageResponse struct {
	Status string         `json:"status" example:"success"`
	Data   models.Message `json:"data"`
}

type MessagesResponse struct {
	Status string           `json:"status" example:"success"`
	Data   []models.Message `json:"data"`
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

// @Summary Create new message
// @Description Creates a new message and saves it to the database
// @Tags messages
// @Accept json
// @Produce json
// @Param message body CreateMessageRequest true "Message information"
// @Success 201 {object} MessageResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /messages [post]
func CreateMessage(c *fiber.Ctx) error {
	var request CreateMessageRequest

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Invalid JSON format",
			Code:    "INVALID_JSON",
		})
	}

	// Debug log
	log.Printf("Received request: %+v", request)

	// Validate required fields
	if request.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Content field is required",
			Code:    "CONTENT_REQUIRED",
		})
	}

	// Content character limit check
	if len(request.Content) > 120 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Content field cannot exceed 120 characters",
			Code:    "CONTENT_TOO_LONG",
		})
	}

	if request.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Phone field is required",
			Code:    "PHONE_REQUIRED",
		})
	}

	// Phone number format validation
	if !phoneRegex.MatchString(request.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Invalid phone number format. Example: +905551234567 or 5551234567",
			Code:    "INVALID_PHONE_FORMAT",
		})
	}

	// Phone number length check
	if len(request.Phone) > 15 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Phone number cannot exceed 15 characters",
			Code:    "PHONE_TOO_LONG",
		})
	}

	// Create message
	message := models.Message{
		Content: request.Content,
		Phone:   request.Phone,
		Status:  false,
	}

	// Debug log before save
	log.Printf("Attempting to save message: %+v", message)

	result := database.DB.Create(&message)
	if result.Error != nil {
		log.Printf("Error creating message: %v", result.Error)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Invalid data format. Please check your input",
			Code:    "VALIDATION_ERROR",
		})
	}

	// Debug log after save
	log.Printf("Successfully created message: %+v", message)

	return c.Status(fiber.StatusCreated).JSON(MessageResponse{
		Status: "success",
		Data:   message,
	})
}

// @Summary Get all sent messages
// @Description Retrieves messages from database where status is true (sent)
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} MessagesResponse "Successful response"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /messages [get]
func GetMessages(c *fiber.Ctx) error {
	var messages []models.Message

	// Try from cache first
	cachedMessages, err := cache.GetMessageCacheWithTimeout(0) // 0 is special key for all messages
	if err != nil {
		log.Printf("Cache error: %v", err)
	} else if cachedMessages != nil {
		// Convert cache data to models.Message
		message := models.Message{
			ID:      cachedMessages.ID,
			Content: cachedMessages.Content,
			Phone:   cachedMessages.Phone,
			Status:  cachedMessages.Status,
		}
		return c.JSON(MessagesResponse{
			Status: "success",
			Data:   []models.Message{message},
		})
	}

	// If not in cache or error occurred, get from database
	result := database.DB.Where("status = ?", true).Find(&messages)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Failed to retrieve messages",
			Code:    "DATABASE_ERROR",
		})
	}

	// Save successful result to cache
	if len(messages) > 0 {
		go func() {
			// Save first message to cache
			cacheData := cache.MessageCache{
				ID:      messages[0].ID,
				Content: messages[0].Content,
				Phone:   messages[0].Phone,
				Status:  messages[0].Status,
			}
			if err := cache.SetMessageCache(0, cacheData); err != nil {
				log.Printf("Cache set error: %v", err)
			}
		}()
	}

	return c.JSON(MessagesResponse{
		Status: "success",
		Data:   messages,
	})
}
