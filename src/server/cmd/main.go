package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"path"
	"runtime"
	"strconv"
	"unicode"

	"github.com/doitmagic/convmic/pb"
	"github.com/doitmagic/convmic/src/server/internal"
	"github.com/doitmagic/convmic/src/server/providers"
	"github.com/doitmagic/convmic/src/server/service"
	"google.golang.org/grpc"
)

const (
	port = "9000"
)

type server struct {
	pb.UnimplementedConvertServiceServer
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

//start server
func main() {

	populateData()

	provider := providers.NewCoingeckoProvider()
	SyncAllCurrencies(provider)

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

func SyncAllCurrencies(provider service.Provider) {
	provider.GetCurrencies(context.Background())
}

//appPath get the path of application
func appPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "./../")
}

//populateData populate data until using api request
func populateData() {
	for r := 'a'; r < 'g'; r++ {
		R := unicode.ToUpper(r)
		for j := 0; j < 11; j++ {
			internal.GetInstance().SetCurrency(fmt.Sprintf("%c", R)+"test"+strconv.Itoa(rand.Intn(100)), 100)
		}
	}
}
