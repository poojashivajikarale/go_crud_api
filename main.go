package main

import (
    "fmt"
    "go_crud_api/router"
    "log"
    "net/http"
    "flag"
    "time"
    "os"
    "os/signal"
    "context"
    "github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
    app, err := newrelic.NewApplication(
       newrelic.ConfigAppName("Customer crud api application"),
       newrelic.ConfigLicense("25a2938c64f1b4d53ea1b14cd2e3a2ac6fbfNRAL"),
       newrelic.ConfigDistributedTracerEnabled(true),
    )
    if err != nil{
       fmt.Println("got error", err)
    }
    var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

    r := router.Router(app)
    fmt.Println("Starting server on the port 8080...")

    srv := &http.Server{
        Addr:         "0.0.0.0:8080",
        WriteTimeout: time.Second * 15,
        ReadTimeout:  time.Second * 15,
        IdleTimeout:  time.Second * 60,
        Handler: r, // Pass our instance of gorilla/mux in.
    }

    log.Fatal(srv.ListenAndServe())
    
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    srv.Shutdown(ctx)
    //log.Println("shutting down")
    os.Exit(0)
}
