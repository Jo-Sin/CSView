package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
) 


type (
	Order struct {
		Id				bson.ObjectId	`json:"id" bson:"_id"`
		OrderId			int				`json:"order-id" bson:"order-id"`
		CreatedAt		time.Time		`json:"order-date" bson:"order-date"`
		OrderName		string			`json:"order-name" bson:"order-name"`
		CustomerId		string			`json:"cust-id" bson:"cust-id"`
	}

	Customer struct {
		Id				bson.ObjectId	`json:"id" bson:"_id"`
		UserId			string			`json:"user-id" bson:"user-id"`
		Login			string			`json:"login" bson:"login"`
		Password		string			`json:"password" bson:"password"`
		Name			string			`json:"name" bson:"name"`
		CompanyId		int				`json:"com-id" bson:"com-id"`
		CreditCards		string		`json:"cred-cards" bson:"cred-cards"`
	}

	Company struct {
		Id				bson.ObjectId	`json:"id" bson:"_id"`
		CompanyId		int				`json:"com-id" bson:"com-id"`
		CompanyName		string			`json:"com-name" bson:"com-name"`
	}

	MongoOrder struct {
		OrderId			int				`json:"order-id"`
		OrderName		string			`json:"order-name"`
		CompanyName		string			`json:"com-name"`
		CustomerName	string			`json:"cust-name"`
		OrderDate		time.Time		`json:"order-date"`
	}

	InitData struct {
		PageCount	int			`json:"page-count"`
		MinDate		time.Time	`json:"min-date"`
		MaxDate		time.Time	`json:"max-date"`
	}
)