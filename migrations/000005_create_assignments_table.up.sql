CREATE TABLE IF NOT EXISTS assignments(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id),
    title TEXT NOT NULL,
    description TEXT,
    due_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now()
);