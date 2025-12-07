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

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		userID := sess.Values["user_id"]

		if userID == nil {
			return c.JSON(401, map[string]string{"error": "unauthorized"})
		}

		// можешь сохранить userID в контекст
		c.Set("userID", userID)

		return next(c)
	}
}

func main() {
	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	CachedProducts = loadProducts()

	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(&userRepo)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/products", getProducts, RequireAuth)

	handlers.Auth(e, authService)

	e.Logger.Fatal(e.Start(":1323"))
}
