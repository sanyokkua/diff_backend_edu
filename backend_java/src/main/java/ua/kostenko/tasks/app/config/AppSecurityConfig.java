package ua.kostenko.tasks.app.config;

import io.jsonwebtoken.security.Keys;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.access.AccessDeniedHandler;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import ua.kostenko.tasks.app.filter.JwtRequestFilter;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.List;

/**
 * Configuration class for security settings in the application.
 * <p>
 * This class is responsible for configuring Spring Security, including authentication, authorization,
 * and CORS settings.
 * </p>
 */
@Slf4j
@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class AppSecurityConfig {

    @Value("${app.jwt.secret}")
    private String jwtSecretKey;

    /**
     * Provides a PasswordEncoder bean using BCrypt hashing algorithm.
     *
     * @return a {@link PasswordEncoder} instance
     */
    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    /**
     * Provides a SecretKey bean for JWT signing and verification.
     *
     * @return a {@link SecretKey} instance created from the JWT secret key
     */
    @Bean
    public SecretKey jwtSecretKey() {
        return Keys.hmacShaKeyFor(jwtSecretKey.getBytes(StandardCharsets.UTF_8));
    }

    /**
     * Configures the security filter chain for the application.
     *
     * @param http                           the {@link HttpSecurity} to be configured
     * @param jwtRequestFilter               the JWT request filter to be added to the filter chain
     * @param customAuthenticationEntryPoint the custom entry point for handling unauthorized access
     * @param customAccessDeniedHandler      the custom handler for access denied exceptions
     *
     * @return a configured {@link SecurityFilterChain}
     *
     * @throws Exception if an error occurs during the configuration
     */
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http, JwtRequestFilter jwtRequestFilter,
                                           AuthenticationEntryPoint customAuthenticationEntryPoint,
                                           AccessDeniedHandler customAccessDeniedHandler) throws Exception {
        http.csrf(AbstractHttpConfigurer::disable);
        http.cors(Customizer.withDefaults());
        http.httpBasic(AbstractHttpConfigurer::disable);

        http.authorizeHttpRequests(auth -> auth.requestMatchers("/",
                                                                "/public/**",
                                                                "/static/**",
                                                                "/js/**",
                                                                "api/v1/auth/login",
                                                                "api/v1/auth/register",
                                                                "/v3/api-docs/**",
                                                                "/swagger-ui/**",
                                                                "/swagger-ui.html")
                                               .permitAll()
                                               .anyRequest()
                                               .authenticated());

        http.sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS));
        http.addFilterBefore(jwtRequestFilter, UsernamePasswordAuthenticationFilter.class);

        // Configure custom exception handling
        http.exceptionHandling(exception -> {
            exception.authenticationEntryPoint(customAuthenticationEntryPoint);
            exception.accessDeniedHandler(customAccessDeniedHandler);
        });

        // Log that the security filter chain has been configured
        log.info("Security filter chain has been configured successfully.");

        return http.build();
    }

    /**
     * Configures CORS settings for the application.
     *
     * @return a {@link CorsConfigurationSource} that defines the CORS settings
     */
    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration configuration = new CorsConfiguration();
        configuration.setAllowedOrigins(Arrays.asList("http://localhost:5173", "http://localhost:8080"));
        configuration.setAllowedMethods(Arrays.asList("GET", "POST", "PUT", "DELETE"));
        configuration.setAllowedHeaders(List.of("Authorization", "*", "Content-Type"));
        configuration.setAllowCredentials(true);

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", configuration);

        // Log the allowed origins for CORS
        log.info("CORS configuration set with allowed origins: {}", configuration.getAllowedOrigins());

        return source;
    }
}
