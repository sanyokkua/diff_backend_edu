package ua.kostenko.tasks.app.service;

import io.jsonwebtoken.*;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import ua.kostenko.tasks.app.exception.InvalidJwtTokenException;

import javax.crypto.SecretKey;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;

/**
 * A service class for handling operations related to JSON Web Tokens (JWT).
 * <p>
 * This class provides methods for generating, validating, and extracting claims from JWT tokens,
 * as well as checking their expiration status.
 * </p>
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class JwtService {

    /**
     * The expiration time of JWT tokens in minutes.
     */
    private static final int JWT_EXP_MINUTES = 15;

    /**
     * The main secret key used for signing and verifying JWT tokens.
     */
    private final SecretKey jwtMainKey;

    /**
     * Extracts the claims from the provided JWT token.
     * <p>
     * This method parses the JWT token using the provided secret key. If the token is expired,
     * it returns the claims from the expired token. If the token is invalid, an exception is thrown.
     * </p>
     *
     * @param token the JWT token
     *
     * @return the claims contained in the token
     *
     * @throws RuntimeException if the token is invalid
     */
    public Claims extractClaims(String token) {
        try {
            JwtParser parser = Jwts.parser().verifyWith(jwtMainKey).build();
            Jws<Claims> signedClaims = parser.parseSignedClaims(token);
            log.debug("Successfully extracted claims from token.");
            return signedClaims.getPayload();
        } catch (ExpiredJwtException ex) {
            log.warn("JWT token is expired: {}", ex.getMessage());
            return ex.getClaims(); // Return claims for further processing
        } catch (Exception ex) {
            log.error("Invalid JWT token: {}", ex.getMessage());
            throw new InvalidJwtTokenException("Invalid JWT token: " + ex.getMessage());
        }
    }

    /**
     * Checks if the JWT token is expired based on its claims.
     *
     * @param claims the claims extracted from the JWT token
     *
     * @return {@code true} if the token is expired, {@code false} otherwise
     */
    public boolean isTokenExpired(Claims claims) {
        Date expirationDate = claims.getExpiration();
        if (expirationDate == null) {
            log.warn("Token expiration date is missing.");
            return true;
        }

        // Get the current date
        LocalDateTime localDateTime = LocalDateTime.now();
        Instant instant = localDateTime.atZone(ZoneId.systemDefault()).toInstant();
        Date currentDate = Date.from(instant);

        boolean isExpired = expirationDate.before(currentDate);
        log.debug("Token expiration check: expired = {}", isExpired);
        return isExpired;
    }

    /**
     * Validates the provided JWT token against the provided username.
     * <p>
     * This method checks that the token is not expired and the username in the token matches
     * the provided username.
     * </p>
     *
     * @param token    the JWT token
     * @param username the username to validate against the token
     *
     * @return {@code true} if the token is valid and the username matches, {@code false} otherwise
     */
    public boolean validateToken(String token, String username) {
        if (token == null || username == null) {
            log.warn("Token or username is null during validation.");
            return false;
        }

        try {
            Claims claims = extractClaims(token);
            String claimUsername = claims.getSubject();

            if (claimUsername == null) {
                log.warn("Username claim is null in the token.");
                return false;
            }

            boolean isUsernameValid = username.equals(claimUsername);
            boolean isTokenValid = !isTokenExpired(claims);

            log.debug("Token validation result: usernameValid = {}, tokenValid = {}", isUsernameValid, isTokenValid);
            return isUsernameValid && isTokenValid;
        } catch (RuntimeException ex) {
            log.warn("Failed to validate JWT token: {}", ex.getMessage());
            return false;
        }
    }

    /**
     * Generates a new JWT token for the specified username.
     * <p>
     * The token is signed using the secret key and has an expiration time of 15 minutes.
     * </p>
     *
     * @param username the username to include in the token's subject claim
     *
     * @return the generated JWT token as a compact string
     */
    public String generateJwtToken(String username) {
        // Get the current time and the expiration time
        LocalDateTime timeNow = LocalDateTime.now();
        LocalDateTime timeExp = timeNow.plusMinutes(JWT_EXP_MINUTES);
        Instant instantNow = timeNow.atZone(ZoneId.systemDefault()).toInstant();
        Instant instantExp = timeExp.atZone(ZoneId.systemDefault()).toInstant();
        Date dateNow = Date.from(instantNow);
        Date dateExp = Date.from(instantExp);

        // Generate the token
        String token =
                Jwts.builder().subject(username).issuedAt(dateNow).expiration(dateExp).signWith(jwtMainKey).compact();

        log.info("Generated JWT token for username: {}, valid until: {}", username, dateExp);
        return token;
    }
}
