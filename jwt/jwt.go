package jwt

import (
	"context"
	"errors"
	jwtgo "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var SigningKey = []byte("123456")

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwtgo.StandardClaims
}

func keyLookupFunction(token *jwtgo.Token) (interface{}, error) {
	return SigningKey, nil
}

func ParseJwt(tokenStr string) (*jwtgo.Token, *CustomClaims, error) {
	token, err := jwtgo.ParseWithClaims(tokenStr, &CustomClaims{}, keyLookupFunction)
	if err != nil {
		return nil, nil, err
	}
	if token == nil {
		return nil, nil, errors.New("Unable to parse token")
	}
	if token.Claims == nil {
		return nil, nil, errors.New("Unable to parse token claims")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		panic("Type Assertion failed")
	}
	return token, claims, err
}

func ExtractRoleFromToken(ctx context.Context) (string, error) {
	tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")

	if err != nil {
		return "", grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	token, claims, err := ParseJwt(tokenStr)
	if err != nil || token == nil {
		return "", grpc.Errorf(codes.Unauthenticated, err.Error())
	} else if !token.Valid {
		return "", grpc.Errorf(codes.Unauthenticated, "Invalid Token")
	}

	return claims.Role, nil
}

func ExtractUsernameFromToken(ctx context.Context) (string, error) {
	tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")

	if err != nil {
		return "", grpc.Errorf(codes.Unauthenticated, err.Error())
	}

	token, claims, err := ParseJwt(tokenStr)
	if err != nil || token == nil {
		return "", grpc.Errorf(codes.Unauthenticated, err.Error())
	} else if !token.Valid {
		return "", grpc.Errorf(codes.Unauthenticated, "Invalid Token")
	}

	return claims.Username, nil
}
