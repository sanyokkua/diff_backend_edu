package ua.kostenko.tasks.app.service;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import ua.kostenko.tasks.app.exception.InvalidJwtTokenException;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;

import static org.junit.jupiter.api.Assertions.*;

class JwtServiceTest {

    private static final String TEST_USERNAME = "testUser";
    private static final String INVALID_USERNAME = "invalidUser";
    private static final String INVALID_TOKEN = "invalidToken";
    private static final int TOKEN_EXPIRATION_MINUTES = 15;

    private String validToken;
    private String expiredToken;
    private SecretKey secretKey;
    private JwtService jwtService;

    @BeforeEach
    void setUp() {
        secretKey =
                Keys.hmacShaKeyFor("q33dvfty23gdfty2dyvbyuewbduyeytgvdfygvwytefvyt".getBytes(StandardCharsets.UTF_8));
        jwtService = new JwtService(secretKey);
        validToken = generateToken(TEST_USERNAME, TOKEN_EXPIRATION_MINUTES);
        expiredToken = generateToken(TEST_USERNAME, -TOKEN_EXPIRATION_MINUTES);
    }

    /**
     * Generates a JWT token for the provided username with a specified expiration offset in minutes.
     * If the offset is positive, the token is valid; if negative, the token is expired.
     */
    private String generateToken(String username, int expirationOffsetMinutes) {
        LocalDateTime timeNow = LocalDateTime.now();
        LocalDateTime timeExp = timeNow.plusMinutes(expirationOffsetMinutes);
        Instant instantNow = timeNow.atZone(ZoneId.systemDefault()).toInstant();
        Instant instantExp = timeExp.atZone(ZoneId.systemDefault()).toInstant();
        Date dateNow = Date.from(instantNow);
        Date dateExp = Date.from(instantExp);

        return Jwts.builder().subject(username).issuedAt(dateNow).expiration(dateExp).signWith(secretKey).compact();
    }

    @Test
    void extractClaims_ValidToken_ShouldReturnClaims() {
        Claims result = jwtService.extractClaims(validToken);

        assertNotNull(result);
        assertEquals(TEST_USERNAME, result.getSubject());
    }

    @Test
    void extractClaims_ExpiredToken_ShouldReturnClaims() {
        Claims result = jwtService.extractClaims(expiredToken);

        assertNotNull(result);
        assertEquals(TEST_USERNAME, result.getSubject());
    }

    @Test
    void extractClaims_InvalidToken_ShouldThrowInvalidJwtTokenException() {
        assertThrows(InvalidJwtTokenException.class, () -> jwtService.extractClaims(INVALID_TOKEN));
    }

    @Test
    void isTokenExpired_ValidToken_ShouldReturnFalse() {
        Claims claims = jwtService.extractClaims(validToken);
        assertFalse(jwtService.isTokenExpired(claims));
    }

    @Test
    void isTokenExpired_ExpiredToken_ShouldReturnTrue() {
        Claims claims = jwtService.extractClaims(expiredToken);
        assertTrue(jwtService.isTokenExpired(claims));
    }

    @Test
    void isTokenExpired_NoExpirationDate_ShouldReturnTrue() {
        Claims claims = Jwts.claims().subject(TEST_USERNAME).build();
        assertTrue(jwtService.isTokenExpired(claims));
    }

    @Test
    void validateToken_ValidTokenAndCorrectUsername_ShouldReturnTrue() {
        assertTrue(jwtService.validateToken(validToken, TEST_USERNAME));
    }

    @Test
    void validateToken_ValidTokenAndIncorrectUsername_ShouldReturnFalse() {
        assertFalse(jwtService.validateToken(validToken, INVALID_USERNAME));
    }

    @Test
    void validateToken_ExpiredToken_ShouldReturnFalse() {
        assertFalse(jwtService.validateToken(expiredToken, TEST_USERNAME));
    }

    @Test
    void validateToken_NullToken_ShouldReturnFalse() {
        assertFalse(jwtService.validateToken(null, TEST_USERNAME));
    }

    @Test
    void validateToken_NullUsername_ShouldReturnFalse() {
        assertFalse(jwtService.validateToken(validToken, null));
    }

    @Test
    void validateToken_NullTokenAndUsername_ShouldReturnFalse() {
        assertFalse(jwtService.validateToken(null, null));
    }

    @Test
    void validateToken_NullSubject_ShouldReturnFalse() {
        var token = generateToken(null, TOKEN_EXPIRATION_MINUTES);
        assertFalse(jwtService.validateToken(token, TEST_USERNAME));
    }

    @Test
    void generateJwtToken_ShouldReturnValidToken() {
        String generatedToken = jwtService.generateJwtToken(TEST_USERNAME);

        Claims claims = jwtService.extractClaims(generatedToken);
        assertEquals(TEST_USERNAME, claims.getSubject());
        assertFalse(jwtService.isTokenExpired(claims));
    }
}