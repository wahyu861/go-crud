package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ===================================================
// 🧾 REGISTER (POST)
// ===================================================
func Register(c echo.Context) error {
	type RegisterRequest struct {
		Nama      string  `json:"nama"`
		Email     string  `json:"email"`
		NoTelp    *string `json:"no_telp"`
		KataSandi string  `json:"kata_sandi"`
		IsAdmin   bool    `json:"is_admin"`
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", []string{err.Error()}))
	}

	// Cek apakah email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Email sudah terdaftar", nil))
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal hash password", []string{err.Error()}))
	}

	// Simpan user baru
	user := models.User{
		Nama:      req.Nama,
		Email:     req.Email,
		NoTelp:    req.NoTelp,
		KataSandi: string(hashedPassword),
		IsAdmin:   req.IsAdmin,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal menyimpan data pengguna", []string{err.Error()}))
	}

	// Buat toko otomatis hanya jika user bukan admin
	if !user.IsAdmin {
		toko := models.Toko{
			IDUser:   user.ID,
			NamaToko: "Toko " + user.Nama,
			UrlFoto:  nil,
		}
		config.DB.Create(&toko)
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Register berhasil", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}))
}

// ===================================================
// 🔐 LOGIN (POST)
// ===================================================
func Login(c echo.Context) error {
	var input models.User
	var user models.User

	// Bind input JSON
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", []string{err.Error()}))
	}

	// Cari user berdasarkan email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Email tidak ditemukan", []string{"email_not_found"}))
	}

	// Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(input.KataSandi)); err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Password salah", []string{"invalid_password"}))
	}

	// Buat token JWT (berlaku 3 hari)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Gagal membuat token", []string{err.Error()}))
	}

	// Siapkan data respons
	data := map[string]interface{}{
		"id":             user.ID,
		"nama":           user.Nama,
		"no_telp":        user.NoTelp,
		"tanggal_lahir":  user.TanggalLahir,
		"tentang":        user.Tentang,
		"pekerjaan":      user.Pekerjaan,
		"email":          user.Email,
		"id_provinsi":    user.IDProvinsi,
		"id_kota":        user.IDKota,
		"is_admin":       user.IsAdmin,
		"token":          tokenString,
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Login berhasil", data))
}

// ===================================================
// 👤 PROFILE (GET)
// ===================================================
func Profile(c echo.Context) error {
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User tidak ditemukan dalam context", []string{"unauthorized"}))
	}

	authUser.KataSandi = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Token valid", map[string]interface{}{
		"user": authUser,
	}))
}
