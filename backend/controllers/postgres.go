package controllers

import(
	"fmt"
	"io"
	"log"
	"os"
	"encoding/csv"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/Jo-Sin/CSView/backend/models"
)

var db *gorm.DB
var err error
var postgresFiles = [3]string {
	"backend/data/Test task - Postgres - orders.csv",
	"backend/data/Test task - Postgres - order_items.csv",
	"backend/data/Test task - Postgres - deliveries.csv",
}

func CloseDB() {
	db.Close()
}




//Connect to database and create relevant tables
//
func GetPostgresSession(dbtype string, user string, password string, dbname string, host string, port string) {

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",host,user,dbname,password,port)

	db, err = gorm.Open(dbtype, dbURI)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully connected")
	}

	db.AutoMigrate(&models.POrder{})
	db.AutoMigrate(&models.Item{})
	db.AutoMigrate(&models.Delivery{})


	InitializeDataPostgres()
}




//Populate tables with data from CSV files
//
func InitializeDataPostgres() {
	//Clear existing data from tables
	db.Exec("DELETE FROM p_orders")
	db.Exec("DELETE FROM items")
	db.Exec("DELETE FROM deliveries")


	f, err := os.Open(postgresFiles[0])
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	_, _ = r.Read()
	for {
		o := models.POrder{}

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


		_ , _ = fmt.Sscan(record[0], &o.OrderId)
		o.OrderDate = dt
		o.OrderName = record[2]
		o.CustomerId = record[3]

		db.Create(&o)
	}


	f, err = os.Open(postgresFiles[1])
	if err != nil {
		log.Fatal(err)
	}

	r = csv.NewReader(f)
	_, _ = r.Read()
	for {
		i := models.Item{}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		_ , _ = fmt.Sscan(record[0], &i.ItemId)
		_ , _ = fmt.Sscan(record[1], &i.OrderId)
		i.PricePerUnit, _ = strconv.ParseFloat(record[2], 64)
		_ , _ = fmt.Sscan(record[3], &i.Quantity)
		i.Product = record[4]

		db.Create(&i)
	}


	f, err = os.Open(postgresFiles[2])
	if err != nil {
		log.Fatal(err)
	}

	r = csv.NewReader(f)
	_, _ = r.Read()
	for {
		d := models.Delivery{}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		_ , _ = fmt.Sscan(record[1], &d.ItemId)
		_ , _ = fmt.Sscan(record[0], &d.DeliverId)
		_ , _ = fmt.Sscan(record[2], &d.DeliverQuantity)

		db.Create(&d)
	}
}




// On API request, this function returns 5 orders in JSON
//
func GetPostOrders(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Validates page number, defaults to 1
	page := 1
	_ , _ = fmt.Sscan(p.ByName("page"), &page)
	if page < 1 {
		page = 1
	}
	if page > mcount {
		page = mcount
	}
	skipCount := (page - 1) * 5


	var orders []models.POrder
	var items []models.Item
	var dels []models.Delivery
	var finalOrders []models.PostOrder

	//Retrieve upto 5 orders according to page
	db.Limit(5).Offset(skipCount).Find(&orders)

	for _, order := range orders {
		p := models.PostOrder{}
		p.OrderId = order.OrderId
		p.DeliveredAmount = 0.0
		p.TotalAmount = 0.0

		//Get all items in given order
		db.Where(&models.Item{OrderId: p.OrderId}).Find(&items)


		//Calculate delivered and total amount for each item of order
		for _, item := range items {
			p.TotalAmount = p.TotalAmount + (item.PricePerUnit * float64(item.Quantity))
			delQuant := 0

			db.Where(&models.Delivery{ItemId: item.ItemId}).Find(&dels)
			for _, delivery := range dels {
				delQuant = delQuant + delivery.DeliverQuantity
			}

			p.DeliveredAmount = p.DeliveredAmount + (item.PricePerUnit * float64(delQuant))
		}

		//Append relevant data to list
		finalOrders = append(finalOrders, p)
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