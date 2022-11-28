package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `json:"username" binding:"required" required:"this is required"`
	Password string `json:"password"`
}

var (
	validator = g.Validator(login{})
)

func test(c *gin.Context) {
	req := &login{}

	// Parse json
	if err := c.BindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"message": validator.DecryptErrors(err),
		})
		return
	}

	c.JSON(200, gin.H{
		"good": 200,
	})
}

func TestDecryptErrors(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/", test)
	t.Run("username_lost", func(t *testing.T) {
		jsonValue, _ := json.Marshal(login{})
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)
		check(t, "{\"message\":{\"username\":\"this is required\"}}", string(responseData))
	})
	t.Run("wrong_data", func(t *testing.T) {
		jsonValue, _ := json.Marshal([]login{{Username: "dasewae"}})
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)
		check(t, "{\"message\":\"unmarshal error\"}", string(responseData))
	})
	t.Run("fine", func(t *testing.T) {
		jsonValue, _ := json.Marshal(login{Username: "dasewae"})
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		responseData, _ := ioutil.ReadAll(w.Body)
		check(t, "{\"good\":200}", string(responseData))
	})
}
