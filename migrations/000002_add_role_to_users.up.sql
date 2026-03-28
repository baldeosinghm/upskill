ALTER TABLE users
   ADD COLUMN role VARCHAR (50) NOT NULL CHECK (role IN ('teacher', 'student'));