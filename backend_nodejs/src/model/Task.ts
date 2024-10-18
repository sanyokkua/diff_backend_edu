import { Column, Entity, Index, ManyToOne, PrimaryGeneratedColumn } from "typeorm";
import { User }                                                     from ".";


@Entity({ schema: "backend_diff", name: "tasks" })
@Index("idx_task_name", ["name"])
@Index(["name", "user"], { unique: true })
export class Task {
    @PrimaryGeneratedColumn({ name: "task_id" })
    id?: number;

    @Column({ type: "varchar", length: 255 })
    name!: string;

    @Column({ type: "text" })
    description!: string;

    @ManyToOne(() => User, user => user.tasks, { onDelete: "CASCADE" })
    user!: User;
}
