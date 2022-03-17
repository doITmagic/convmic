package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/doitmagic/convmic/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = ":9000"

func main() {

	var command string
	var page, perpage int

	flag.StringVar(&command, "c", "list", "Command to run")
	flag.IntVar(&page, "p", 1, "page number")
	flag.IntVar(&perpage, "pp", 1, "records per page")
	flag.Parse()

	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewConvertServiceClient(conn)
	switch command {
	case "convert":
		Convert(client)
	case "list":
		if page == 0 {
			page = 1
		}
		if perpage == 0 {
			perpage = 10
		}
		List(client, page, perpage)
	}
}

//Convert client request to convert *hardcoded
func Convert(client pb.ConvertServiceClient) {
	ctx := context.Background()
	currencies := []*pb.CurrencyConvert{}
	currencies = append(currencies, &pb.CurrencyConvert{CurrencyName: "Atest1", CurrencyQty: 2})
	currencies = append(currencies, &pb.CurrencyConvert{CurrencyName: "Atest1", CurrencyQty: 3})
	resp, err := client.Convert(ctx, &pb.GetCurrenciesConvertRequest{From: currencies, To: "Aave"})
	fmt.Println(resp, err)
}

//List Client request to list currencies with paginations
func List(client pb.ConvertServiceClient, page, PerPageOption int) {
	ctx := context.Background()
	fmt.Println(int32(page), int32(PerPageOption))
	resp, err := client.List(ctx, &pb.GetListCurrenciesRequest{PageNumber: int32(page), ResultPerPage: int32(PerPageOption)})
	fmt.Println(resp, err)
}
