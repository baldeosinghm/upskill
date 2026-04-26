ALTER TABLE courses RENAME COLUMN title TO name;
ALTER TABLE courses RENAME COLUMN teacher_id TO owner_id;
ALTER TABLE courses
    ALTER COLUMN owner_id SET NOT NULL,
    DROP CONSTRAINT courses_teacher_id_fkey,
    ADD CONSTRAINT courses_owner_id_fkey
        FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE RESTRICT;
ALTER TABLE courses ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();