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

func CreateAuctionLotInAuctionService(lotInformation LotDetails, lotClient lotPb.LotServiceClient) (*lotPb.CreateLotResponse, error) {
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

type CompleteMultiPartUploadReq struct {
	Parts []EtagPart `json:"parts"`
}
type EtagPart struct {
	Etag string `json:"etag"`
	Part int32  `json:"part"`
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

func CompleteMultiPartUpload(multiPartAttemptId int64, client lotimage.LotImageServiceClient, parts CompleteMultiPartUploadReq) (*lotimage.CompleteMultiPartUploadResponse, error) {
	partEtags := make([]*lotimage.PartEtag, 0, len(parts.Parts))
	for _, p := range parts.Parts {
		partEtags = append(partEtags, &lotimage.PartEtag{
			Etag: &p.Etag,
			Part: &p.Part,
		})
	}
	response, err := client.CompleteMultiPartUpload(context.TODO(), &lotimage.CompleteMultiPartUploadRequest{
		MultipartAttemptId: &multiPartAttemptId,
		PartEtags:          partEtags,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
