package user

import (
	"crypto/rsa"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func verifyUserToken(path string) fiber.Handler {
	data, err := os.ReadFile(path)
	var key *rsa.PublicKey
	if err == nil {
		key, err = jwt.ParseRSAPublicKeyFromPEM(data)
	}
	return func(c *fiber.Ctx) error {
		if err != nil || key == nil {
			return fiber.NewError(500, "token verification is unavailable")
		}
		header := strings.TrimSpace(c.Get("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			return fiber.NewError(401, "missing bearer token")
		}
		claims := &TokenClaims{}
		token, parseErr := jwt.ParseWithClaims(strings.TrimPrefix(header, "Bearer "), claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
					return nil, fiber.NewError(401, "invalid token algorithm")
				}
				return key, nil
			}, jwt.WithIssuer("shopcore-user-service"), jwt.WithAudience("shopcore"))
		if parseErr != nil || !token.Valid {
			return fiber.NewError(401, "invalid or expired token")
		}
		c.Locals("userID", claims.Subject)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}
