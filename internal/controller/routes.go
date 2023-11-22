package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"example.com/dynamicWordpressBuilding/internal/middleware"
	"example.com/dynamicWordpressBuilding/utils"
	"example.com/dynamicWordpressBuilding/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (ctl *Controller) Routes() {
	ctl.Router.GET("/", home)
	ctl.publicRoutes()
	ctl.SetRoutes("/payment", "PAYMENT_SERVER_URL")
	//middleware
	authRouter := ctl.Router.Group("/").Use(middleware.AuthMiddleware(utils.NewTokenMaker()))
	ctl.privateRoutes(authRouter)
}

func (ctl *Controller) publicRoutes() {
	ctl.Router.POST("/user", ctl.CreateUser)
	ctl.Router.POST("/user/login", ctl.LoginUser)
	ctl.Router.POST("/user/reset-password", ctl.ResetPassword)
	ctl.Router.POST("/user/forget-password", ctl.ForgetPassword)
}

func (ctl *Controller) privateRoutes(authRouter gin.IRoutes) {
	authRouter.GET("/alluser", ctl.GetAllUser)
	authRouter.GET("/user/:id", ctl.GetUser)
	authRouter.DELETE("/user/:id", ctl.DeleteUser)
}

func (ctl *Controller) SetRoutes(path string, envValue string) {
	ctl.Router.Any(path+"/*proxyPath", func(c *gin.Context) {
		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Replace the request body for later use
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Construct the URL for the proxy request using the specified environment variable
		url := fmt.Sprintf("%s%s", os.Getenv(envValue), path+"/"+c.Param("proxyPath"))
		proxyReq, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		// Copy headers from the original request to the proxy request
		copyHeader(proxyReq.Header, c.Request.Header)

		// Extract user and organization IDs from the token

		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization is not provided")

			c.JSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			// ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s ", authorizationType)
			// ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err))
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}
		accessToken := fields[1]

		payload, err := ctl.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(err))
			return
		}
		// Set user-related headers if a user ID is present
		logrus.Info("payload ", payload)
		if payload.ID != 0 {
			user, _, _ := ctl.svc.GetUser(int(payload.ID))
			proxyReq.Header.Set("x-user-id", fmt.Sprint(user.ID))
			// if userGotten != nil && userGotten.IsAdmin {
			// 	proxyReq.Header.Set("x-user-role", "ADMIN")
			// }
		}

		// Perform the proxy request
		httpClient := http.Client{}
		logrus.Info(proxyReq)
		resp, err := httpClient.Do(proxyReq)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Read the response body from the proxy
		response, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Copy headers from the proxy response to the original response
		copyHeader(c.Writer.Header(), resp.Header)
		c.Status(resp.StatusCode)
		c.Writer.Write(response)
	})
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
