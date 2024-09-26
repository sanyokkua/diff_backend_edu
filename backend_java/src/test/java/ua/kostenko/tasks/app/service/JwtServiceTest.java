package ua.kostenko.tasks.app.service;

import io.jsonwebtoken.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.MockitoAnnotations;
import org.mockito.junit.jupiter.MockitoExtension;
import ua.kostenko.tasks.app.exception.InvalidJwtTokenException;

import javax.crypto.SecretKey;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class JwtServiceTest {

    private static final String TEST_TOKEN = "testToken";
    private static final String TEST_USERNAME = "testUser";
    @InjectMocks
    private JwtService jwtService;
    @Mock
    private SecretKey jwtMainKey;
    @Mock
    private JwtParser jwtParser;
    @Mock
    private JwtParserBuilder jwtParserBuilder;
    @Mock
    private Claims claims;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        try (MockedStatic<Jwts> mocked = mockStatic(Jwts.class)) {
            mocked.when(Jwts::parser).thenReturn(jwtParserBuilder);
            assertEquals(jwtParserBuilder, Jwts.parser());
        }
    }

    @Test
    void extractClaims_ValidToken_ShouldReturnClaims() {
        // Arrange
        Jws<Claims> signedClaims = mock(Jws.class);
        when(jwtParserBuilder.verifyWith(jwtMainKey)).thenReturn(jwtParserBuilder);
        when(jwtParserBuilder.build()).thenReturn(jwtParser);
        when(jwtParser.parseSignedClaims(TEST_TOKEN)).thenReturn(signedClaims);
        when(signedClaims.getPayload()).thenReturn(claims);

        // Act
        Claims result = jwtService.extractClaims(TEST_TOKEN);

        // Assert
        assertEquals(claims, result);
        verify(jwtParser).parseSignedClaims(TEST_TOKEN);
    }

    @Test
    void extractClaims_ExpiredToken_ShouldReturnClaims() {
        // Arrange
        ExpiredJwtException expiredJwtException = mock(ExpiredJwtException.class);
        when(jwtParserBuilder.verifyWith(jwtMainKey)).thenReturn(jwtParserBuilder);
        when(jwtParserBuilder.build()).thenReturn(jwtParser);
        when(jwtParser.parseSignedClaims(TEST_TOKEN)).thenThrow(expiredJwtException);
        when(expiredJwtException.getClaims()).thenReturn(claims);

        // Act
        Claims result = jwtService.extractClaims(TEST_TOKEN);

        // Assert
        assertEquals(claims, result);
        verify(jwtParser).parseSignedClaims(TEST_TOKEN);
    }

    @Test
    void extractClaims_InvalidToken_ShouldThrowException() {
        // Arrange
        RuntimeException runtimeException = new RuntimeException("Invalid token");
        when(jwtParserBuilder.verifyWith(jwtMainKey)).thenReturn(jwtParserBuilder);
        when(jwtParserBuilder.build()).thenReturn(jwtParser);
        when(jwtParser.parseSignedClaims(TEST_TOKEN)).thenThrow(runtimeException);

        // Act & Assert
        InvalidJwtTokenException exception =
                assertThrows(InvalidJwtTokenException.class, () -> jwtService.extractClaims(TEST_TOKEN));

        assertEquals("Invalid JWT token: Invalid token", exception.getMessage());
    }

    @Test
    void isTokenExpired_ExpiredToken_ShouldReturnTrue() {
        // Arrange
        Date expirationDate = Date.from(Instant.now().minusSeconds(60)); // 1 minute ago
        when(claims.getExpiration()).thenReturn(expirationDate);

        // Act
        boolean result = jwtService.isTokenExpired(claims);

        // Assert
        assertTrue(result);
    }

    @Test
    void isTokenExpired_NonExpiredToken_ShouldReturnFalse() {
        // Arrange
        Date expirationDate = Date.from(Instant.now().plusSeconds(60)); // 1 minute in the future
        when(claims.getExpiration()).thenReturn(expirationDate);

        // Act
        boolean result = jwtService.isTokenExpired(claims);

        // Assert
        assertFalse(result);
    }

    @Test
    void isTokenExpired_NullExpiration_ShouldReturnTrue() {
        // Arrange
        when(claims.getExpiration()).thenReturn(null);

        // Act
        boolean result = jwtService.isTokenExpired(claims);

        // Assert
        assertTrue(result);
    }

    @Test
    void validateToken_ValidTokenAndUsername_ShouldReturnTrue() {
        // Arrange
        when(jwtService.extractClaims(TEST_TOKEN)).thenReturn(claims);
        when(claims.getSubject()).thenReturn(TEST_USERNAME);
        when(jwtService.isTokenExpired(claims)).thenReturn(false);

        // Act
        boolean result = jwtService.validateToken(TEST_TOKEN, TEST_USERNAME);

        // Assert
        assertTrue(result);
    }

    @Test
    void validateToken_ValidTokenAndNonMatchingUsername_ShouldReturnFalse() {
        // Arrange
        when(jwtService.extractClaims(TEST_TOKEN)).thenReturn(claims);
        when(claims.getSubject()).thenReturn("otherUser");
        when(jwtService.isTokenExpired(claims)).thenReturn(false);

        // Act
        boolean result = jwtService.validateToken(TEST_TOKEN, TEST_USERNAME);

        // Assert
        assertFalse(result);
    }

    @Test
    void validateToken_NullToken_ShouldReturnFalse() {
        // Act
        boolean result = jwtService.validateToken(null, TEST_USERNAME);

        // Assert
        assertFalse(result);
    }

    @Test
    void validateToken_NullUsername_ShouldReturnFalse() {
        // Act
        boolean result = jwtService.validateToken(TEST_TOKEN, null);

        // Assert
        assertFalse(result);
    }

    @Test
    void validateToken_NullUsernameClaim_ShouldReturnFalse() {
        // Arrange
        when(jwtService.extractClaims(TEST_TOKEN)).thenReturn(claims);
        when(claims.getSubject()).thenReturn(null);
        when(jwtService.isTokenExpired(claims)).thenReturn(false);

        // Act
        boolean result = jwtService.validateToken(TEST_TOKEN, TEST_USERNAME);

        // Assert
        assertFalse(result);
    }

    @Test
    void generateJwtToken_ValidUsername_ShouldReturnToken() {
        // Arrange
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime exp = now.plusMinutes(15);
        Date nowDate = Date.from(now.atZone(ZoneId.systemDefault()).toInstant());
        Date expDate = Date.from(exp.atZone(ZoneId.systemDefault()).toInstant());

        // Act
        String token = jwtService.generateJwtToken(TEST_USERNAME);

        // Assert
        assertNotNull(token);
        // Further checks could be added to decode the token and check its claims, expiration, etc.
    }
}