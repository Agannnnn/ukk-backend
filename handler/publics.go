package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
	"lelang.id/backend/utils"
)

func GetProfile(c echo.Context) error {
	userIdCookie, err := c.Cookie("User-ID")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Anda tidak memiliki hak"})
	}
	userId, err := uuid.Parse(userIdCookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "ID user tidak valid"})
	}

	var user utils.Masyarakat

	db := utils.DB()
	if db.Preload(clause.Associations).Where(&utils.Masyarakat{ID: userId}).Find(&user).Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err})
	} else {
		return c.JSON(http.StatusOK, &user)
	}
}

func UpdateProfile(c echo.Context) error {
	userIdCookie, err := c.Cookie("User-ID")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Anda tidak memiliki hak"})
	}
	userId, err := uuid.Parse(userIdCookie.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "ID user tidak valid"})
	}

	var user utils.Masyarakat
	user.ID = userId

	c.Bind(&user)

	db := utils.DB()
	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Terjadi kesalahan"})
	} else {
		return c.JSON(http.StatusOK, echo.Map{"message": "Data baru tersimpan"})
	}
}
