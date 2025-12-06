package main

import (
	"backApp/handlers"
	"backApp/models"
	"backApp/repository"
	"backApp/services"
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
	if start > len(CachedProducts) {
		return c.String(http.StatusNotFound, "Page out of range")
	}

	end = min(page*limit, len(CachedProducts))
	return c.JSON(http.StatusOK, CachedProducts[start:end])

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

var CachedProducts []models.Product

func main() {
	e := echo.New()

	CachedProducts = loadProducts()

	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(&userRepo)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/products", getProducts)

	handlers.Auth(e, authService)

	e.Logger.Fatal(e.Start(":1323"))
}
