package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lelang.id/backend/utils"
)

func IsSubscribed(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}
	userIdCookie, err := c.Cookie("User-ID")
	if err != nil {
		res := utils.Response{Error: "Anda tidak memiliki hak"}
		return c.JSON(http.StatusUnauthorized, &res)
	}
	userId, err := uuid.Parse(userIdCookie.Value)
	if err != nil {
		res := utils.Response{Error: "ID user tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var langganan utils.Langganan

	db := utils.DB()
	if err := db.Where(utils.Langganan{LelangID: auctionId, MasyarakatID: userId}).First(&langganan).Error; err != nil {
		res := utils.Response{Error: "Pelelangan belum disubscribe"}
		return c.JSON(http.StatusNotFound, &res)
	} else {
		res := utils.Response{Data: &langganan}
		return c.JSON(http.StatusOK, &res)
	}
}

func Subscribe(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var userId uuid.UUID
	if cookieUserID, err := c.Cookie("User-ID"); err != nil {
		res := utils.Response{Error: "Mohon login terlebih dahulu"}
		return c.JSON(http.StatusUnauthorized, res)
	} else {
		userId, err = uuid.Parse(cookieUserID.Value)
		if err != nil {
			res := utils.Response{Error: "ID session tidak valid, mohon login ulang"}
			return c.JSON(http.StatusBadRequest, &res)
		}
	}

	auction := utils.Langganan{LelangID: auctionId, MasyarakatID: userId}
	db := utils.DB()
	if db.Create(&auction).Error != nil {
		res := utils.Response{Error: "Status berlangganan tidak berhasil disimpan"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Message: "Pelelangan diikuti"}
		return c.JSON(http.StatusCreated, &res)
	}
}

func Unsubscribe(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var userId uuid.UUID
	if cookieUserID, err := c.Cookie("User-ID"); err != nil {
		res := utils.Response{Error: "Mohon login terlebih dahulu"}
		return c.JSON(http.StatusUnauthorized, &res)
	} else {
		userId, err = uuid.Parse(cookieUserID.Value)
		if err != nil {
			res := utils.Response{Error: "ID session tidak valid, mohon login ulang"}
			return c.JSON(http.StatusBadRequest, &res)
		}
	}

	auction := &utils.Langganan{LelangID: auctionId, MasyarakatID: userId}
	db := utils.DB()
	if db.Where(&auction).Delete(&auction).Error != nil {
		res := utils.Response{Error: "Terjadi kesalahan"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Message: "Pelelangan dilepaskan"}
		return c.JSON(http.StatusOK, &res)
	}
}
