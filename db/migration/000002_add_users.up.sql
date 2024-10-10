CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password_change_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:0:00Z',
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");