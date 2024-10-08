CREATE SCHEMA "backend_diff";

CREATE TABLE "backend_diff"."users" (
    "user_id"       BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "email"         VARCHAR(255) NOT NULL UNIQUE,
    "password_hash" VARCHAR(255) NOT NULL
                                    );

CREATE TABLE "backend_diff"."tasks" (
    "task_id"     BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "user_id"     BIGINT       NOT NULL REFERENCES "backend_diff"."users" ("user_id") ON DELETE CASCADE,
    "name"        VARCHAR(255) NOT NULL,
    "description" TEXT         NOT NULL,
    UNIQUE ("name", "user_id")
                                    );

CREATE UNIQUE INDEX "idx_user_email" ON "backend_diff"."users" ("email");
CREATE INDEX "idx_task_name" ON "backend_diff"."tasks" ("name");
