# CSView

## PREREQUISITES
* MongoDB v5.0.3
* GoLang v1.17.2
* PostgreSQL v14



## STEPS

1. Start MongoDB Server

2. Start PostgreSQL Server

3. If necessary, change the connection string in the function **getMongoSession()** and the connection parameters in the function **getPostgresSession()** in **main.go**

2. Run Go Server with `go run main.go`