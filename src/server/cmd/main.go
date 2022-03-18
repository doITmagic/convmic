package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"path"
	"runtime"
	"time"

	profile "github.com/pkg/profile"

	"github.com/doitmagic/convmic/pb"
	"github.com/doitmagic/convmic/src/server/internal"
	"github.com/doitmagic/convmic/src/server/model"
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

var provider = providers.NewCoingeckoProvider(context.Background())

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

	//populate first data
	provider.SyncCurrencies(4)

	//start Ticker to Load all currencies from provider and set the price
	//on provided seconds interval
	go syncAllCurrencies(provider, 10)

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

	//get all currecies from requested page
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

	fromCurrenciesConvert := []model.CurrencyConvert{}
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	fromCurrencies := req.GetFrom()
	toCurrency := req.GetTo()
	for _, currencyPb := range fromCurrencies {
		fromCurrenciesConvert = append(fromCurrenciesConvert, model.CurrencyConvert{Name: currencyPb.GetCurrencyName(), Amount: float64(currencyPb.GetCurrencyQty())})
		log.Infof("receive request to convert %v of %s to %s", currencyPb.GetCurrencyQty(), currencyPb.GetCurrencyName(), toCurrency)
	}

	convertedResp, err := convertCurrencies(provider, ctx, fromCurrenciesConvert, toCurrency)
	if err != nil {
		return &pb.GetCurrenciesConvertResponse{}, err
	}

	return convertedResp, nil
}

//syncAllCurrencies private, start ticker to sync currencies by seconds
func syncAllCurrencies(provider service.Provider, secondsInterval int) {
	d := time.NewTicker(time.Duration(secondsInterval) * time.Second)
	for tm := range d.C {
		log.Infof("START to sync currencies value from API provider at %v", tm.Local().GoString())
		provider.SyncCurrencies(2)
	}

}

//convertCurrencies private, dealing with the convert process
func convertCurrencies(provider service.Provider, ctx context.Context, from []model.CurrencyConvert, to string) (*pb.GetCurrenciesConvertResponse, error) {
	converted, err := provider.Convert(ctx, from, to)
	if err != nil {
		return nil, err
	}
	convertedResp := pb.GetCurrenciesConvertResponse{}

	for _, conv := range converted {
		convertedResp.Converted = append(convertedResp.Converted, &pb.CurrencyConvertResponse{From: conv.From, FromAmount: float32(conv.FromAmount), To: conv.To, ToAmount: float32(conv.ToAmount)})
	}

	return &convertedResp, nil
}

//appPath get the path of application
func appPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "./../")
}
