package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ========================== RESPONSE STRUCT ==========================

type SafeUserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type StoreResponse struct {
	ID          uint64            `json:"id"`
	NamaToko    string            `json:"nama_toko"`
	UrlFoto     string            `json:"url_foto"`
	Description string            `json:"description,omitempty"`
	UserID      uint64            `json:"user_id"`
	User        *SafeUserResponse `json:"user,omitempty"`
}

// ===================== HELPER RESPONSE ======================
type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

func successResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  true,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func errorResponse(message string, errs []string) BaseResponse {
	return BaseResponse{
		Status:  false,
		Message: message,
		Errors:  errs,
		Data:    nil,
	}
}

// Ambil user dari context JWT
func getAuthUser(c echo.Context) (*models.User, error) {
	authUserRaw := c.Get("authUser")
	if authUserRaw == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "User tidak ditemukan dalam context (token tidak valid?)")
	}

	authUser, ok := authUserRaw.(models.User)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Gagal membaca data user dari context")
	}
	return &authUser, nil
}

// ========================== HANDLER ===============================

// GET /api/stores (Admin only)
func GetAllStores(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("Failed to GET data", []string{err.Error()}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, errorResponse("Failed to GET data", []string{"Hanya admin yang dapat melihat semua toko"}))
	}

	var stores []models.Store
	if err := config.DB.Preload("User").Find(&stores).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to GET data", []string{err.Error()}))
	}

	var result []StoreResponse
	for _, s := range stores {
		resp := StoreResponse{
			ID:       s.ID,
			NamaToko: s.Name,
			UrlFoto:  s.UrlFoto,
			UserID:   s.UserID,
		}
		if s.User.ID != 0 {
			resp.User = &SafeUserResponse{
				ID:    s.User.ID,
				Name:  s.User.Name,
				Email: s.User.Email,
			}
		}
		result = append(result, resp)
	}

	data := map[string]interface{}{
		"page":  1,
		"limit": 10,
		"data":  result,
	}

	return c.JSON(http.StatusOK, successResponse("Succeed to GET data", data))
}

// GET /api/stores/my
func GetMyStore(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("Failed to GET data", []string{err.Error()}))
	}

	var store models.Store
	if err := config.DB.Preload("User").Where("user_id = ?", authUser.ID).First(&store).Error; err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("Failed to GET data", []string{"Toko tidak ditemukan"}))
	}

	resp := StoreResponse{
		ID:       store.ID,
		NamaToko: store.Name,
		UrlFoto:  store.UrlFoto,
		UserID:   store.UserID,
	}

	return c.JSON(http.StatusOK, successResponse("Succeed to GET data", resp))
}

// GET /api/stores/:id
func GetStoreByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to GET data", []string{"ID toko tidak valid"}))
	}

	var store models.Store
	if err := config.DB.Preload("User").First(&store, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("Failed to GET data", []string{"Toko tidak ditemukan"}))
	}

	resp := StoreResponse{
		ID:       store.ID,
		NamaToko: store.Name,
		UrlFoto:  store.UrlFoto,
		UserID:   store.UserID,
	}

	return c.JSON(http.StatusOK, successResponse("Succeed to GET data", resp))
}

// PUT /api/stores/:id
func UpdateStore(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to UPDATE data", []string{"ID toko tidak valid"}))
	}

	var store models.Store
	if err := config.DB.First(&store, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("Failed to UPDATE data", []string{"Toko tidak ditemukan"}))
	}

	if store.UserID != authUser.ID && !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, errorResponse("Failed to UPDATE data", []string{"Tidak bisa mengubah toko milik orang lain"}))
	}

	var input struct {
		NamaToko    string `json:"nama_toko" form:"nama_toko"`
		UrlFoto     string `json:"url_foto" form:"url_foto"`
		Description string `json:"description" form:"description"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to UPDATE data", []string{"Input tidak valid"}))
	}

	updates := map[string]interface{}{}
	if input.NamaToko != "" {
		updates["name"] = input.NamaToko
	}
	if input.UrlFoto != "" {
		updates["url_foto"] = input.UrlFoto
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to UPDATE data", []string{"Tidak ada data yang diubah"}))
	}

	if err := config.DB.Model(&store).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, successResponse("Succeed to UPDATE data", "Update toko succeed"))
}

// DELETE /api/stores/:id (admin only - nonaktifkan toko)
func DeleteStore(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errorResponse("Unauthorized", []string{err.Error()}))
	}

	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, errorResponse("Forbidden", []string{"Hanya admin yang dapat menonaktifkan toko"}))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to DELETE data", []string{"ID toko tidak valid"}))
	}

	var store models.Store
	if err := config.DB.First(&store, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, errorResponse("Failed to DELETE data", []string{"Toko tidak ditemukan"}))
	}

	// Jika sudah tidak aktif, jangan ubah lagi
	if !store.IsActive {
		return c.JSON(http.StatusBadRequest, errorResponse("Failed to DELETE data", []string{"Toko sudah nonaktif"}))
	}

	store.IsActive = false
	if err := config.DB.Save(&store).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse("Failed to UPDATE data", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, successResponse("Succeed to DELETE data", "Toko berhasil dinonaktifkan"))
}
