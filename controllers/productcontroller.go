package controllers

import (
	"goLang/database"
	"goLang/entities"
	"net/http"

	"github.com/labstack/echo"
)

func CreateProduct(c echo.Context) error {
	product := new(entities.Product)
	if err := c.Bind(product); err != nil {
		return err
	}
	database.Instance.Create(&product)
	return c.JSON(http.StatusCreated, product)
}

func GetProductById(c echo.Context) error {
	productId := c.Param("id")
	if checkIfProductExists(productId) == false {
		return c.JSON(http.StatusNotFound, "product Not Found")
	}
	var product entities.Product
	database.Instance.First(&product, productId)
	return c.JSON(http.StatusOK, product)
}

func GetProducts(c echo.Context) error {
	var products []entities.Product
	database.Instance.Find(&products)
	return c.JSON(http.StatusOK, products)
}

func UpdateProduct(c echo.Context) error {
	productId := c.Param("id")
	if checkIfProductExists(productId) == false {
		return c.JSON(http.StatusNotFound, "product Not Found")
	}
	product := entities.Product{}
	if err := c.Bind(&product); err != nil {
		return err
	}
	database.Instance.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	productId := c.Param("id")
	if checkIfProductExists(productId) == false {
		return c.JSON(http.StatusNotFound, "product Not Found")
	}
	product := new(entities.Product)
	if err := c.Bind(product); err != nil {
		return err
	}
	database.Instance.Delete(&product, productId)
	return c.JSON(http.StatusOK, product)
}

func checkIfProductExists(productId string) bool {
	var product entities.Product
	database.Instance.First(&product, productId)
	if product.ID == 0 {
		return false
	}
	return true
}
