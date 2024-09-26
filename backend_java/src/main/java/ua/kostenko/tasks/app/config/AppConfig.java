package ua.kostenko.tasks.app.config;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

/**
 * Configuration class for setting up web-related configurations in the Spring application.
 * <p>
 * This class implements the {@link WebMvcConfigurer} interface to customize the default behavior of Spring MVC.
 * In particular, it configures resource handlers for serving static resources such as the Swagger UI.
 * </p>
 */
@Slf4j
@Configuration
public class AppConfig implements WebMvcConfigurer {

    /**
     * Adds resource handlers for serving static resources.
     * <p>
     * This method registers a resource handler for serving the Swagger UI documentation.
     * It maps requests to "/swagger-ui/**" to the appropriate location in the classpath where the Swagger UI resources
     * are stored.
     * </p>
     *
     * @param registry the {@link ResourceHandlerRegistry} to add resource handlers to
     */
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        registry.addResourceHandler("/swagger-ui/**")
                .addResourceLocations("classpath:/META-INF/resources/webjars/springdoc-openapi-ui/")
                .resourceChain(false);

        // This is useful for debugging and ensuring that the correct path is set up
        // Check http://localhost:8080/swagger-ui/index.html for API Description
        log.info("Swagger UI available at: http://localhost:8080/swagger-ui/index.html");
    }
}
