package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RSOam/chargers-service/chargers"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	consulapi "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	httpAddr := flag.String("http", ":8080", "http listen addr")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "chargers",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://amAdmin:"+os.Getenv("DBpw")+"@cluster0.4lsbm.mongodb.net/Chargers?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}

	level.Info(logger).Log("msg", "DB Connected")
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	//CONSUL
	consulClient, err := getConsulClient(os.Getenv("CONSUL_ADDR"))
	if err != nil {
		panic(err)
	}
	level.Info(logger).Log("msg", "Consul Connected")

	//
	collection := client.Database("Chargers")

	flag.Parse()
	ctx2 := context.Background()
	var srv chargers.ChargersService
	{
		database := chargers.NewDatabase(collection, logger, *consulClient)
		srv = chargers.NewService(database, logger, *consulClient)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := chargers.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := chargers.NewHttpServer(ctx2, endpoints, collection, *consulClient)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}
func test() string {
	return "Ok"
}
func getConsulClient(address string) (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	consul, err := consulapi.NewClient(config)
	return consul, err

}
