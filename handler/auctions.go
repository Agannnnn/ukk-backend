package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
	"lelang.id/backend/utils"
)

func ListAuctions(c echo.Context) error {
	filter := c.QueryParam("filter")
	search := c.QueryParam("search")
	category := c.QueryParam("category")
	subscribed := c.QueryParam("subscribed")

	if search == "" {
		search = "_"
	}

	var filterQuerry string
	switch filter {
	case "mulai_desc":
		filterQuerry = "mulai_lelang DESC"
	case "berakhir_asc":
		filterQuerry = "selesai_lelang ASC"
	case "berakhir_desc":
		filterQuerry = "selesai_lelang DESC"
	case "harga_awal":
		filterQuerry = "harga_awal DESC"
	default:
		filterQuerry = "mulai_lelang ASC"
	}

	var auctions []utils.Lelang

	db := utils.DB()

	querry := db.
		Model(&utils.Lelang{}).
		Preload("Barang.FotoBarang").
		Preload("Barang.Kategori.Kategori").
		Preload("Langganan").
		Preload(clause.Associations).
		Joins("Barang").
		Where(`"Barang"."nama" LIKE ?`, fmt.Sprintf("%%%s%%", search)).
		Order(filterQuerry)

	if err := querry.Find(&auctions).Error; err != nil {
		res := utils.Response{Error: "Pelelangan tidak ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	}

	if subscribed == "true" {
		userIdCookie, err := c.Cookie("User-ID")
		if err != nil {
			res := utils.Response{Error: "Mohon login terlebih dahulu"}
			return c.JSON(http.StatusUnauthorized, &res)
		}
		userId, err := uuid.Parse(userIdCookie.Value)
		if err != nil {
			res := utils.Response{Error: "ID user tidak valid"}
			return c.JSON(http.StatusBadRequest, &res)
		}

		var newAuctions []utils.Lelang
		for _, auction := range auctions {
			for _, sub := range auction.Langganan {
				if sub.MasyarakatID == userId {
					newAuctions = append(newAuctions, auction)
					break
				}
			}
		}

		auctions = newAuctions
	}

	if category != "" {
		categoryId, err := uuid.Parse(category)
		if err != nil {
			res := utils.Response{Error: "Kategori tidak valid"}
			return c.JSON(http.StatusOK, &res)
		}

		var newAuctions []utils.Lelang

		for _, auction := range auctions {
			for _, c := range auction.Barang.Kategori {
				if c.Kategori.ID == categoryId {
					newAuctions = append(newAuctions, auction)
					break
				}
			}
		}

		auctions = newAuctions
	}

	res := utils.Response{Data: &auctions}
	return c.JSON(http.StatusOK, &res)
}

func GetAuction(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID pelelangan tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var auction utils.Lelang
	db := utils.DB()
	if err := db.Preload(clause.Associations).Preload("Barang.FotoBarang").First(&auction, auctionId).Error; err != nil {
		res := utils.Response{Error: "Pelelangan tidak ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	} else {
		res := utils.Response{Data: &auction}
		return c.JSON(http.StatusOK, &res)
	}
}

func SaveAuction(c echo.Context) error {
	var auction utils.Lelang
	if err := c.Bind(&auction); err != nil {
		res := utils.Response{Error: "Terjadi kesalahan"}
		return c.JSON(http.StatusInternalServerError, &res)
	}

	var userId uuid.UUID
	if cookieUserID, err := c.Cookie("User-ID"); err != nil {
		res := utils.Response{Error: "Mohon login terlebih dahulu"}
		return c.JSON(http.StatusUnauthorized, &res)
	} else {
		userId, err = uuid.Parse(cookieUserID.Value)
		if err != nil {
			res := utils.Response{Error: "ID User tidak valid, mohon login ulang"}
			return c.JSON(http.StatusBadRequest, &res)
		}
	}

	db := utils.DB()
	auction.PetugasPembuatID = userId
	if err := db.Create(&auction).Error; err != nil {
		res := utils.Response{Error: "Pelelangan tidak berhasil dibuat"}
		return c.JSON(http.StatusBadRequest, &res)
	} else {
		res := utils.Response{Message: "Pelelangan berhasil dibuat"}
		return c.JSON(http.StatusCreated, &res)
	}
}

func UpdateAuction(c echo.Context) error {
	var auction utils.Lelang
	if err := c.Bind(&auction); err != nil {
		res := utils.Response{Error: "Terjadi kesalahan"}
		return c.JSON(http.StatusInternalServerError, &res)
	}

	db := utils.DB()
	if err := db.Select("mulai_lelang", "selesai_lelang", "timeout", "harga_awal", "min_penawaran").Updates(&auction).Error; err != nil {
		res := utils.Response{Error: "Pelelangan gagal diperbarui"}
		return c.JSON(http.StatusBadRequest, &res)
	} else {
		res := utils.Response{Message: "Pelelangan berhasil diperbarui"}
		return c.JSON(http.StatusOK, &res)
	}
}

func CloseAuction(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID pelelangan tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db := utils.DB()

	auction := utils.Lelang{ID: auctionId}
	if err := db.First(&auction).Error; err != nil {
		res := utils.Response{Error: "Lelang tidak dapat ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	}

	userIdCookie, err := c.Cookie("User-ID")
	if err != nil {
		res := utils.Response{Error: "Mohon login terlebih dahulu"}
		return c.JSON(http.StatusUnauthorized, &res)
	}
	userId, err := uuid.Parse(userIdCookie.Value)
	if err != nil {
		res := utils.Response{Error: "ID user tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	if db.Where(&auction).Updates(utils.Lelang{Status: "ditutup", PetugasPenutupID: userId}).Error != nil {
		res := utils.Response{Error: "Tidak dapat memperbarui status"}
		return c.JSON(http.StatusInternalServerError, &res)
	}

	res := utils.Response{Message: "Status telah diperbarui"}
	return c.JSON(http.StatusOK, &res)
}

func AuctionsReport(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	fromDate, err := time.Parse(time.RFC3339, from)
	if err != nil {
		res := utils.Response{Error: "Tanggal 'dari' tidak  valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}
	toDate, err := time.Parse(time.RFC3339, to)
	if err != nil {
		res := utils.Response{Error: "Tanggal 'hingga' tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var auctions []utils.Lelang

	db := utils.DB()
	if err := db.Preload(clause.Associations).Where("mulai_lelang > ? OR selesai_lelang < ?", fromDate, toDate).Find(&auctions).Error; err != nil {
		res := utils.Response{Error: "Terjadi kesalahan"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		var auctionsStr []string

		for _, auction := range auctions {
			var auctionStr string
			if len(auction.Penawaran) > 0 {
				auctionStr = fmt.Sprintf("ID %v | ID Barang %v | Dimulai pada %v sampai %v | Status %s | Dibuat oleh %v | Diakhiri oleh %v | ID penawaran terakhir %v dengan nominal %v | Dimenangkan oleh %v",
					auction.ID,
					auction.Barang.ID,
					auction.MulaiLelang,
					auction.SelesaiLelang,
					auction.Status,
					auction.PetugasPembuat.ID,
					auction.PetugasPenutup.ID,
					auction.Penawaran[len(auction.Penawaran)-1].ID,
					auction.Penawaran[len(auction.Penawaran)-1].Harga,
					auction.Penawaran[len(auction.Penawaran)-1].MasyarakatID,
				)
			} else {
				auctionStr = fmt.Sprintf("ID %v | ID Barang %v | Dimulai pada %v sampai %v | Status %s | Dibuat oleh %v | Diakhiri oleh %v | Harga awal %v",
					auction.ID,
					auction.Barang.ID,
					auction.MulaiLelang,
					auction.SelesaiLelang,
					auction.Status,
					auction.PetugasPembuat.ID,
					auction.PetugasPenutup.ID,
					auction.HargaAwal,
				)

			}
			auctionsStr = append(auctionsStr, auctionStr)
		}

		return c.Blob(http.StatusOK, "text/plain", []byte(strings.Join(auctionsStr, "\n-\n")))
	}
}
