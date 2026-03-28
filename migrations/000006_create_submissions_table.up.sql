CREATE TABLE IF NOT EXISTS submissions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id UUID REFERENCES assignments(id),
    student_id UUID REFERENCES users(id),
    file_url TEXT,
    status TEXT DEFAULT 'submitted',
    grade INT,
    graded_by UUID REFERENCES users(id),
    submitted_at TIMESTAMPTZ DEFAULT now(),
    graded_at TIMESTAMPTZ
);