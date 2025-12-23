package main

import (
	"backApp/handlers"
	"backApp/models"
	"backApp/repository"
	"backApp/services"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const productsPerPage = 3

func RequireAuth(userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			userID, ok := sess.Values["user_id"].(int)
			if !ok || userID == 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok || user.Role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}
		return next(c)
	}
}

func currentUser(c echo.Context) (*models.User, bool) {
	user, ok := c.Get("user").(*models.User)
	return user, ok
}

func main() {
	e := echo.New()

	// Sessions
	e.Use(session.Middleware(
		sessions.NewCookieStore([]byte("secret")),
	))

	// Repositories & services
	userRepo := repository.NewUserRepository()
	authService := services.NewAuthService(&userRepo)

	productService, err := services.NewProductService("products.json")
	if err != nil {
		e.Logger.Fatal(err)
	}

	cartService, err := services.NewCartService("carts.json")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Static images
	e.Static("/images", "static/images")

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Backend is running")
	})

	//products
	e.GET("/products", func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}

		products := productService.List()
		start := (page - 1) * productsPerPage
		if start >= len(products) && len(products) != 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "page out of range"})
		}

		end := min(start+productsPerPage, len(products))
		return c.JSON(http.StatusOK, products[start:end])
	}, RequireAuth(userRepo))

	e.GET("/products/all", func(c echo.Context) error {
		return c.JSON(http.StatusOK, productService.List())
	}, RequireAuth(userRepo), RequireAdmin)

	e.POST("/products", func(c echo.Context) error {
		var product models.Product
		if err := c.Bind(&product); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		}

		created, err := productService.Add(product)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, created)
	}, RequireAuth(userRepo), RequireAdmin)

	e.PUT("/products/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}

		var updates models.Product
		if err := c.Bind(&updates); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		}

		updated, err := productService.Update(id, updates)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, updated)
	}, RequireAuth(userRepo), RequireAdmin)

	e.DELETE("/products/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}

		if err := productService.Delete(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	}, RequireAuth(userRepo), RequireAdmin)

	// User and Cart
	e.GET("/me", func(c echo.Context) error {
		user, ok := currentUser(c)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		return c.JSON(http.StatusOK, user)
	}, RequireAuth(userRepo))

	e.GET("/cart", func(c echo.Context) error {
		user, _ := currentUser(c)
		items := cartService.GetCart(user.Id)

		result := make([]map[string]any, 0)
		for _, item := range items {
			product, err := productService.GetByID(item.ProductID)
			if err != nil {
				continue
			}
			result = append(result, map[string]any{
				"product":  product,
				"quantity": item.Quantity,
			})
		}
		return c.JSON(http.StatusOK, result)
	}, RequireAuth(userRepo))

	e.POST("/cart", func(c echo.Context) error {
		user, _ := currentUser(c)

		var payload struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		}

		if payload.Quantity <= 0 {
			payload.Quantity = 1
		}

		if err := cartService.AddItem(user.Id, payload.ProductID, payload.Quantity); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "added"})
	}, RequireAuth(userRepo))

	e.DELETE("/cart/:productId", func(c echo.Context) error {
		user, _ := currentUser(c)

		productID, err := strconv.Atoi(c.Param("productId"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}

		if err := cartService.RemoveItem(user.Id, productID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	}, RequireAuth(userRepo))

	handlers.Auth(e, authService)

	e.Logger.Fatal(e.Start(":1323"))
}
