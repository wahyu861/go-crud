package controllers

import (
	"encoding/json"
	"net/http"

	"go-crud/utils"

	"github.com/labstack/echo/v4"
)

// Struct untuk mapping JSON dari API emsifa
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

// ===================================================
// 🔹 GET /provcity/listprovinces
// ===================================================
func GetListProvinces(c echo.Context) error {
	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET provinces", []string{err.Error()}))
	}
	defer resp.Body.Close()

	var provinces []Province
	if err := json.NewDecoder(resp.Body).Decode(&provinces); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to parse provinces", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET provinces", provinces))
}

// ===================================================
// 🔹 GET /provcity/listcities/:province_id
// ===================================================
func GetListCities(c echo.Context) error {
	provinceID := c.Param("province_id")

	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + provinceID + ".json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET cities", []string{err.Error()}))
	}
	defer resp.Body.Close()

	var cities []City
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to parse cities", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET cities", cities))
}

// ===================================================
// 🔹 GET /provcity/detailprovince/:id
// ===================================================
func GetDetailProvince(c echo.Context) error {
	id := c.Param("id")

	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/province/" + id + ".json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET province detail", []string{err.Error()}))
	}
	defer resp.Body.Close()

	var province Province
	if err := json.NewDecoder(resp.Body).Decode(&province); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to parse province detail", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET province detail", province))
}

// ===================================================
// 🔹 GET /provcity/detailcity/:id
// ===================================================
func GetDetailCity(c echo.Context) error {
	id := c.Param("id")

	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regency/" + id + ".json")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to GET city detail", []string{err.Error()}))
	}
	defer resp.Body.Close()

	var city City
	if err := json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to parse city detail", []string{err.Error()}))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Succeed to GET city detail", city))
}
