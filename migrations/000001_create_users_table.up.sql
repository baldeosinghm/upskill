CREATE TABLE IF NOT EXISTS users(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username VARCHAR (50) UNIQUE NOT NULL,
   password TEXT NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
   updated_at TIMESTAMPTZ DEFAULT now()
);