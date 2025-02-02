package handlers

import (
	"fiber-app/pkg/database"
	"fiber-app/pkg/models"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	Content string `json:"content" example:"Merhaba, siparişiniz hazırlanıyor."` // Mesaj içeriği
	Phone   string `json:"phone" example:"+905551234567"`                        // Telefon numarası
}

// SuccessResponse represents the success response structure
type SuccessResponse struct {
	Status string      `json:"status" example:"success"` // İşlem durumu
	Data   interface{} `json:"data"`                     // Yanıt verisi
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Status  string `json:"status" example:"failed"`                    // İşlem durumu
	Message string `json:"message" example:"Content alanı zorunludur"` // Hata mesajı
	Code    string `json:"code,omitempty" example:"CONTENT_REQUIRED"`  // Hata kodu
}

// MessageResponse represents a single message response
type MessageResponse struct {
	Status string         `json:"status" example:"success"` // İşlem durumu
	Data   models.Message `json:"data"`                     // Mesaj verisi
}

// MessagesResponse represents multiple messages response
type MessagesResponse struct {
	Status string           `json:"status" example:"success"` // İşlem durumu
	Data   []models.Message `json:"data"`                     // Mesaj listesi
}

// Telefon numarası formatı kontrolü için regex
var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

// @Summary Yeni mesaj oluştur
// @Description Yeni bir mesaj oluşturur ve veritabanına kaydeder
// @Tags messages
// @Accept json
// @Produce json
// @Param message body CreateMessageRequest true "Mesaj bilgileri"
// @Success 201 {object} MessageResponse "Başarılı yanıt"
// @Failure 400 {object} ErrorResponse "Geçersiz istek"
// @Router /messages [post]
func CreateMessage(c *fiber.Ctx) error {
	var request CreateMessageRequest

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "JSON formatı geçersiz",
			Code:    "INVALID_JSON",
		})
	}

	// Debug log
	log.Printf("Received request: %+v", request)

	// Validate required fields
	if request.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Content alanı zorunludur",
			Code:    "CONTENT_REQUIRED",
		})
	}

	// Content karakter limiti kontrolü
	if len(request.Content) > 120 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Content alanı maksimum 120 karakter olabilir",
			Code:    "CONTENT_TOO_LONG",
		})
	}

	if request.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Phone alanı zorunludur",
			Code:    "PHONE_REQUIRED",
		})
	}

	// Telefon numarası format kontrolü
	if !phoneRegex.MatchString(request.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Geçersiz telefon numarası formatı. Örnek: +905551234567 veya 5551234567",
			Code:    "INVALID_PHONE_FORMAT",
		})
	}

	// Telefon numarası uzunluk kontrolü
	if len(request.Phone) > 15 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Telefon numarası 15 karakterden uzun olamaz",
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
			Message: "Geçersiz veri formatı. Lütfen girdiğiniz bilgileri kontrol edin",
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

// @Summary Tüm mesajları getir
// @Description Veritabanındaki tüm mesajları getirir
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} MessagesResponse "Başarılı yanıt"
// @Failure 500 {object} ErrorResponse "Sunucu hatası"
// @Router /messages [get]
func GetMessages(c *fiber.Ctx) error {
	var messages []models.Message

	result := database.DB.Find(&messages)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Status:  "failed",
			Message: "Mesajlar getirilemedi",
			Code:    "DATABASE_ERROR",
		})
	}

	return c.JSON(MessagesResponse{
		Status: "success",
		Data:   messages,
	})
}
