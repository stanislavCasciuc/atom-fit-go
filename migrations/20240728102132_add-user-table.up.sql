CREATE TABLE IF NOT EXISTS users (
    id       SERIAL  PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT false,
    is_superuser BOOLEAN NOT NULL DEFAULT false,
    is_male BOOLEAN NOT NULL DEFAULT true,
    age INTEGER NOT NULL DEFAULT 20,
    height INTEGER NOT NULL DEFAULT 175,
    weight INTEGER NOT NULL DEFAULT 70,
    goal TEXT NOT NULL DEFAULT 'lose' CHECK (goal IN ('lose', 'maintain', 'gain')),
    weight_goal INTEGER NOT NULL DEFAULT 65
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);