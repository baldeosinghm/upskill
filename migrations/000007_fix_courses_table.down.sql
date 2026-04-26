ALTER TABLE courses RENAME COLUMN name TO title;
ALTER TABLE courses RENAME COLUMN owner_id TO teacher_id;
ALTER TABLE courses
    ALTER COLUMN teacher_id DROP NOT NULL,
    DROP CONSTRAINT courses_owner_id_fkey,
    ADD CONSTRAINT courses_teacher_id_fkey
        FOREIGN KEY (teacher_id)
        REFERENCES users(id);
ALTER TABLE courses DROP COLUMN updated_at;