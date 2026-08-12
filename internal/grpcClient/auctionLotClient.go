package grpcClient

import (
	"context"
	"fmt"

	lotPb "github.com/Auction-Application/be-gateway/rpc/gen/lot/v1"
	lotimage "github.com/Auction-Application/be-gateway/rpc/gen/lot_image/v1"
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

type LotImage struct {
	Sha256       string `json:"sha256"`
	FileName     string `json:"fileName"`
	ClientFileId int32  `json:"clientFileId"`
	FileSize     uint64 `json:"fileSize"`
}

type UploadFiles struct {
	LotId     string     `json:"lotId"`
	LotImages []LotImage `json:"lotImages"`
}

func GenerateUploadPresignedUrls(uploadFiles UploadFiles, client lotimage.LotImageServiceClient) (*lotimage.GeneratePresignedUrlResponse, error) {
	lotImages := make([]*lotimage.ImageFile, 0, len(uploadFiles.LotImages))
	for _, v := range uploadFiles.LotImages {
		lotImages = append(lotImages, &lotimage.ImageFile{
			Sha256:       &v.Sha256,
			FileName:     &v.FileName,
			ClientFileId: &v.ClientFileId,
			FileSize:     &v.FileSize,
		})
	}
	presignedUrls, err := client.GeneratePresignedUrl(context.TODO(), &lotimage.GeneratePresignedUrlRequest{
		LotId:     &uploadFiles.LotId,
		LotImages: lotImages,
	})
	if err != nil {
		return nil, err
	}
	return presignedUrls, nil
}
