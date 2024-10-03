import { LogLevel } from "../config";


const log = LogLevel.getLogger("BrowserStore");

/**
 * Enum for item types stored in the browser.
 * @enum {string}
 */
export enum ItemType {
    JWT_TOKEN = "JWT_USER_TOKEN",
    USER_EMAIL = "USER_EMAIL",
    USER_ID = "USER_ID"
}

/**
 * Class for managing browser storage operations.
 * @class
 */
class BrowserStore {
    /**
     * Retrieves data from local storage.
     * @function getData
     * @memberof BrowserStore
     * @param {ItemType} itemType - The type of item to retrieve.
     * @returns {string | null} - The retrieved data or null if not found.
     */
    getData(itemType: ItemType): string | null {
        log.info(`Getting ${ itemType } from storage`);
        const value: string | null = localStorage.getItem(itemType);
        log.debug(`Retrieved ${ itemType }: ${ value }`);
        return value;
    }

    /**
     * Saves data to local storage.
     * @function saveData
     * @memberof BrowserStore
     * @param {ItemType} itemType - The type of item to save.
     * @param {string} dataValue - The data to save.
     * @returns {void}
     */
    saveData(itemType: ItemType, dataValue: string): void {
        log.info(`Saving ${ itemType } to storage`);
        localStorage.setItem(itemType, dataValue);
        log.debug(`Saved ${ itemType }: ${ dataValue }`);
    }
}

export default BrowserStore;
