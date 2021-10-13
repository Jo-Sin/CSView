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
	"strings"
	"time"
	"math"

	"github.com/Jo-Sin/CSView/models"
)

var mcount = 0

type MongoController struct {
	session *mgo.Session
}
func (mc MongoController) InitializeDatabase() {
	mc.session.DB("test-db").DropDatabase()
	
	f, err := os.Open("data/Test task - Orders.csv")
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

		layout := "2006-01-02 15:04:05"
		dtime := strings.Replace(record[1], "T", " ", 1)
		dtime = dtime[:len(dtime)-1]
		dt, _ := time.Parse(layout, dtime)

		o.Id = bson.NewObjectId()
		_ , _ = fmt.Sscan(record[0], &o.OrderId)
		o.CreatedAt = dt
		o.OrderName = record[2]
		o.CustomerId = record[3]
		
		mc.session.DB("test-db").C("Orders").Insert(o)
		mcount = mcount + 1
	}
	mcount = int(math.Ceil(float64(mcount)/5))


	f, err = os.Open("data/Test task - Mongo - customers.csv")
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



	f, err = os.Open("data/Test task - Mongo - customer_companies.csv")
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



func (mc MongoController) GetOrders(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	page := 1
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

	err := mc.session.DB("test-db").C("Orders").Find(nil).Limit(5).Skip(skipCount).All(&orders)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		return
	}

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


	oj, err := json.Marshal(finalOrders)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", oj)
}

func GetMongoController(session *mgo.Session) *MongoController {
	return &MongoController{session}
}


func GetCount(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	oj, err := json.Marshal(mcount)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", oj)
}