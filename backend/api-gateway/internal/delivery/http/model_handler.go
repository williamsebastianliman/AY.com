package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type PredictRequest struct {
  Text string `json:"text"`
}

type PredictResponse struct {
  Prediction int `json:"prediction"`
}

func PredictHandler() gin.HandlerFunc {
  mlAddr := os.Getenv("ML_SVC_ADDR")
  if mlAddr == "" {
    mlAddr = "localhost:5000"
  }
  url := "http://" + mlAddr + "/predict"

  return func(c *gin.Context) {
    var req PredictRequest
    if err := c.ShouldBindJSON(&req); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }

    body, _ := json.Marshal(req)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
      c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
      return
    }
    defer resp.Body.Close()

    var pr PredictResponse
    if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
      c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, pr)
  }
}