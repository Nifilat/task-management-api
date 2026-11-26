DELETE FROM tokens;
ALTER TABLE tokens ALTER COLUMN refresh_token TYPE UUID USING refresh_token::UUID;
