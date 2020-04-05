package server

import (
	"database/sql"
	"github.com/GreatGodApollo/als/crypto"
	"github.com/GreatGodApollo/als/database"
	"github.com/GreatGodApollo/als/models"
	"github.com/GreatGodApollo/als/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var db *sql.DB

func Setup(d *sql.DB) {
	db = d
}

func RunAPI() {
	if viper.GetBool("server.production") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/", IndexRouter)

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/", gin.BasicAuth(viper.GetStringMapString("auth.accounts")))
			{
				auth.POST("/create", CreateRouter)
				auth.POST("/invalidate", InvalidateRouter)
				auth.POST("/specific", GetRouter)
				auth.GET("/all/:product", GetAllRouter)
			}
		}
	}

	license := r.Group("/license")
	{
		license.POST("/check", CheckRouter)
	}

	license.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP()
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(1*time.Minute), 10), time.Hour
	}, func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "you have reached your limit!",
			"code":    http.StatusTooManyRequests,
		})
	}))

	r.NoRoute(NotFoundRouter)

	r.Run(viper.GetString("server.bind"))
}

func NotFoundRouter(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "404: not found",
	})
}

func IndexRouter(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "apollo's licensing server",
	})
}

func CreateRouter(c *gin.Context) {
	var req models.LicenseRequest

	if c.ShouldBind(&req) == nil {
		crypt, err := utils.GenerateEncryptedLicense(db, req.Product, req.Email)
		if handleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, models.LicenseResponse{
			LicenseKey: crypto.EncodeBase64(crypt),
			Status:     "created",
			Message:    "license created",
			Code:       http.StatusCreated,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "required parameters were not provided!",
			"code":    http.StatusBadRequest,
		})
	}
}

func InvalidateRouter(c *gin.Context) {
	var req models.BasicRequest
	if c.ShouldBind(&req) == nil {
		// Decode and decrypt key
		enc, err := crypto.DecodeBase64(req.Key)
		if handleError(c, err) {
			return
		}
		key, err := utils.DecryptLicense(enc)
		if handleError(c, err) {
			return
		}

		exist, err := database.CheckLicenseExist(db, string(key))
		if handleError(c, err) {
			return
		}

		exist, valid, err := database.CheckLicenseValid(db, string(key))
		if handleLicenseError(c, req.Key, err) {
			return
		}

		if exist && valid {
			_, err := database.InvalidateLicense(db, string(key))
			if handleError(c, err) {
				return
			}

			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "invalidated",
				Message:    "license invalidated",
				Code:       http.StatusOK,
			})
		} else if exist {
			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "invalid",
				Message:    "license already invalid",
				Code:       http.StatusOK,
			})
		} else {
			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "invalid",
				Message:    "license nonexistent",
				Code:       http.StatusOK,
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "required parameters not provided",
			"code":    http.StatusBadRequest,
		})
	}
}

func GetRouter(c *gin.Context) {
	var req models.BasicRequest
	if c.ShouldBind(&req) == nil {
		enc, err := crypto.DecodeBase64(req.Key)
		if handleError(c, err) {
			return
		}
		key, err := utils.DecryptLicense(enc)
		if handleError(c, err) {
			return
		}

		licObj, err := database.GetWholeRecord(db, string(key))
		if handleLicenseError(c, req.Key, err) {
			return
		}

		licObj.LicenseKey = req.Key
		licObj.Code = http.StatusOK

		c.JSON(http.StatusOK, licObj)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "required parameters not provided",
			"code":    http.StatusBadRequest,
		})
	}
}

func GetAllRouter(c *gin.Context) {
	objects, err := database.GetAllValidRecords(db, c.Param("product"))
	if handleError(c, err) {
		return
	}
	objects.Code = http.StatusOK
	c.JSON(http.StatusOK, objects)
}

func CheckRouter(c *gin.Context) {
	var req models.CheckRequest

	if c.ShouldBind(&req) == nil {

		// Decode & Decrypt Key
		enc, err := crypto.DecodeBase64(req.Key)
		if handleError(c, err) {
			return
		}
		key, err := utils.DecryptLicense(enc)
		if handleError(c, err) {
			return
		}

		// Check if key exists in DB
		exist, err := database.CheckLicenseExist(db, string(key))
		if handleError(c, err) {
			return
		}

		// Check if valid
		exist, valid, err := database.CheckLicenseValidProduct(db, string(key), req.Product)
		if handleLicenseError(c, req.Key, err) {
			return
		}

		if exist && valid {
			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "valid",
				Message:    "license valid",
				Code:       http.StatusOK,
			})
		} else if exist {
			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "invalid",
				Message:    "license invalid",
				Code:       http.StatusOK,
			})
		} else {
			c.JSON(http.StatusNotFound, models.LicenseResponse{
				LicenseKey: req.Key,
				Status:     "invalid",
				Message:    "license nonexistent",
				Code:       http.StatusNotFound,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "required parameters not provided",
			"code":    http.StatusBadRequest,
		})
	}
}

func handleLicenseError(c *gin.Context, license string, err error) bool {
	if err != nil {
		if err.Error() == "incorrect product" || err.Error() == "license nonexistent" {
			c.JSON(http.StatusOK, models.LicenseResponse{
				LicenseKey: license,
				Status:     "invalid",
				Message:    err.Error(),
				Code:       http.StatusOK,
			})
			return true
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return true
	}
	return false
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return true
	}
	return false
}
