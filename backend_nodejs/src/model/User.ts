import { Column, Entity, Index, OneToMany, PrimaryGeneratedColumn } from "typeorm";
import { Task }                                                     from ".";


@Entity({ schema: "backend_diff", name: "users" })
@Index("idx_user_email", ["email"], { unique: true })
export class User {
    @PrimaryGeneratedColumn({ name: "user_id" })
    id?: number;

    @Column({ type: "varchar", length: 255, unique: true })
    email!: string;

    @Column({ name: "password_hash", type: "varchar", length: 255 })
    passwordHash!: string;

    @OneToMany(() => Task, task => task.user)
    tasks?: Task[];
}
