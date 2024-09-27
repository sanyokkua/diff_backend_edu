package ua.kostenko.tasks.app.filter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import ua.kostenko.tasks.app.config.custom.UserAuthentication;
import ua.kostenko.tasks.app.repository.UserRepository;
import ua.kostenko.tasks.app.service.JwtService;

import java.io.IOException;
import java.util.Objects;

/**
 * A filter that intercepts each incoming HTTP request and checks for a valid JWT token.
 * <p>
 * If a valid JWT token is found and authenticated, the filter sets the authenticated user in the
 * {@link SecurityContextHolder} for the current request.
 * </p>
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class JwtRequestFilter extends OncePerRequestFilter {

    private final JwtService jwtUtil;
    private final UserRepository userRepository;

    /**
     * Extracts the JWT token from the "Authorization" header of the HTTP request.
     * <p>
     * The JWT token is expected to start with the "Bearer " prefix. If the header does not contain
     * the token or it is invalid, an empty string is returned.
     * </p>
     *
     * @param request the incoming HTTP request
     *
     * @return the JWT token if present, otherwise an empty string
     */
    private String extractJwt(HttpServletRequest request) {
        final String authorizationHeader = request.getHeader("Authorization");
        if (Objects.nonNull(authorizationHeader) && authorizationHeader.startsWith("Bearer ")) {
            String jwt = authorizationHeader.substring(7);
            log.debug("JWT extracted from Authorization header: {}", jwt);
            return jwt;
        }
        log.debug("Authorization header is missing or does not contain a Bearer token.");
        return "";
    }

    /**
     * Extracts the user's email (subject) from the JWT token.
     * <p>
     * If the JWT is valid, the email is extracted from the token. If the token is invalid or
     * blank, an empty string is returned.
     * </p>
     *
     * @param jwt the JWT token
     *
     * @return the extracted email, or an empty string if extraction fails
     */
    private String extractEmail(String jwt) {
        if (!jwt.isBlank()) {
            try {
                String email = jwtUtil.extractClaims(jwt).getSubject();
                log.debug("Email extracted from JWT: {}", email);
                return email;
            } catch (Exception e) {
                log.warn("Failed to extract email from JWT: {}", e.getMessage());
                return "";
            }
        }
        log.debug("JWT is blank, cannot extract email.");
        return "";
    }

    /**
     * Filters each HTTP request to validate the JWT token and set authentication if the token is valid.
     * <p>
     * This method extracts the JWT token, validates it, and if valid, sets the authenticated user in the
     * {@link SecurityContextHolder}. If the token is invalid or absent, the request continues without
     * authentication.
     * </p>
     *
     * @param req   the HTTP request
     * @param res   the HTTP response
     * @param chain the filter chain to proceed with the next filters
     *
     * @throws ServletException if an error occurs during request filtering
     * @throws IOException      if an I/O error occurs
     */
    @SuppressWarnings("NullableProblems")
    @Override
    protected void doFilterInternal(HttpServletRequest req, HttpServletResponse res,
                                    FilterChain chain) throws ServletException, IOException {
        log.info("Processing JWT authentication for request: {}", req.getRequestURI());

        // Extract JWT and email from the request
        String jwt = extractJwt(req);
        String email = extractEmail(jwt);

        // Check if the current security context already has an authenticated user
        boolean isAuthAbsent = Objects.isNull(SecurityContextHolder.getContext().getAuthentication());
        // Validate that we have extracted a non-blank email and the JWT token is valid
        boolean isEmailExtracted = !email.isBlank();
        boolean isTokenValid = jwtUtil.validateToken(jwt, email);

        if (isAuthAbsent && isEmailExtracted && isTokenValid) {
            log.debug("JWT is valid, attempting to authenticate user with email: {}", email);

            // Fetch the user associated with the extracted email
            var foundUser = userRepository.findByEmail(email);

            if (foundUser.isPresent()) {
                log.debug("User found for email: {}, setting security context", email);
                var authObj =
                        UserAuthentication.builder().user(foundUser.get()).jwtToken(jwt).isAuthenticated(true).build();

                // Set the authentication token in the security context
                SecurityContextHolder.getContext().setAuthentication(authObj);
                log.info("Successfully authenticated user: {}", email);
            } else {
                log.warn("No user found for email: {}", email);
            }
        } else {
            if (!isEmailExtracted) {
                log.warn("Email extraction failed or JWT is blank.");
            }
            if (!isTokenValid) {
                log.warn("JWT is invalid for email: {}", email);
            }
            if (!isAuthAbsent) {
                log.debug("Authentication already exists in security context, skipping JWT authentication.");
            }
        }

        // Proceed with the filter chain
        log.debug("Continuing the filter chain for request: {}", req.getRequestURI());
        chain.doFilter(req, res);
    }
}
