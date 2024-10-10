-- Step 1: Drop the foreign key constraint on "accounts" table
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Step 2: Drop the unique index on ("owner", "currency") in "accounts"
DROP INDEX IF EXISTS "accounts_owner_currency_idx";

-- Step 3: Drop the "users" table
DROP TABLE IF EXISTS "users";