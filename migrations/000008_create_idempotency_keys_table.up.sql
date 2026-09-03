CREATE TABLE IF NOT EXISTS idempotency_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- The client-supplied key. UNIQUE is the whole game: it's what makes
    -- the "claim" atomic. Two concurrent inserts with the same key -> the
    -- database rejects the second one. No application-level check needed.
    idempotency_key TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
           CHECK (status IN ('pending', 'completed')),
    -- What we replay on a repeat request once status = 'completed'.
    response_code INT,
    response_body JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)