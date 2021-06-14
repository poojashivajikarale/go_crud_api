package router

import (
    "go_crud_api/middleware"

    "github.com/gorilla/mux"
    
    "github.com/newrelic/go-agent/v3/newrelic"
)

// Router is exported and used in main.go
func Router(app *newrelic.Application) *mux.Router {

    router := mux.NewRouter()

    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/customer/{id}", middleware.GetCustomer)).Methods("GET", "OPTIONS")
    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/customer", middleware.GetAllCustomer)).Methods("GET", "OPTIONS")
    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/newcustomer", middleware.CreateCustomer)).Methods("POST", "OPTIONS")
    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/customer/{id}", middleware.UpdateCustomer)).Methods("PUT", "OPTIONS")
    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/deletecustomer/{id}", middleware.DeleteCustomer)).Methods("DELETE", "OPTIONS")
    router.HandleFunc(newrelic.WrapHandleFunc(app, "/api/searchcustomer", middleware.SearchCustomer)).Methods("POST", "OPTIONS")

    return router
}
