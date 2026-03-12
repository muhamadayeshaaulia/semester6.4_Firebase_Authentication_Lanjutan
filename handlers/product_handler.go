package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/services"

)

type ProductHandler struct {
	productService *services.ProductService
}
