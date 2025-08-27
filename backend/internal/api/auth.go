package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xanuthatusu/tepia/internal/db"
	"github.com/xanuthatusu/tepia/internal/models"
	"github.com/xanuthatusu/tepia/internal/sessions"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	r.POST("/register", func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		if err := db.CreateUserWithPassword(ctx, pool, user.Name, user.Email, string(hash)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	})

	r.POST("/login", func(c *gin.Context) {
		var bodyUser models.User
		if err := c.ShouldBindJSON(&bodyUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		user, err := db.GetUserByEmail(ctx, pool, bodyUser.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyUser.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		session, _ := sessions.Store.Get(c.Request, "session")
		session.Values["user_id"] = user.ID.String()
		if err := session.Save(c.Request, c.Writer); err != nil {
			fmt.Println("err: ", err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"message": "logged in"})
	})
}
