package middleware

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// extracting the request headers
// extracting our claims
// validating everything

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		//	extract the headers

		token := extractTokenFromHeaders(request.Headers)

		if token == "" {
			return events.APIGatewayProxyResponse{
				Body:       "Missing Auth Token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		//	parse the token
		claims, err := parseToken(token)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "Invalid Token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "Token Expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := "secret"
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("could not parse claim")
	}

	return claims, nil
}
