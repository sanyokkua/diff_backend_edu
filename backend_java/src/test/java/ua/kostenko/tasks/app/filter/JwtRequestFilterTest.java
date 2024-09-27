package ua.kostenko.tasks.app.filter;

import io.jsonwebtoken.Claims;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.context.SecurityContextImpl;
import ua.kostenko.tasks.app.config.custom.UserAuthentication;
import ua.kostenko.tasks.app.entity.User;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.JwtService;

import java.io.IOException;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.mockito.Mockito.*;

class JwtRequestFilterTest {

    @Mock
    private JwtService jwtUtil;

    @Mock
    private UserRepository authenticationService;

    @Mock
    private FilterChain filterChain;

    @InjectMocks
    private JwtRequestFilter jwtRequestFilter;

    private MockHttpServletRequest request;
    private MockHttpServletResponse response;

    @BeforeEach
    void setUp() {
        //noinspection resource
        MockitoAnnotations.openMocks(this);
        request = new MockHttpServletRequest();
        response = new MockHttpServletResponse();
        SecurityContextHolder.clearContext();
    }

    @Test
    void doFilterInternal_validToken_authenticatesUser() throws ServletException, IOException {
        request.addHeader("Authorization", "Bearer validJwtToken");
        String email = "test@example.com";

        Claims mockClaimsWithEmail = mockClaimsWithEmail(email);

        when(jwtUtil.extractClaims(anyString())).thenReturn(mockClaimsWithEmail);
        when(jwtUtil.validateToken("validJwtToken", email)).thenReturn(true);
        when(authenticationService.findByEmail(email)).thenReturn(Optional.of(mockUserDetails()));

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNotNull(SecurityContextHolder.getContext().getAuthentication());
        verify(jwtUtil).validateToken("validJwtToken", email);
        verify(authenticationService).findByEmail(email);
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_invalidToken_doesNotAuthenticate() throws ServletException, IOException {
        request.addHeader("Authorization", "Bearer invalidJwtToken");
        String email = "test@example.com";

        when(jwtUtil.extractClaims("invalidJwtToken")).thenThrow(new RuntimeException("Invalid token"));
        when(jwtUtil.validateToken("invalidJwtToken", email)).thenReturn(false);

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNull(SecurityContextHolder.getContext().getAuthentication());
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_missingAuthorizationHeader_proceedsWithoutAuthentication() throws ServletException, IOException {
        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNull(SecurityContextHolder.getContext().getAuthentication());
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_blankJwt_doesNotAuthenticate() throws ServletException, IOException {
        request.addHeader("Authorization", "Bearer ");

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNull(SecurityContextHolder.getContext().getAuthentication());
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_existingAuthentication_skipsAuthentication() throws ServletException, IOException {
        SecurityContextHolder.setContext(new SecurityContextImpl(mockAuthentication()));

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        verify(jwtUtil, times(0)).extractClaims(anyString());
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_emailNotExtracted_doesNotAuthenticate() throws ServletException, IOException {
        request.addHeader("Authorization", "Bearer someJwtToken");

        when(jwtUtil.extractClaims("someJwtToken")).thenThrow(new RuntimeException("Failed to extract email"));
        when(jwtUtil.validateToken("someJwtToken", "")).thenReturn(false);

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNull(SecurityContextHolder.getContext().getAuthentication());
        verify(filterChain).doFilter(request, response);
    }

    @Test
    void doFilterInternal_userNotFound_doesNotAuthenticate() throws ServletException, IOException {
        request.addHeader("Authorization", "Bearer validJwtToken");
        String email = "notfound@example.com";

        Claims mockClaimsWithEmail = mockClaimsWithEmail(email);
        when(jwtUtil.extractClaims("validJwtToken")).thenReturn(mockClaimsWithEmail);
        when(jwtUtil.validateToken("validJwtToken", email)).thenReturn(true);
        when(authenticationService.findByEmail(email)).thenReturn(Optional.empty());

        jwtRequestFilter.doFilterInternal(request, response, filterChain);

        assertNull(SecurityContextHolder.getContext().getAuthentication());
        verify(filterChain).doFilter(request, response);
    }

    private UserAuthentication mockAuthentication() {
        User userDetails = User.builder().email("user").passwordHash("password").build();
        return UserAuthentication.builder().user(userDetails).isAuthenticated(true).jwtToken("jwt").build();
    }

    private User mockUserDetails() {
        return User.builder().email("test@example.com").passwordHash("password").build();
    }

    private Claims mockClaimsWithEmail(String email) {
        Claims claims = mock(Claims.class);
        when(claims.getSubject()).thenReturn(email);
        return claims;
    }
}