package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"path"
	"runtime"

	profile "github.com/pkg/profile"

	"github.com/doitmagic/convmic/pb"
	"github.com/doitmagic/convmic/src/server/helper"
	"github.com/doitmagic/convmic/src/server/internal"
	"github.com/doitmagic/convmic/src/server/providers"
	"github.com/doitmagic/convmic/src/server/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port = "9000"
)

type server struct {
	pb.UnimplementedConvertServiceServer
}

//start server
func main() {

	//profiling part if it use option " -cpuprofile true"
	var cpuprofile bool
	flag.BoolVar(&cpuprofile, "cpuprofile", false, "write cpu profile to filen")

	flag.Parse()
	if cpuprofile {
		defer profile.Start(profile.ProfilePath(appPath())).Stop()
		log.Printf("Start CPU profiling")
	}

	//populate dummy data
	helper.PopulateData(internal.GetInstance())

	//set the provider, it wil be loaded from config
	provider := providers.NewCoingeckoProvider(context.Background())

	//Load all currencies from provider and set the price
	err := syncAllCurrencies(provider)
	if err != nil {
		log.Fatal(err)
	}

	g := grpc.NewServer()

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	pb.RegisterConvertServiceServer(g, &server{})
	log.Printf("Star server on port %s", port)

	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc: " + err.Error())
	}
}

//List all currency from memory whith pagination
func (s *server) List(ctx context.Context, req *pb.GetListCurrenciesRequest) (*pb.GetListCurrenciesResponse, error) {
	var perPage, page int32

	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	page = req.GetPageNumber()
	perPage = req.GetResultPerPage()

	//add default per page if is 0
	if perPage == 0 {
		perPage = 10
	}

	currencies := []*pb.Currency{}

	//get onlu currecies from requested page
	paginatedCurrencies := internal.GetInstance().GetCurrenciesByPage(page, perPage)

	if len(paginatedCurrencies) > 0 {
		for name, val := range paginatedCurrencies {
			currencies = append(currencies, &pb.Currency{Name: name, Value: val})
		}
	}

	return &pb.GetListCurrenciesResponse{PageNumber: page, Currencies: currencies}, nil
}

//Convert all received currency name and qty to another currency
//if are more will make a batch convert
func (s *server) Convert(ctx context.Context, req *pb.GetCurrenciesConvertRequest) (*pb.GetCurrenciesConvertResponse, error) {

	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	fromCurrencies := req.GetFrom()
	toCurrency := req.GetTo()
	for _, currencyPb := range fromCurrencies {
		fmt.Println("convert: ", currencyPb.GetCurrencyName(), " to ", toCurrency)
	}

	return &pb.GetCurrenciesConvertResponse{}, nil
}

func syncAllCurrencies(provider service.Provider) error {
	//to be be implemented
	go provider.SyncCurrencies(0, 2)
	return nil
}

//appPath get the path of application
func appPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "./../")
}
