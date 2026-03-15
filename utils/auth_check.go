package utils

import (
	"errors"
	"net/http"
	"strings"

	firebase_shared "github.com/C2NOfficial/C2NGCShared/firebase"
)

func IsAuthorizedAndAdmin(request *http.Request) (string, error) {
	ctx := request.Context()

	authHeader := request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	idToken := parts[1]
	token, err := firebase_shared.AuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}

	claimValue, ok := token.Claims["admin"].(bool)
	if !ok {
		return "", errors.New("invalid admin claim")
	}
	if !claimValue {
		return "", errors.New("user is not an admin")
	}
	return token.UID, nil
}