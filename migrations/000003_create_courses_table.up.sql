CREATE TABLE IF NOT EXISTS courses(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   title TEXT NOT NULL,
   teacher_id UUID REFERENCES users(id),
   created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);