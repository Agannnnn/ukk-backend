package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lelang.id/backend/handler"
	"lelang.id/backend/utils"
)

func PublicRouter(e *echo.Group) {
	e.Use(middleware.BasicAuth(publicAuth))

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/auctions", handler.ListAuctions)
	e.GET("/auction/:id", handler.GetAuction)

	e.GET("/products", handler.ListProducts)
	e.GET("/product/:id", handler.GetProduct)

	e.GET("/bid/:id", handler.GetBid)
	e.POST("/bid/:id", handler.SaveBid)

	e.GET("/paid/:id", handler.PaymentStatus)
	e.POST("/pay/:id", handler.PayAuction)

	e.GET("/subscribe/:id", handler.IsSubscribed)
	e.POST("/subscribe/:id", handler.Subscribe)
	e.DELETE("/unsubscribe/:id", handler.Unsubscribe)

	e.GET("/profile", handler.GetProfile)
	e.PATCH("/profile", handler.UpdateProfile)

	e.GET("/categories", handler.ListCategories)
}

func publicAuth(username, password string, ctx echo.Context) (bool, error) {
	db := utils.DB()
	var user utils.Masyarakat
	db.Where(&utils.Masyarakat{Username: username, Password: password}).Find(&user)
	if err := db.Where(&utils.Masyarakat{Username: username, Password: password}).First(&user).Error; err != nil {
		return false, err
	}

	userId := new(http.Cookie)
	userId.Name = "User-ID"
	userId.Value = user.ID.String()
	ctx.SetCookie(userId)

	return true, nil
}
