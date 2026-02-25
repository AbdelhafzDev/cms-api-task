package crypto

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")
	ErrMissingToken      = errors.New("missing token")
	ErrMissingUserID     = errors.New("missing or invalid user identifier in token")
)

func GenerateToken(privateKey *rsa.PrivateKey, claims jwt.MapClaims, expiry time.Duration) (string, error) {
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(expiry).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func ValidateToken(publicKey *rsa.PublicKey, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrInvalidAuthFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}

// Supports both "sub" and "user_id" claim keys
func ExtractUserID(claims jwt.MapClaims) (string, bool) {
	if sub, ok := claims["sub"].(string); ok && sub != "" {
		return sub, true
	}

	if userID, ok := claims["user_id"].(string); ok && userID != "" {
		return userID, true
	}

	return "", false
}

func ExtractEmail(claims jwt.MapClaims) string {
	if email, ok := claims["email"].(string); ok {
		return email
	}
	return ""
}

// Supports both JSON array and legacy comma-separated string
func ExtractRoles(claims jwt.MapClaims) []string {
	switch v := claims["roles"].(type) {
	case []interface{}:
		roles := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				roles = append(roles, s)
			}
		}
		return roles
	}
	if role, ok := claims["role"].(string); ok && role != "" {
		return strings.Split(role, ",")
	}
	return nil
}

func ExtractSessionID(claims jwt.MapClaims) string {
	if sessionID, ok := claims["session_id"].(string); ok {
		return sessionID
	}
	// Also try "jti" (JWT ID)
	if jti, ok := claims["jti"].(string); ok {
		return jti
	}
	return ""
}
