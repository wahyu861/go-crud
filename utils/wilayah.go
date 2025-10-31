package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ================================
// 🗺️ Struct Wilayah
// ================================
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

// ================================
// ⚙️ HTTP Client (dengan timeout)
// ================================
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// ================================
// 🔹 GetProvinceByID
// ================================
func GetProvinceByID(id string) (*Province, error) {
	if id == "" {
		return nil, errors.New("id provinsi tidak boleh kosong")
	}

	url := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/province/%s.json", id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data provinsi: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provinsi dengan id %s tidak ditemukan (status %d)", id, resp.StatusCode)
	}

	var province Province
	if err := json.NewDecoder(resp.Body).Decode(&province); err != nil {
		return nil, fmt.Errorf("gagal decode data provinsi: %v", err)
	}

	return &province, nil
}

// ================================
// 🔹 GetCityByID
// ================================
func GetCityByID(id string) (*City, error) {
	if id == "" {
		return nil, errors.New("id kota tidak boleh kosong")
	}

	url := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/regency/%s.json", id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data kota: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kota dengan id %s tidak ditemukan (status %d)", id, resp.StatusCode)
	}

	var city City
	if err := json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return nil, fmt.Errorf("gagal decode data kota: %v", err)
	}

	return &city, nil
}

// ================================
// 🔹 GetAllProvinces
// ================================
func GetAllProvinces() ([]Province, error) {
	url := "https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar provinsi: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal mengambil daftar provinsi (status %d)", resp.StatusCode)
	}

	var provinces []Province
	if err := json.NewDecoder(resp.Body).Decode(&provinces); err != nil {
		return nil, fmt.Errorf("gagal decode daftar provinsi: %v", err)
	}

	return provinces, nil
}

// ================================
// 🔹 GetCitiesByProvinceID
// ================================
func GetCitiesByProvinceID(provinceID string) ([]City, error) {
	if provinceID == "" {
		return nil, errors.New("id provinsi tidak boleh kosong")
	}

	url := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/%s.json", provinceID)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar kota: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gagal mengambil daftar kota (status %d)", resp.StatusCode)
	}

	var cities []City
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return nil, fmt.Errorf("gagal decode daftar kota: %v", err)
	}

	return cities, nil
}
