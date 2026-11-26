ALTER TABLE tokens ALTER COLUMN refresh_token TYPE TEXT USING refresh_token::TEXT;
