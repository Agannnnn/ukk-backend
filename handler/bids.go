package handler

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm/clause"
	"lelang.id/backend/utils"
)

func GetBid(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var bid utils.Penawaran
	db := utils.DB()
	if db.Where(&utils.Penawaran{LelangID: auctionId}).Last(&bid).Error != nil {
		res := utils.Response{Error: "Tidak ada bid tersimpan"}
		return c.JSON(http.StatusNotFound, &res)
	} else {
		res := utils.Response{Data: &bid}
		return c.JSON(http.StatusOK, &res)
	}
}

func SaveBid(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
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

	var bid utils.Penawaran
	c.Bind(&bid)
	bid.LelangID = auctionId
	bid.MasyarakatID = userId

	db := utils.DB()

	var auction utils.Lelang
	if db.Where(utils.Lelang{ID: auctionId}).First(&auction).Error != nil {
		res := utils.Response{Error: "Pelelangan tidak ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	}

	var latestBid utils.Penawaran
	if err := db.Where(&utils.Penawaran{LelangID: auctionId}).Last(&latestBid).Error; err != nil {
		bid.Harga += auction.HargaAwal
		if db.Create(&bid).Error != nil {
			res := utils.Response{Error: "Gagal menyimpan penawaran"}
			return c.JSON(http.StatusInternalServerError, &res)
		} else {
			res := utils.Response{Message: "Penawaran berhasil disimpan"}
			return c.JSON(http.StatusCreated, &res)
		}
	} else {
		bid.Harga += latestBid.Harga
		if db.Create(&bid).Error != nil {
			res := utils.Response{Error: "Gagal menyimpan penawaran"}
			return c.JSON(http.StatusInternalServerError, &res)
		} else {
			res := utils.Response{Message: "Penawaran berhasil disimpan"}
			return c.JSON(http.StatusCreated, &res)
		}
	}
}

func PaymentStatus(c echo.Context) error {
	auctionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db := utils.DB()

	var auction utils.Lelang
	db.Preload(clause.Associations).First(&auction, &utils.Lelang{ID: auctionId})

	if len(auction.Penawaran) == 0 {
		return c.JSON(http.StatusNotFound, &utils.Response{Error: "Pelelangan tidak memiliki penawaran"})
	}

	midtransClient := coreapi.Client{}
	midtransClient.New(os.Getenv("MIDTRANS-ENV"), midtrans.Sandbox)

	midtransExistedTransaction, _ := midtransClient.CheckTransaction(auctionId.String())

	if midtransExistedTransaction.TransactionStatus == "settlement" {
		return c.JSON(http.StatusOK, &utils.Response{Message: "Pelelangan sudah dibayar"})
	}

	return c.JSON(http.StatusNotFound, &utils.Response{Error: "Pelelangan belum dibayar"})
}

func PayAuction(c echo.Context) error {
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

	db := utils.DB()

	var auction utils.Lelang
	db.Preload(clause.Associations).First(&auction, &utils.Lelang{ID: auctionId})

	if len(auction.Penawaran) == 0 {
		return c.JSON(http.StatusNotFound, &utils.Response{Error: "Pelelangan tidak memiliki penawaran"})
	}

	var userInput utils.Masyarakat
	c.Bind(&userInput)

	var user utils.Masyarakat
	db.Preload(clause.Associations).First(&user, &utils.Masyarakat{ID: userId})

	if userInput.NoTelp != "" {
		user.NoTelp = userInput.NoTelp
	}

	midtransClient := coreapi.Client{}
	midtransClient.New(os.Getenv("MIDTRANS-ENV"), midtrans.Sandbox)

	midtransExistedTransaction, _ := midtransClient.CheckTransaction(auctionId.String())
	if midtransExistedTransaction.StatusCode != "404" {
		midtransClient.CancelTransaction(auctionId.String())
	}

	midtransTransaction, _ := midtransClient.ChargeTransaction(&coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeGopay,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  auctionId.String(),
			GrossAmt: int64(auction.Penawaran[len(auction.Penawaran)-1].Harga),
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    auction.Barang.ID.String(),
				Name:  auction.Barang.Nama,
				Price: int64(auction.Penawaran[len(auction.Penawaran)-1].Harga),
				Qty:   1,
			},
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.NamaDepan,
			LName: user.NamaBelakang,
			Email: user.Email,
			Phone: user.NoTelp,
		},
	})

	if midtransTransaction.StatusCode == "201" {
		return c.JSON(http.StatusCreated, &utils.Response{Data: &midtransTransaction.Actions})
	}
	return c.JSON(http.StatusBadGateway, &utils.Response{Error: "Gagal menyimpan pembayaran"})
}
