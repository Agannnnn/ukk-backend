package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lelang.id/backend/handler"
	"lelang.id/backend/utils"
)

func StaffRouter(e *echo.Group) {
	e.Use(middleware.BasicAuth(adminAuth))

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/auctions", handler.ListAuctions)
	e.GET("/auction/:id", handler.GetAuction)
	e.POST("/auction", handler.SaveAuction)
	e.PATCH("/auction", handler.UpdateAuction)
	e.DELETE("/auction/:id", handler.CloseAuction)

	e.GET("/categories", handler.ListCategories)
	e.POST("/category", handler.SaveCategory)

	e.GET("/products", handler.ListProducts)
	e.GET("/product/:id", handler.GetProduct)
	e.POST("/product", handler.SaveProduct)
	e.PATCH("/product/:id", handler.UpdateProduct)
	e.DELETE("/product/:id", handler.DeleteProduct)

	e.GET("/staffs", handler.ListStaffs)
	e.POST("/staff", handler.SaveStaff)
	e.PATCH("/staff", handler.UpdateStaff)
	e.DELETE("/staff/:id", handler.DeleteStaff)

	e.GET("/report/products", handler.ProductsReport)
	e.GET("/report/auctions", handler.AuctionsReport)
}

func adminAuth(username, password string, ctx echo.Context) (bool, error) {
	db := utils.DB()
	var user utils.Petugas
	if res := db.Where(&utils.Petugas{Username: username, Password: password}).First(&user); res.Error != nil || res.RowsAffected < 1 {
		return false, res.Error
	}

	userId := new(http.Cookie)
	userId.Name = "User-ID"
	userId.Value = user.ID.String()
	ctx.SetCookie(userId)

	isAdmin := new(http.Cookie)
	isAdmin.Name = "Is-Admin"
	isAdmin.Value = fmt.Sprintf("%v", strings.Contains(user.Level, "administrator"))
	ctx.SetCookie(isAdmin)

	return true, nil
}
