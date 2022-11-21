package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"

	"github.com/robbert229/jwt"
)

const (
	CryptoKey = "secret123456789"
)

var (
	ErrInvalidToken = errors.New("failed to decode the provided Token")
)

// JWTEncodeUserID creates JWT with userID encoded inside.
func JWTEncodeUserID(value interface{}) (string, error) {
	return JWTEncode("sub", value)
}

// JWTEncode creates JWT with encoded inside passed value.
func JWTEncode(key string, value interface{}) (string, error) {
	algorithm := jwt.HmacSha256(CryptoKey)

	claims := jwt.NewClaim()
	claims.Set(key, value)

	token, err := algorithm.Encode(claims)
	if err != nil {
		return ``, err
	}

	if err = algorithm.Validate(token); err != nil {
		return ``, err
	}

	return token, nil
}

// JWTDecodeUserID provides userID from JWT, if decoding is successful.
func JWTDecodeUserID(token string) (int, error) {
	value, err := JWTDecode(token, "sub")
	if err != nil {
		return -1, err
	}
	return int(value.(float64)), nil
}

// JWTDecode decodes the passed JWT and returns the interface value.
func JWTDecode(token, key string) (interface{}, error) {
	algorithm := jwt.HmacSha256(CryptoKey)

	if err := algorithm.Validate(token); err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}

	claims, err := algorithm.Decode(token)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}

	return claims.Get(key)
}

type ctxkey string

var (
	userID ctxkey = "userID"
)

// GetUserIDFromCTX returns from context userID if found.
func GetUserIDFromCTX(ctx context.Context) int {
	value, ok := ctx.Value(userID).(int)
	if !ok {
		return -1
	}
	return value
}

// SetUserIDToCTX add userID to the context.
func SetUserIDToCTX(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, userID, value)
}

// EncryptPass creates encrypted password based on md5 hash algorithm.
func EncryptPass(pass string) string {
	h := md5.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}
