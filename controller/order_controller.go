package controller

import (
	"net/http"

	"github.com/fadilahonespot/orderer/entity"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type OrderInterface interface {
	CreateOrder(c *gin.Context)
	GetOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

type defaultOrder struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) OrderInterface {
	return &defaultOrder{db}
}

func (s *defaultOrder) CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	dataOrder := entity.Order{
		CustomerName: req.CustomerName,
		OrderedAt:    req.OrderedAt,
	}

	for i := 0; i < len(req.Items); i++ {
		dataOrder.Items = append(dataOrder.Items, entity.Item{
			Code:        req.Items[i].ItemCode,
			Description: req.Items[i].Description,
			Quantity:    req.Items[i].Quantity,
		})
	}

	err := s.db.Create(&dataOrder).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, "failed create order")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := generateResponseSuccess(http.StatusOK, "success create order", req)
	c.JSON(http.StatusOK, response)
}

func (s *defaultOrder) GetOrder(c *gin.Context) {
	var orders []entity.Order
	err := s.db.Preload("Items").Find(&orders).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, "failed get order")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := generateResponseSuccess(http.StatusOK, "success create order", orders)
	c.JSON(http.StatusOK, response)
}

func (s *defaultOrder) UpdateOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := s.db.Take(&entity.Order{}, "id = ?", orderId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "order not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	tx := s.db.Begin()
	err = tx.Delete(&entity.Item{}, "order_id = ?", orderId).Error
	if err != nil {
		tx.Rollback()
		response := generateResponseError(http.StatusNotFound, "failed delete items")
		c.JSON(http.StatusNotFound, response)
		return
	}

	dataOrder := entity.Order{
		ID:           cast.ToInt(orderId),
		CustomerName: req.CustomerName,
		OrderedAt:    req.OrderedAt,
	}

	for i := 0; i < len(req.Items); i++ {
		dataOrder.Items = append(dataOrder.Items, entity.Item{
			Code:        req.Items[i].ItemCode,
			Description: req.Items[i].Description,
			Quantity:    req.Items[i].Quantity,
		})
	}

	err = tx.Save(&dataOrder).Error
	if err != nil {
		tx.Rollback()
		response := generateResponseError(http.StatusInternalServerError, "failed update order")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	tx.Commit()
	response := generateResponseSuccess(http.StatusOK, "success create order", req)
	c.JSON(http.StatusOK, response)
}

func (s *defaultOrder) DeleteOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	err := s.db.Take(&entity.Order{}, "id = ?", orderId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "order not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	tx := s.db.Begin()
	err = tx.Delete(&entity.Item{}, "order_id = ?", orderId).Error
	if err != nil {
		tx.Rollback()
		response := generateResponseError(http.StatusNotFound, "failed delete items")
		c.JSON(http.StatusNotFound, response)
		return
	}

	err = tx.Delete(&entity.Order{}, "id = ?", orderId).Error
	if err != nil {
		tx.Rollback()
		response := generateResponseError(http.StatusNotFound, "failed delete items")
		c.JSON(http.StatusNotFound, response)
		return
	}

	tx.Commit()
	response := generateResponseSuccess(http.StatusOK, "success delete order", nil)
	c.JSON(http.StatusOK, response)
}
