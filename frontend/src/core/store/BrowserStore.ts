import { LogLevel } from "../config";


const log = LogLevel.getLogger("BrowserStore");

export enum ItemType {
    JWT_TOKEN = "JWT_USER_TOKEN",
    USER_EMAIL = "USER_EMAIL",
    USER_ID = "USER_ID"
}

class BrowserStore {
    getData(itemType: ItemType): string | null {
        log.info(`Getting ${ itemType } from storage`);
        const value: string | null = localStorage.getItem(itemType);
        log.debug(`Retrieved ${ itemType }: ${ value }`);
        return value;
    }

    saveData(itemType: ItemType, dataValue: string): void {
        log.info(`Saving ${ itemType } to storage`);
        localStorage.setItem(itemType, dataValue);
        log.debug(`Saved ${ itemType }: ${ dataValue }`);
    }
}

export default BrowserStore;