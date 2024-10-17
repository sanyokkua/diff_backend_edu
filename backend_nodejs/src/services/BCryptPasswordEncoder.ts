import bcrypt               from "bcrypt";
import { IPasswordEncoder } from "../api";


export const DEFAULT_COST = 10;

export class BCryptPasswordEncoder implements IPasswordEncoder {
    private readonly cost: number;

    constructor(cost: number) {
        this.cost = cost;
    }

    matches(rawPassword: string, encodedPassword: string): boolean {
        return bcrypt.compareSync(rawPassword, encodedPassword);
    }

    encode(rawPassword: string): string {
        return bcrypt.hashSync(rawPassword, this.cost);
    }

}