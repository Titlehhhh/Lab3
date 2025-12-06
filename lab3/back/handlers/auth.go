package handlers

import (
	"backApp/services"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Auth(e *echo.Echo, authService services.AuthService) {
	e.POST("/register", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		err := authService.Register(username, password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "registered"})
	})

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		userID, err := authService.Login(username, password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		sess, _ := session.Get("session", c)
		sess.Values["user_id"] = userID
		sess.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, map[string]string{"status": "logged in"})
	})

	e.POST("/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		delete(sess.Values, "user_id")
		sess.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, map[string]string{"status": "logged out"})
	})
}
