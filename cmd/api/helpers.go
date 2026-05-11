package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
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

type idTokenClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"cognito:username"`
}

func extractIDTokenClaims(idToken string) (idTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return idTokenClaims{}, fmt.Errorf("invalid jwt format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return idTokenClaims{}, fmt.Errorf("decode payload:%w", err)
	}

	var claims idTokenClaims

	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return idTokenClaims{}, fmt.Errorf("unmarshal error:%w", err)

	}

	return claims, nil
}
