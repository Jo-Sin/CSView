package controllers

import(
	"fmt"
	"io"
	"log"
	"os"
	"encoding/csv"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	_ "strings"
	"time"
	"math"
	"github.com/Jo-Sin/CSView/backend/models"
)

var mcount = 0
var mongoFiles = [3]string {
	"backend/data/Test task - Orders.csv",
	"backend/data/Test task - Mongo - customers.csv",
	"backend/data/Test task - Mongo - customer_companies.csv",
}

// Set unrealistic initial dates that will be overriden by data
var minDate, _ = time.Parse("2006-01-02", "4000-12-12")
var maxDate, _ = time.Parse("2006-01-02", "1000-01-01")




type MongoController struct {
	session *mgo.Session
}

func GetMongoController(session *mgo.Session) *MongoController {
	return &MongoController{session}
}




// Populate the database with data from the CSV files
//
func (mc MongoController) InitializeDatabase() {
	// Clear any existing data in the database
	mc.session.DB("test-db").DropDatabase()


	f, err := os.Open(mongoFiles[0])
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	_, _ = r.Read()
	for {
		o := models.Order{}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		layout := "2006-01-02T15:04:05Z"
		dt, _ := time.Parse(layout, record[1])

		o.Id = bson.NewObjectId()
		_ , _ = fmt.Sscan(record[0], &o.OrderId)
		o.CreatedAt = dt
		o.OrderName = record[2]
		o.CustomerId = record[3]
		
		mc.session.DB("test-db").C("Orders").Insert(o)

		// Count the number of orders in the database
		mcount = mcount + 1
		// Check and set date range
		if minDate.After(dt) {
			minDate = dt
		}
		if maxDate.Before(dt) {
			maxDate = dt
		}
	}


	//Divide order count by 5 to get page count
	mcount = int(math.Ceil(float64(mcount)/5))


	f, err = os.Open(mongoFiles[1])
	if err != nil {
		log.Fatal(err)
	}

	r = csv.NewReader(f)
	_, _ = r.Read()
	for {
		c := models.Customer{}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		c.Id = bson.NewObjectId()
		c.UserId = record[0]
		c.Login = record[1]
		c.Password = record[2]
		c.Name = record[3]
		_ , _ = fmt.Sscan(record[4], &c.CompanyId)
		c.CreditCards = record[5]
		
		mc.session.DB("test-db").C("Customers").Insert(c)
	}



	f, err = os.Open(mongoFiles[2])
	if err != nil {
		log.Fatal(err)
	}

	r = csv.NewReader(f)
	_, _ = r.Read()
	for {
		c := models.Company{}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		c.Id = bson.NewObjectId()
		_ , _ = fmt.Sscan(record[0], &c.CompanyId)
		c.CompanyName = record[1]
		
		mc.session.DB("test-db").C("Companies").Insert(c)
	}

}




// On API request, this function returns 5 orders in JSON
//
func (mc MongoController) GetOrders(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//Get params and format for the query
	page := 1
	src := p.ByName("src")[1:]
	fromdate, _ := time.Parse("2006-01-02",p.ByName("lodate"))
	todate, _ := time.Parse("2006-01-02",p.ByName("hidate"))
	todate = todate.AddDate(0,0,1)
	fromdate = fromdate.Add(-time.Minute * 750)
	todate = todate.Add(-time.Minute * 750)

	_ , _ = fmt.Sscan(p.ByName("page"), &page)
	if page < 1 {
		page = 1
	}
	if page > mcount {
		page = mcount
	}
	skipCount := (page - 1) * 5

	var orders []models.Order
	var customer models.Customer
	var company models.Company
	var finalOrders []models.MongoOrder


	queryVar := mc.session.DB("test-db").C("Orders").Find(bson.M{
		"order-name": bson.RegEx{Pattern: src},
		"order-date": bson.M{"$gt":fromdate, "$lt":todate},
		})

	//Update page count according to query params
	mcount, _ = queryVar.Count()
	mcount = int(math.Ceil(float64(mcount)/5))
	
	//Retrieve upto 5 orders according to query params
	err := queryVar.Limit(5).Skip(skipCount).All(&orders)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		return
	}

	//Retrieve associated data from other collections and append relevant data to list
	for _, order := range orders {
		m := models.MongoOrder{}
		m.OrderId = order.OrderId
		m.OrderName = order.OrderName
		m.OrderDate = order.CreatedAt

		err := mc.session.DB("test-db").C("Customers").Find(bson.M{"user-id": order.CustomerId}).One(&customer)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println(err)
			return
		}

		m.CustomerName = customer.Name

		err = mc.session.DB("test-db").C("Companies").Find(bson.M{"com-id": customer.CompanyId}).One(&company)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println(err)
			return
		}

		m.CompanyName = company.CompanyName

		finalOrders = append(finalOrders, m)
	}


	// Convert final data to JSON and return
	oj, err := json.Marshal(finalOrders)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", oj)
}




// Sends page count and dates for initialization and updates to frontend
//
func GetCount(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	initInfo := models.InitData{}
	initInfo.PageCount = mcount
	initInfo.MinDate = minDate
	initInfo.MaxDate = maxDate

	oj, err := json.Marshal(initInfo)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", oj)
}