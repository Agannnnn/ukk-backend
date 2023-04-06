package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lelang.id/backend/utils"
)

func ListProducts(c echo.Context) error {
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")

	var sortQuery string
	switch sort {
	case "upload_baru":
		sortQuery = "diunggah DESC"
	case "upload_lama":
		sortQuery = "diunggah ASC"
	default:
		sortQuery = "nama ASC"
	}

	var products []utils.Barang

	db := utils.DB()
	if err := db.Preload(clause.Associations).Where("LOWER(nama) LIKE ?", fmt.Sprintf("%s%s%s", "%", strings.ToLower(search), "%")).Order(sortQuery).Find(&products).Error; err != nil {
		res := utils.Response{Error: "Barang tidak ditemukan", Data: []string{}}
		return c.JSON(http.StatusNotFound, &res)
	}

	res := utils.Response{Data: &products}
	return c.JSON(http.StatusOK, &res)
}

func GetProduct(c echo.Context) error {
	productId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var products utils.Barang
	db := utils.DB()
	if err := db.Preload(clause.Associations).First(&products, productId).Error; err != nil {
		res := utils.Response{Error: "Barang tidak dapat ditemukan"}
		return c.JSON(http.StatusInternalServerError, &res)
	} else {
		res := utils.Response{Data: &products}
		return c.JSON(http.StatusOK, &res)
	}
}

func SaveProduct(c echo.Context) error {
	var product utils.Barang
	if err := c.Bind(&product); err != nil {
		res := utils.Response{Error: "Input tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db := utils.DB()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		if categories := strings.Split(c.FormValue("categories"), ";")[:]; len(categories) > 0 {
			for _, category := range categories {
				categoryId, err := uuid.Parse(category)
				if err != nil {
					continue
				}
				if err := tx.Create(&utils.DetailKategori{BarangID: product.ID, KategoriID: categoryId}).Error; err != nil {
					return err
				}
			}
		}

		form, err := c.MultipartForm()
		if err == nil {
			images := form.File["images"]
			for i, image := range images {
				imageDesc := c.FormValue(fmt.Sprintf("image_desc[%d]", i))
				fileName := fmt.Sprintf("%v%s", time.Now().Unix(), path.Ext(image.Filename))
				if err := tx.Create(&utils.FotoBarang{Filename: fileName, BarangID: product.ID, Deskripsi: imageDesc}).Error; err != nil {
					return err
				}

				src, err := image.Open()
				if err != nil {
					return err
				}
				defer src.Close()

				dst, err := os.Create(fmt.Sprintf("./public/images/%v", fileName))
				if err != nil {
					return err
				}
				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		res := utils.Response{Error: "Gagal menyimpan barang"}
		return c.JSON(http.StatusBadRequest, &res)
	} else {
		res := utils.Response{Message: "Barang tersimpan"}
		return c.JSON(http.StatusCreated, &res)
	}
}

func UpdateProduct(c echo.Context) error {
	productId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID produk tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	var product utils.Barang
	if err := c.Bind(&product); err != nil {
		res := utils.Response{Error: "Input tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}
	product.ID = productId

	db := utils.DB()
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(utils.Barang{ID: productId}).Updates(&product).Error; err != nil {
			return err
		}

		tx.Where(&utils.DetailKategori{BarangID: productId}).Delete(&utils.DetailKategori{})
		if categories := strings.Split(c.FormValue("categories"), ";")[:]; len(categories) > 0 {
			for _, category := range categories {
				categoryId, err := uuid.Parse(category)
				if err != nil {
					continue
				}
				if err := tx.Create(&utils.DetailKategori{BarangID: product.ID, KategoriID: categoryId}).Error; err != nil {
					return err
				}
			}
		}

		form, err := c.MultipartForm()
		if err == nil {
			Inputimages := form.File["images"]
			if len(Inputimages) > 0 {
				var images []utils.FotoBarang
				tx.Where(&utils.FotoBarang{BarangID: productId}).Find(&images)
				for _, image := range images {
					os.Remove(fmt.Sprintf("./public/images/%v", image.Filename))
				}
				tx.Where(&utils.FotoBarang{BarangID: productId}).Delete(&utils.FotoBarang{})

				for i, image := range Inputimages {
					imageDesc := c.FormValue(fmt.Sprintf("image_desc[%d]", i))
					fileName := fmt.Sprintf("%v%s", time.Now().Unix(), path.Ext(image.Filename))
					if err := tx.Create(&utils.FotoBarang{Filename: fileName, BarangID: product.ID, Deskripsi: imageDesc}).Error; err != nil {
						return err
					}

					src, err := image.Open()
					if err != nil {
						return err
					}
					defer src.Close()

					dst, err := os.Create(fmt.Sprintf("./public/images/%v", fileName))
					if err != nil {
						return err
					}
					defer dst.Close()

					if _, err = io.Copy(dst, src); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		res := utils.Response{Error: "Barang tidak dapat diperbarui"}
		return c.JSON(http.StatusBadRequest, &res)
	} else {
		res := utils.Response{Message: "Barang telah diperbarui"}
		return c.JSON(http.StatusOK, &res)
	}
}

func DeleteProduct(c echo.Context) error {
	productId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		res := utils.Response{Error: "ID tidak valid"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db := utils.DB()

	product := utils.Barang{ID: productId}
	if err := db.Preload(clause.Associations).First(&product).Error; err != nil {
		res := utils.Response{Error: "Barang tidak dapat ditemukan"}
		return c.JSON(http.StatusNotFound, &res)
	}

	if len(product.Lelang) > 0 {
		res := utils.Response{Error: "Tidak dapat menghapus barang yang sudah pernah dilelang"}
		return c.JSON(http.StatusBadRequest, &res)
	}

	db.Delete(&utils.Barang{ID: productId})

	var images []utils.FotoBarang
	if db.Where(&utils.FotoBarang{BarangID: productId}).Find(&images).Error == nil {
		for _, image := range images {
			os.Remove(fmt.Sprintf("./public/images/%s", image.Filename))
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Barang telah dihapus"})
}

func ProductsReport(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	fromDate, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": []string{"Tanggal 'dari' tidak valid", err.Error()}})
	}
	toDate, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Error": []string{"Tanggal 'dari' tidak valid", err.Error()}})
	}

	var products []utils.Barang

	db := utils.DB()
	if err := db.Where("diunggah BETWEEN ? AND ?", fromDate, toDate).Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"Error": "Terjadi kesalahan saat mengambil data barang"})
	} else {
		var productsStr []string

		for _, product := range products {
			var productsImage []string
			for _, image := range product.FotoBarang {
				productsImage = append(productsImage, image.Filename)
			}
			productStr := fmt.Sprintf(
				"ID %v Nama Barang %s Foto Barang [%v] Diunggah pada %v",
				product.ID,
				product.Nama,
				product.Diunggah,
				strings.Join(productsImage, ", "),
			)
			productsStr = append(productsStr, productStr)
		}

		return c.Blob(http.StatusOK, "text/plain", []byte(strings.Join(productsStr, "\n-\n")))
	}
}
