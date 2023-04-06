package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
	"lelang.id/backend/utils"
)

func SaveStaff(c echo.Context) error {
	var user utils.Petugas
	c.Bind(&user)

	db := utils.DB()
	if db.Create(&user).Error != nil {
		res := utils.Response{Error: "Gagal menyimpan akun"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Message: "Akun tersimpan"}
		return c.JSON(http.StatusCreated, &res)
	}
}

func ListStaffs(c echo.Context) error {
	var users []utils.Petugas

	search := c.QueryParam("search")

	db := utils.DB()
	if db.Preload(clause.Associations).Where("LOWER(CONCAT(nama_depan, ' ', nama_belakang)) LIKE ?", fmt.Sprintf("%s%s%s", "%", strings.ToLower(search), "%")).Find(&users).Error != nil {
		res := utils.Response{Error: "Tidak ada petugas"}
		return c.JSON(http.StatusNotFound, &res)
	} else {
		res := utils.Response{Data: &users}
		return c.JSON(http.StatusOK, &res)
	}
}

func DeleteStaff(c echo.Context) error {
	staffId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db := utils.DB()
	if db.Where(&utils.Petugas{ID: staffId}).Delete(&utils.Petugas{}).Error != nil {
		res := utils.Response{Error: "Tidak dapat menghapus akun"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Message: "Akun terhapus"}
		return c.JSON(http.StatusOK, &res)
	}
}

func UpdateStaff(c echo.Context) error {
	staffId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var staff utils.Petugas
	c.Bind(&staff)
	staff.ID = staffId

	db := utils.DB()
	if db.Save(&staff).Error != nil {
		res := utils.Response{Error: "Tidak dapat menyimpan perubahan"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Message: "Perubahan tersimpan"}
		return c.JSON(http.StatusOK, &res)
	}
}
