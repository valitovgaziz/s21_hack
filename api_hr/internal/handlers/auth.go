package handlers

import (
	"math/rand"
	"net/http"
	"time"
	"api_hr/internal/models"
	"api_hr/internal/utils"

	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var existingUser models.UserT
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.WriteError(w, http.StatusConflict, "User already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	user := models.UserT{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
	}
	if user.Role == "" {
		user.Role = "employee"
	}

	if err := h.DB.Create(&user).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var user models.UserT
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) LoginPhone(w http.ResponseWriter, r *http.Request) {
	var req models.LoginPhoneRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Phone == "" {
		utils.WriteError(w, http.StatusBadRequest, "Phone is required")
		return
	}

	code := generateOTP()
	expiresAt := time.Now().Add(5 * time.Minute)

	session := models.OTPSession{
		Phone:     req.Phone,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	if err := h.DB.Create(&session).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error creating OTP")
		return
	}

	// В MVP просто возвращаем код (в реальности отправляем через SMS)
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "OTP sent",
		"code":    code,
		"expires": expiresAt,
	})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyOTPRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var session models.OTPSession
	if err := h.DB.Where("phone = ? AND code = ? AND used = ? AND expires_at > ?",
		req.Phone, req.Code, false, time.Now()).
		First(&session).Error; err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid or expired OTP")
		return
	}

	h.DB.Model(&session).Update("used", true)

	var user models.UserT
	if err := h.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		user = models.UserT{
			Name:  "User " + req.Phone[len(req.Phone)-4:],
			Phone: req.Phone,
			Email: req.Phone + "@pulsehr.local",
			Role:  "employee",
		}
		hashedPassword, _ := utils.HashPassword(req.Phone + "pulsehr")
		user.Password = hashedPassword
		h.DB.Create(&user)
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var user models.UserT
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"valid": true,
		"user":  user,
	})
}

func generateOTP() string {
	code := rand.Intn(1000000)
	return formatCode(code)
}

func formatCode(code int) string {
	digits := make([]byte, 6)
	for i := 5; i >= 0; i-- {
		digits[i] = byte('0' + code%10)
		code /= 10
	}
	return string(digits)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
