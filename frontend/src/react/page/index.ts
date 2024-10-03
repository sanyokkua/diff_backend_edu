/**
 * Exports various components and helper functions used throughout the application.
 *
 * This module serves as an index file for easier imports of commonly used components.
 * It exports the following components:
 * - Dashboard: The main dashboard component for authenticated users.
 * - Login: The login page component for user authentication.
 * - PrivateRoute: A helper component to restrict access to certain routes for authenticated users.
 * - PublicRoute: A helper component to restrict access to certain routes for unauthenticated users.
 * - Profile: The user profile page component for managing user settings.
 * - Register: The registration page component for new user sign-up.
 * - ErrorPage: A component for displaying error messages or pages when a route is not found.
 *
 * @module
 * @exports {JSX.Element} Dashboard - Main dashboard component for authenticated users.
 * @exports {JSX.Element} Login - Login page component for user authentication.
 * @exports {JSX.Element} PrivateRoute - Helper component to restrict access for authenticated users.
 * @exports {JSX.Element} PublicRoute - Helper component to restrict access for unauthenticated users.
 * @exports {JSX.Element} Profile - User profile page component for managing user settings.
 * @exports {JSX.Element} Register - Registration page component for new user sign-up.
 * @exports {JSX.Element} ErrorPage - Component for displaying error messages or pages.
 */

export { default as Dashboard }    from "./Dashboard";    // Main dashboard component for authenticated users.
export { default as Login }        from "./LoginPage";   // Login page component for user authentication.
export { default as PrivateRoute } from "./helper/PrivateRoute"; // Helper component for private route protection.
export { default as PublicRoute }  from "./helper/PublicRoute";  // Helper component for public route protection.
export { default as Profile }      from "./ProfilePage"; // User profile page component for managing user settings.
export { default as Register }     from "./RegistrationPage"; // Registration page for new users.
export { default as ErrorPage }    from "./ErrorPage";   // Component for displaying error messages or not found pages.
