package main

import (
	"backApp/models"
	"encoding/json"
	"net/http"
	"os"

	"strconv"

	"github.com/labstack/echo/v4"
)

func getProducts(c echo.Context) error {
	pageParam := c.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}
	const limit = 3

	var start, end int
	start = (page - 1) * limit
	if start > len(Products) {
		return c.String(http.StatusNotFound, "Page out of range")
	}

	end = min(page*limit, len(Products))
	return c.JSON(http.StatusOK, Products[start:end])

}

func loadProducts() []models.Product {
	file, err := os.Open("products.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var products []models.Product
	err = decoder.Decode(&products)
	if err != nil {
		panic(err)
	}
	return products
}

var Products []models.Product

func main() {
	e := echo.New()

	Products = loadProducts()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/products", getProducts)
	e.Logger.Fatal(e.Start(":1323"))
}
