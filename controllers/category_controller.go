package controllers

import (
	"go-crud/config"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET /api/categories
func GetAllCategories(c echo.Context) error {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET data", []string{err.Error()}))
	}
	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", categories))
}

// GET /api/categories/:id
func GetCategoryByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category ID", []string{err.Error()}))
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET data", category))
}

// POST /api/categories (Admin only)
func CreateCategory(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}
	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Only admin can create category"}))
	}

	var req models.Category
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{err.Error()}))
	}

	if err := config.DB.Create(&req).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create category", []string{err.Error()}))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Category created successfully", req))
}

// PUT /api/categories/:id (Admin only)
func UpdateCategory(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}
	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Only admin can update category"}))
	}

	id, _ := strconv.Atoi(c.Param("id"))
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found", []string{err.Error()}))
	}

	var req models.Category
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input", []string{err.Error()}))
	}

	category.NamaCategory = req.NamaCategory
	if err := config.DB.Save(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update category", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Category updated successfully", category))
}

// DELETE /api/categories/:id (Admin only)
func DeleteCategory(c echo.Context) error {
	authUser, err := getAuthUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", []string{err.Error()}))
	}
	if !authUser.IsAdmin {
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Forbidden", []string{"Only admin can delete category"}))
	}

	id, _ := strconv.Atoi(c.Param("id"))
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete category", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Category deleted successfully", nil))
}
