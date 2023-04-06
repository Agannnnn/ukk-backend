package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"lelang.id/backend/utils"
)

func SaveCategory(c echo.Context) error {
	var category utils.Kategori
	c.Bind(&category)

	db := utils.DB()
	if db.Create(&category).Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Gagal menyimpan kategori"})
	} else {
		return c.JSON(http.StatusOK, echo.Map{"message": "Kategori disimpan"})
	}
}

func ListCategories(c echo.Context) error {
	var categories []utils.Kategori

	db := utils.DB()
	if db.Find(&categories).Error != nil {
		res := utils.Response{Error: "Kategori tidak dapat ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	} else {
		res := utils.Response{Data: &categories}
		return c.JSON(http.StatusOK, &res)
	}
}
