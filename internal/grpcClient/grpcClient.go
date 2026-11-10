package grpcClient

import (
	"errors"
	"fmt"

	uploadFile "github.com/Auction-Application/be-gateway/rpc/gen/upload_file/v1"
	uploadRoom "github.com/Auction-Application/be-gateway/rpc/gen/upload_room/v1"
	userPb "github.com/Auction-Application/be-gateway/rpc/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// func CreategrpcClient(port int) (*grpc.ClientConn, error) {
// 	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	if err != nil {
// 		fmt.Printf("Cannot start grpc client:%v", err)
// 		return nil, err

// 	}

// 	return conn, nil
// }

type GrpcTargetServerName string

const (
	UserService     GrpcTargetServerName = "user"
	LotService      GrpcTargetServerName = "auctionLot"
	LotImageService GrpcTargetServerName = "lotImage"
)

type GrpcServerTargets struct {
	Targets map[GrpcTargetServerName]string
}

type Clients struct {
	UserService     userPb.UserServiceClient
	LotService      uploadRoom.UploadRoomServiceClient
	LotImageService uploadFile.FileUploadServiceClient
	conns           []*grpc.ClientConn
}

func Dial(targets map[GrpcTargetServerName]string) (*Clients, error) {
	c := &Clients{}
	dial := func(addr string) (*grpc.ClientConn, error) {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("dial %s: %w", addr, err)
		}
		c.conns = append(c.conns, conn)
		return conn, nil
	}

	lotConn, err := dial(targets[LotService])
	if err != nil {
		return nil, err
	}
	c.LotService = uploadRoom.NewUploadRoomServiceClient(lotConn)

	userConn, err := dial(targets[UserService])
	if err != nil {
		return nil, err
	}
	c.UserService = userPb.NewUserServiceClient(userConn)

	lotImageConn, err := dial(targets[LotImageService])
	if err != nil {
		return nil, err
	}

	c.LotImageService = uploadFile.NewFileUploadServiceClient(lotImageConn)

	return c, nil
}

func (c *Clients) Close() error {
	var errs []error
	for _, conn := range c.conns {
		errs = append(errs, conn.Close())
	}
	return errors.Join(errs...)
}
