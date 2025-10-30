package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email sudah terdaftar"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal hash password"})
	}

	// Simpan user baru
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		IsAdmin:  req.IsAdmin,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Buat toko otomatis hanya jika user bukan admin
	if !user.IsAdmin {
		store := models.Store{
			UserID:      user.ID,
			Name:        "Toko " + user.Name,
			Description: "Toko milik " + user.Name,
		}
		config.DB.Create(&store)
	}

	user.Password = ""

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Registrasi berhasil",
		"user":    user,
	})
}

// ===================================================
// 🔐 PROSES LOGIN (POST)
// ===================================================
func Login(c echo.Context) error {
	var input models.User
	var user models.User

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Cari user berdasarkan email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Email tidak ditemukan"})
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Password salah"})
	}

	// Generate token JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // berlaku 3 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal membuat token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

// ===================================================
// 👤 PROFILE (cek token JWT)
// ===================================================
func Profile(c echo.Context) error {
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "User tidak ditemukan dalam context",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Token valid",
		"user":    authUser,
	})
}

