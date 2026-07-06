package grpcClient

import (
	"context"
	"fmt"

	lotPb "github.com/Auction-Application/be-gateway/rpc/gen/lot/v1"
)

type LotDetails struct {
	Title           string `json:"title"`
	Category        string `json:"category"`
	Description     string `json:"description"`
	BidOpeningPrice int32  `json:"bidOpeningPrice"`
}

// func CreateAuctionLotInAuctionService(lotInformation LotDetails, conn *grpc.ClientConn) (*lotPb.CreateLotResponse, error) {
func CreateAuctionLotInAuctionService(lotInformation LotDetails, lotClient lotPb.LotServiceClient) (*lotPb.CreateLotResponse, error) {
	// lotClient := lotPb.NewLotServiceClient(conn)
	createLotResponse, err := lotClient.CreateLot(context.TODO(), &lotPb.CreateLotRequest{
		Lot: &lotPb.Lot{
			LotTitle:        &lotInformation.Title,
			Description:     &lotInformation.Description,
			Category:        &lotInformation.Category,
			BidOpeningPrice: &lotInformation.BidOpeningPrice,
		},
	})

	fmt.Println(createLotResponse)

	if err != nil {
		return nil, err
	}

	return createLotResponse, nil

}
