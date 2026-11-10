package grpcClient

import (
	"context"
	"fmt"
	"strconv"

	uploadFile "github.com/Auction-Application/be-gateway/rpc/gen/upload_file/v1"
	uploadRoom "github.com/Auction-Application/be-gateway/rpc/gen/upload_room/v1"
)

type UploadRoomDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateAuctionLotInAuctionService(lotInformation UploadRoomDetails, lotClient uploadRoom.UploadRoomServiceClient) (*uploadRoom.CreateRoomResponse, error) {
	createLotResponse, err := lotClient.CreateRoom(context.TODO(), &uploadRoom.CreateRoomRequest{
		Room: &uploadRoom.UploadRoom{
			Title:       &lotInformation.Title,
			Description: &lotInformation.Description,
		},
	})

	fmt.Println(createLotResponse)

	if err != nil {
		return nil, err
	}

	return createLotResponse, nil
}

type File struct {
	Sha256       string `json:"sha256"`
	FileName     string `json:"fileName"`
	ClientFileId int32  `json:"clientFileId"`
	FileSize     uint64 `json:"fileSize"`
}

type UploadFiles struct {
	UploadRoomId string `json:"uploadRoomId"`
	Files        []File `json:"files"`
}

type CompleteMultiPartUploadReq struct {
	Parts []EtagPart `json:"parts"`
}
type EtagPart struct {
	Etag string `json:"etag"`
	Part uint32 `json:"part"`
}

func GenerateUploadPresignedUrls(uploadFiles UploadFiles, client uploadFile.FileUploadServiceClient) (*uploadFile.GeneratePresignedUrlResponse, error) {
	lotImages := make([]*uploadFile.File, 0, len(uploadFiles.Files))
	for _, v := range uploadFiles.Files {
		lotImages = append(lotImages, &uploadFile.File{
			Sha256:       &v.Sha256,
			ClientFileId: &v.ClientFileId,
			// fileSize is string in json and uint64 in protobuf
			FileSize: &v.FileSize,
		})
	}
	// not using Math/big package here, to convert string to bigint
	roomId, err := strconv.Atoi(uploadFiles.UploadRoomId)
	if err != nil {
		return nil, err
	}
	presignedUrls, err := client.GeneratePresignedUrl(context.TODO(), &uploadFile.GeneratePresignedUrlRequest{
		// not using Math/big package here, to convert string to bigint
		UploadRoomId: new(uint64(roomId)),
		Files:        lotImages,
	})
	if err != nil {
		return nil, err
	}
	return presignedUrls, nil
}

func CompleteMultiPartUpload(multiPartAttemptId int64, client uploadFile.FileUploadServiceClient, parts CompleteMultiPartUploadReq) (*uploadFile.CompleteMultiPartUploadResponse, error) {
	partEtags := make([]*uploadFile.PartEtag, 0, len(parts.Parts))
	for _, p := range parts.Parts {
		partEtags = append(partEtags, &uploadFile.PartEtag{
			Etag: &p.Etag,
			Part: &p.Part,
		})
	}
	response, err := client.CompleteMultiPartUpload(context.TODO(), &uploadFile.CompleteMultiPartUploadRequest{
		MultipartAttemptId: &multiPartAttemptId,
		PartEtags:          partEtags,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
