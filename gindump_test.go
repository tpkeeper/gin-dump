package gindump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, path string,body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestMIMJSON(t *testing.T) {
	router := gin.New()
	router.Use(Dump())

	router.POST("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":true,
			"data":"gin-dump",
		})
	})


	type params struct {
		StartTime string `json:"start_time"`
		EndTime string `json:"end_time"`
	}

	var httpdata = params{
		StartTime:"2019-03-03",
		EndTime:"2019-03-03",
	}
	b ,err := json.Marshal(httpdata)
	if err != nil {
		fmt.Println("json format error:", err)
		return
	}

	body := bytes.NewBuffer(b)
	performRequest(router, "POST", "/dump",body)

}
