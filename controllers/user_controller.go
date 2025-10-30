package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /users (hanya Admin)
func GetAllUsers(c echo.Context) error {
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User tidak ditemukan dalam context"})
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Akses ditolak. Hanya admin yang dapat melihat semua user.",
		})
	}

	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// GET /users/:id (user lihat profil sendiri)
func GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User tidak ditemukan dalam context"})
	}

	// konversi ID dari string ke uint
	var user models.User
	if err := config.DB.First(&user, "id = ?", idParam).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User tidak ditemukan"})
	}

	if !authUser.IsAdmin && authUser.ID != user.ID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Tidak boleh melihat data user lain"})
	}

	return c.JSON(http.StatusOK, user)
}

// PUT /users/:id (update data user)
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User tidak ditemukan dalam context"})
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User tidak ditemukan"})
	}

	// hanya admin atau pemilik akun
	if !authUser.IsAdmin && authUser.ID != user.ID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Tidak bisa mengedit user lain"})
	}

	type UpdateUserRequest struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Input tidak valid"})
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DELETE /users/:id (hanya admin)
func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	authUser, ok := c.Get("authUser").(models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User tidak ditemukan dalam context"})
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Hanya admin yang bisa menghapus user"})
	}

	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User berhasil dihapus"})
}
