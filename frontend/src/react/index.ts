/**
 * Exports all modules from the specified directories.
 *
 * This module re-exports all exports from the `component`, `layout`, `page`, and `route` directories.
 * It serves as a central export point, making it easier to import these modules in other parts of the application.
 *
 * @module CentralExports
 */

export * from "./component";
/**
 * Re-exports all exports from the `component` directory.
 *
 * This export makes all components available for import from the `component` directory.
 *
 * @exports Components
 */

export * from "./layout";
/**
 * Re-exports all exports from the `layout` directory.
 *
 * This export makes all layout components available for import from the `layout` directory.
 *
 * @exports Layouts
 */

export * from "./page";
/**
 * Re-exports all exports from the `page` directory.
 *
 * This export makes all page components available for import from the `page` directory.
 *
 * @exports Pages
 */

export * from "./route";
/**
 * Re-exports all exports from the `route` directory.
 *
 * This export makes all route components and helpers available for import from the `route` directory.
 *
 * @exports Routes
 */
