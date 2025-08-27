package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanuthatusu/tepia/internal/sessions"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		session, err := sessions.Store.Get(c.Request, "session")
		if err != nil {
			fmt.Println("err:", err.Error())
		}
		if _, ok := session.Values["user_id"]; !ok {
			fmt.Println("session:", session.Values)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
