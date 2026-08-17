package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Auction-Application/be-gateway/internal/grpcClient"
	"github.com/go-chi/chi/v5"
)

func (app *application) createAuctionLot(w http.ResponseWriter, r *http.Request) {
	var req grpcClient.LotDetails
	if err := decodeJSON(r, &req); err != nil {
		fmt.Println("bad req")
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	auctionLotRes, err := grpcClient.CreateAuctionLotInAuctionService(req, app.rpcClient.LotService)
	if err != nil {
		fmt.Println(err)
		fmt.Println("createAuctionResponseError")
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	writeProtoJSON(w, http.StatusOK, auctionLotRes)
	// writeJSON(w, http.StatusOK, auctionLotRes)

	// if req.Username == "" {
	// 	writeError(w, http.StatusBadRequest, "username is required")
	// 	return
	// }
	// if req.Password == "" {
	// 	writeError(w, http.StatusBadRequest, "password is required")
	// 	return
	// }

	// if err != nil {
	// 	fmt.Printf("Error getting signInOutput,%v", err)
	// 	writeError(w, http.StatusBadRequest, "error signing in")
	// 	return
	// }

	// if !ok {
	// 	fmt.Println("Error in conversion of jwt from any type to string type")
	// 	writeError(w, http.StatusInternalServerError, "Something went wrong")
	// 	return
	// }

	// if err != nil {
	// 	fmt.Printf("Error extracting IdToken,%w", err)
	// 	writeError(w, http.StatusInternalServerError, "Something went wrong")
	// 	return
	// }

	// // err = grpcClient.InsertUserInUserService(userClaim, app.rpcClient)

	// if err != nil {
	// 	fmt.Printf("error in inserting user,%v \n", err)
	// 	writeError(w, http.StatusInternalServerError, err.Error())
	// 	return

	// }

	// fmt.Printf("%+v", userClaim)

	// type hello struct {
	// 	Name string `json:"name"`
	// }
}

func (app *application) generatePresignedUrls(w http.ResponseWriter, r *http.Request) {
	var req grpcClient.UploadFiles

	if err := decodeJSON(r, &req); err != nil {
		fmt.Println("Cannot decode json")
		writeError(w, http.StatusBadRequest, "invalid json:"+err.Error())
		return
	}

	presignedUrlResponse, err := grpcClient.GenerateUploadPresignedUrls(req, app.rpcClient.LotImageService)
	if err != nil {
		fmt.Println(err)
		fmt.Println("generatePresignedUrls for lot image Error")
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	writeProtoJSON(w, http.StatusOK, presignedUrlResponse)
}

func (app *application) completeMultiPartUpload(w http.ResponseWriter, r *http.Request) {
	var req grpcClient.CompleteMultiPartUploadReq

	if err := decodeJSON(r, &req); err != nil {
		fmt.Println("Cannot decode json")
		writeError(w, http.StatusBadRequest, "invalid json:"+err.Error())
		return
	}
	attemptId, err := strconv.Atoi(chi.URLParam(r, "attemptId"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Attempt Id should be parsable to integer")
		return
	}

	response, err := grpcClient.CompleteMultiPartUpload(int64(attemptId), app.rpcClient.LotImageService, req)
	if err != nil {
		fmt.Println(err)
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	writeProtoJSON(w, http.StatusOK, response)
}
