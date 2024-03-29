package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-vilt-manager/graph"
	generated "github.com/zicops/zicops-vilt-manager/graph"

	"github.com/zicops/zicops-vilt-manager/lib/jwt"
)

// CCRouter ... the router for the controller
func CCRouter(restRouter *gin.Engine) (*gin.Engine, error) {
	// set up cors
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	restRouter.Use(cors.New(configCors))
	// user a middleware to get context values
	restRouter.Use(func(c *gin.Context) {
		currentRequest := c.Request
		incomingToken := jwt.GetToken(currentRequest)
		claimsFromToken, _ := jwt.GetClaims(incomingToken)
		c.Set("zclaims", claimsFromToken)

	})
	restRouter.GET("/healthz", HealthCheckHandler)
	// create group for restRouter
	version1 := restRouter.Group("/api/v1")
	version1.POST("/query", graphqlHandler())
	return restRouter, nil
}

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		ctxValue := c.Value("zclaims").(map[string]interface{})
		// set ctxValue to request context
		lspIdInt := c.Request.Header.Get("tenant")
		lspID := ""
		if lspIdInt != "" {
			lspID = lspIdInt
		}
		ctxValue["lsp_id"] = lspID
		// get current origin in https format
		origin := c.Request.Header.Get("Origin")
		mobileHeader := c.Request.Header.Get("zmobile")
		ctxValue["origin"] = origin
		ctxValue["mobile"] = mobileHeader
		ctxValue["role"] = c.Request.Header.Get("role")
		request := c.Request
		requestWithValue := request.WithContext(context.WithValue(request.Context(), "zclaims", ctxValue))
		h.ServeHTTP(c.Writer, requestWithValue)
	}
}

func HealthCheckHandler(c *gin.Context) {
	log.Debugf("HealthCheckHandler Method --> %s", c.Request.Method)

	switch c.Request.Method {
	case http.MethodGet:
		GetHealthStatus(c.Writer)
	default:
		err := errors.New("method not supported")
		ResponseError(c.Writer, http.StatusBadRequest, err)
	}
}

// GetHealthStatus ...
func GetHealthStatus(w http.ResponseWriter) {
	healthStatus := "Zicops viltz manager service is healthy"
	response, _ := json.Marshal(healthStatus)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(response); err != nil {
		log.Errorf("GetHealthStatus ... unable to write JSON response: %v", err)
	}
}

// ResponseError ... essentially a single point of sending some error to route back
func ResponseError(w http.ResponseWriter, httpStatusCode int, err error) {
	log.Errorf("Response error %s", err.Error())
	response, _ := json.Marshal(err)
	w.Header().Add("Status", strconv.Itoa(httpStatusCode)+" "+err.Error())
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(httpStatusCode)

	if _, err := w.Write(response); err != nil {
		log.Errorf("ResponseError ... unable to write JSON response: %v", err)
	}
}
