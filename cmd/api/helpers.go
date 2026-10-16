package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/Auction-Application/be-gateway/internal/grpcClient"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func extractIDTokenClaims(idToken string) (grpcClient.IdTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return grpcClient.IdTokenClaims{}, fmt.Errorf("invalid jwt format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return grpcClient.IdTokenClaims{}, fmt.Errorf("decode payload:%w", err)
	}

	var claims grpcClient.IdTokenClaims

	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return grpcClient.IdTokenClaims{}, fmt.Errorf("unmarshal error:%w", err)
	}

	return claims, nil
}

func writeProtoJSON(w http.ResponseWriter, status int, msg proto.Message) {
	b, err := protojson.MarshalOptions{
		UseProtoNames:     false, // false = camelCase (lotId); true = snake_case (lot_id)
		EmitUnpopulated:   false, // true if you want zero-value fields included
		EmitDefaultValues: true,
	}.Marshal(msg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
