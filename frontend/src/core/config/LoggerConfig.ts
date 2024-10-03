import log    from "loglevel";
import prefix from "loglevel-plugin-prefix";


/**
 * The log level to be used for logging.
 * @constant {string}
 */
const logLevel = "DEBUG";

// Register the prefix plugin with loglevel
prefix.reg(log);

// Enable all log levels
log.enableAll();

// Set the log level
log.setLevel(logLevel);

// Apply the prefix settings to loglevel
prefix.apply(log, {
    template: "[%t] %l (%n):",
    levelFormatter(level) {
        return level.toUpperCase();
    },
    nameFormatter(name) {
        return name ?? "[anonymous]";
    },
    timestampFormatter(date) {
        return date.toISOString();
    }
});

// Log an informational message indicating the logger has been initialized
log.info("Initialized logger");

/**
 * Exported loglevel instance with applied prefix settings.
 */
export const LogLevel = log;
