package middlewares

import (
	"go_backend/apperrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go_backend/api"
	"go_backend/utils"
)

// JWTAuthMiddleware returns a Gin middleware function that handles JWT authentication.
// It extracts the JWT from the Authorization header, validates it, and retrieves the corresponding user.
// Parameters:
//   - jwtService: The service responsible for handling JWT operations.
//   - userRepository: The repository used to fetch user details from the database.
//
// Returns:
//   - gin.HandlerFunc: The middleware function that performs the authentication.
func JWTAuthMiddleware(jwtService api.JwtService, userRepository api.UserRepository) gin.HandlerFunc {
	log.Debug().Msg("Initializing JWTAuthMiddleware")

	return func(ctx *gin.Context) {
		// Retrieve the Authorization header from the request context
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().Msg("Authorization header is missing")
			// If the header is missing, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewInsufficientAuthenticationError("authorization header required"))
			ctx.Abort()
			return
		}

		// Extract the JWT from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			log.Warn().Msg("JWT is missing from Authorization header")
			// If the token is missing, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewInsufficientAuthenticationError("bearer token required"))
			ctx.Abort()
			return
		}

		// Extract and validate the claims from the JWT
		claims, err := jwtService.ExtractClaims(tokenString)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to extract claims from JWT")
			// If claims extraction fails, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewInsufficientAuthenticationError("invalid token"))
			ctx.Abort()
			return
		}

		// Retrieve the subject (user email) from the claims
		subject, err := claims.GetSubject()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to extract subject from JWT")
			// If subject extraction fails, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewInsufficientAuthenticationError("invalid token subject"))
			ctx.Abort()
			return
		}

		// Validate the JWT and check if it has expired
		if !jwtService.ValidateToken(tokenString, subject) {
			log.Warn().Str("subject", subject).Msg("Invalid or expired JWT")
			// If the token is invalid or expired, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewInsufficientAuthenticationError("invalid or expired token"))
			ctx.Abort()
			return
		}

		// Fetch the user from the database using the subject (email)
		user, err := userRepository.GetUserByEmail(subject)
		if err != nil || user == nil {
			log.Warn().Str("subject", subject).Msg("No user found for subject")
			// If no user is found, return a 401 Unauthorized response
			utils.WriteErrorResponse(ctx, http.StatusUnauthorized, apperrors.NewAuthenticationCredentialsNotFoundError("user not found"))
			ctx.Abort()
			return
		}

		// Set the authenticated user in the context for later use
		log.Debug().Str("email", user.Email).Msg("User authenticated")
		ctx.Set("CurrentUser", user)

		// Proceed to the next middleware or handler in the chain
		ctx.Next()
	}
}
