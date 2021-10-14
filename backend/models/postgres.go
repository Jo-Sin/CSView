package models

import (
	"time"
	"github.com/jinzhu/gorm"
) 


type (
	POrder struct {
		gorm.Model
		OrderId			int				`json:"order-id"`
		OrderDate		time.Time		`json:"order-date"`
		OrderName		string			`json:"order-name"`
		CustomerId		string			`json:"cust-id"`
	}

	Item struct {
		gorm.Model
		ItemId			int				`json:"item-id"`
		OrderId			int				`json:"order-id"`
		PricePerUnit	float64			`json:"ppu"`
		Quantity		int				`json:"quantity"`
		Product			string			`json:"item-name"`
	}

	Delivery struct {
		gorm.Model
		DeliverId		int				`json:"deliver-id"`
		ItemId			int				`json:"item-id"`
		DeliverQuantity	int				`json:"deliver-quantity"`
	}

	PostOrder struct {
		OrderId			int				`json:"order-id"`
		DeliveredAmount	float64			`json:"delivered-amount"`
		TotalAmount		float64			`json:"total-amount"`
	}
)