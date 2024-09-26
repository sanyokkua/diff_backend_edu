package ua.kostenko.tasks.app;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Main application class for the Spring Boot application.
 * <p>
 * This class serves as the entry point for the Spring Boot application. It contains the main method
 * which uses {@link SpringApplication#run(Class, String...)} to launch the application.
 * </p>
 */
@SpringBootApplication
public class AppApplication {

    /**
     * Main method which serves as the entry point for the Spring Boot application.
     *
     * @param args command-line arguments passed to the application.
     */
    public static void main(String[] args) {
        SpringApplication.run(AppApplication.class, args);
    }
}
