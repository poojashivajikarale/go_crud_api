package middleware

import (
    "database/sql"
    "encoding/json" // package to encode and decode the json into struct and vice versa
    "fmt"
    "bytes"
    "strings"
    "go_crud_api/models" // models package where customer schema is defined
    "log"
    "net/http" // used to access the request and response object of the api
    "os"       // used to read the environment variable
    "strconv"  // package used to covert string into int type

    "github.com/gorilla/mux" // used to get the params from the route

    "github.com/joho/godotenv" // package used to read the .env file
    _ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
    // load .env file
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Open the connection
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

    if err != nil {
        panic(err)
    }

    // check the connection
    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")
    // return the connection
    return db
}

// CreateCustomer create a customer in the postgres db
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
    // set the header to content type x-www-form-urlencoded
    // Allow all origin to handle cors issue
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // create an empty customer of type models.Customer
    var customer models.Customer

    // decode the json request to customer
    err := json.NewDecoder(r.Body).Decode(&customer)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    // call insert customer function and pass the customer
    insertID := insertCustomer(customer)

    // format a response object
    res := response{
        ID:      insertID,
        Message: "Customer created successfully",
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

// GetCustomer will return a single customer by its id
func GetCustomer(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // get the legal_entity_id from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id type from string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // call the getCustomer function with legal_entity_id to retrieve a single customer
    customer, err := getCustomer(int64(id))

    if err != nil {
        log.Fatalf("Unable to get customer. %v", err)
    }

    // send the response
    json.NewEncoder(w).Encode(customer)
}

// GetAllCustomer will return all the customers
func GetAllCustomer(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // get all the customers in the db
    customers, err := getAllCustomers()

    if err != nil {
        log.Fatalf("Unable to get all customers. %v", err)
    }

    // send all the customers as response
    json.NewEncoder(w).Encode(customers)
}

// UpdateCustomer update customer's detail in the postgres db
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // get the lega_entity_id from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id type from string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // create an empty customer of type models.Customer
    var customer models.Customer

    // decode the json request to customer
    err = json.NewDecoder(r.Body).Decode(&customer)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    // call update customer to update the customer
    updatedRows := updateCustomer(int64(id), customer)

    // format the message string
    msg := fmt.Sprintf("Customer updated successfully. Total rows/record affected %v", updatedRows)

    // format the response message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

// DeleteCustomer delete customer's detail in the postgres db
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // get the legal_entity_id from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id in string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // call the deleteCustomer, convert the int to int64
    deletedRows := deleteCustomer(int64(id))

    // format the message string
    msg := fmt.Sprintf("Customer deleted successfully. Total rows/record affected %v", deletedRows)

    // format the reponse message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

//SearchCustomer function searches customers using multiple field
func SearchCustomer(w http.ResponseWriter, r *http.Request){
    // set the header to content type x-www-form-urlencoded
    // Allow all origin to handle cors issue
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // call search customer function and pass the customer
    customers,err := searchCustomer(r)

    if err != nil {
        log.Fatalf("Unable to get all customers. %v", err)
    }

    // send all the customers as response
    json.NewEncoder(w).Encode(customers)
    
}

//------------------------- handler functions ----------------
// insert one customer in the DB
func insertCustomer(customer models.Customer) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the insert sql query
    // returning legal_entity_id will return the id of the inserted customer
    sqlStatement := `INSERT INTO customers (legal_entity_id, first_name, last_name, company_name, bankruptcy_indicator_flag, legal_entity_stage, legal_entity_type, date_of_birth, created_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING legal_entity_id`

    // the inserted id will store in this id
    var id int64

    // execute the sql statement
    // Scan function will save the insert id in the id
    err := db.QueryRow(sqlStatement, customer.LegalEntityId, customer.FirstName, customer.LastName, customer.CompanyName, customer.BankruptcyIndicatorFlag, customer.LegalEntityState, customer.LegalEntityType, customer.DateOfBirth, customer.CreatedDate).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    // return the inserted id
    return id
}

// get one customer from the DB by its legal_entity_id
func getCustomer(id int64) (models.Customer, error) {
    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create a customer of models.Customer type
    var customer models.Customer

    // create the select sql query
    sqlStatement := `SELECT * FROM customers WHERE legal_entity_id=$1`

    // execute the sql statement
    row := db.QueryRow(sqlStatement, id)

    // unmarshal the row object to customer
    err := row.Scan(&customer.LegalEntityId, &customer.FirstName, &customer.LastName, &customer.CompanyName, &customer.BankruptcyIndicatorFlag, &customer.LegalEntityState, &customer.LegalEntityType, &customer.DateOfBirth, &customer.CreatedDate)

    switch err {
    case sql.ErrNoRows:
        fmt.Println("No rows were returned!")
        return customer, nil
    case nil:
        return customer, nil
    default:
        log.Fatalf("Unable to scan the row. %v", err)
    }

    // return empty customer on error
    return customer, err
}

// get all customers
func getAllCustomers() ([]models.Customer, error) {
    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    var customers []models.Customer

    // create the select sql query
    sqlStatement := `SELECT * FROM customers`

    // execute the sql statement
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // close the statement
    defer rows.Close()

    // iterate over the rows
    for rows.Next() {
        var customer models.Customer

        // unmarshal the row object to customer
        err = rows.Scan(&customer.LegalEntityId, &customer.FirstName, &customer.LastName, &customer.CompanyName, &customer.BankruptcyIndicatorFlag, &customer.LegalEntityState, &customer.LegalEntityType, &customer.DateOfBirth, &customer.CreatedDate)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        // append the customer in the customers slice
        customers = append(customers, customer)

    }

    // return empty customers on error
    return customers, err
}

// update customer in the DB
func updateCustomer(id int64, customer models.Customer) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the update sql query
    sqlStatement := `UPDATE customers SET first_name=$2, last_name=$3, company_name=$4, bankruptcy_indicator_flag=$5, legal_entity_stage=$6, legal_entity_type=$7, created_date=$8, date_of_birth=$9 WHERE legal_entity_id=$1`

    // execute the sql statement
    res, err := db.Exec(sqlStatement, id, customer.FirstName, customer.LastName, customer.CompanyName, customer.BankruptcyIndicatorFlag, customer.LegalEntityState, customer.LegalEntityType, customer.CreatedDate, customer.DateOfBirth)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

// delete customer in the DB
func deleteCustomer(id int64) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the delete sql query
    sqlStatement := `DELETE FROM customers WHERE legal_entity_id=$1`

    // execute the sql statement
    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

//seach customers in the DB
func searchCustomer(r *http.Request)([]models.Customer, error){

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()
    var array_string []string

    for k,v := range r.URL.Query(){
       var buffer bytes.Buffer
       buffer.WriteString(k)
       buffer.WriteString(" = ")
       buffer.WriteString(strings.Join(v, ""))
       array_string = append(array_string, buffer.String())
    }

    search_query := strings.Join(array_string, " OR ")
    var customers []models.Customer

    // create the select sql query
    sqlStatement := `SELECT * FROM customers WHERE `+ search_query

    // execute the sql statement
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // close the statement
    defer rows.Close()

    // iterate over the rows
    for rows.Next() {
        var customer models.Customer

        // unmarshal the row object to customer
         err = rows.Scan(&customer.LegalEntityId, &customer.FirstName, &customer.LastName, &customer.CompanyName, &customer.BankruptcyIndicatorFlag, &customer.LegalEntityState, &customer.LegalEntityType, &customer.DateOfBirth, &customer.CreatedDate)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        // append the customer in the customers slice
        customers = append(customers, customer)

    }

    // return empty customers on error
    
    return customers, nil

}
