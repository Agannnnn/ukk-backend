package main

import (
	"net/http"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
	"lelang.id/backend/router"
	"lelang.id/backend/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}

	utils.InitDB()  // Initializing database
	e := echo.New() // Initializing server

	e.Static("/assets", path.Clean("./public/images"))

	e.POST("/register", func(c echo.Context) error {
		var user utils.Masyarakat
		c.Bind(&user)
		db := utils.DB()

		if err := db.Create(&user).Error; err != nil {
			res := utils.Response{Error: "Gagal menyimpan akun"}
			return c.JSON(http.StatusInternalServerError, &res)
		} else {
			res := utils.Response{Message: "Akun berhasil disimpan"}
			return c.JSON(http.StatusCreated, &res)
		}
	})

	e.GET("/auctions/refresh", func(c echo.Context) error {
		db := utils.DB()

		var auctions []utils.Lelang

		db.Preload(clause.Associations).Where(utils.Lelang{Status: "dibuka"}).Find(&auctions)

		for _, auction := range auctions {
			if auction.SelesaiLelang.Compare(time.Now()) < 0 {
				db.Model(utils.Lelang{ID: auction.ID}).Update("status", "ditutup")
			} else if len(auction.Penawaran) > 0 {
				if time.Until(auction.Penawaran[len(auction.Penawaran)-1].Timestamp).Seconds() > float64(auction.Timeout) {
					db.Model(utils.Lelang{ID: auction.ID}).Update("status", "ditutup")
				}
			}
		}

		return c.NoContent(http.StatusOK)
	})

	router.StaffRouter(e.Group("/admin"))
	router.PublicRouter(e.Group("/public"))

	go refreshAuction()

	appPort := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(appPort))

}

func refreshAuction() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		http.Get("http://127.0.0.1:3000/auctions/refresh")
	}
}
