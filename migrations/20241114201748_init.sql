-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS public.users (
    user_id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher', 'leader')),
    fullname TEXT,
    telegram TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS public.faculties (
    faculty_id BIGSERIAL PRIMARY KEY,
    faculty_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.programs (
    program_id BIGSERIAL PRIMARY KEY,
    faculty_id BIGINT NOT NULL,
    program_name TEXT NOT NULL UNIQUE,
    FOREIGN KEY (faculty_id) REFERENCES public.faculties (faculty_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.groups (
    group_id BIGSERIAL PRIMARY KEY,
    leader_id BIGINT NOT NULL,
    faculty_id BIGINT NOT NULL,
    program_id BIGINT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    short_name TEXT NOT NULL UNIQUE,
    exists_schedule BOOLEAN NOT NULL DEFAULT FALSE,
    number_of_people INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (leader_id) REFERENCES public.users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (faculty_id) REFERENCES public.faculties (faculty_id),
    FOREIGN KEY (program_id) REFERENCES public.programs (program_id)
);

CREATE TABLE IF NOT EXISTS public.members (
    member_id BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL,
    group_id  BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES public.users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES public.groups (group_id) ON DELETE CASCADE
);

INSERT INTO public.faculties (faculty_name) VALUES ('ИИТ');
INSERT INTO public.faculties (faculty_name) VALUES ('Математический факультет');
INSERT INTO public.faculties (faculty_name) VALUES ('Другое');

INSERT INTO public.programs (faculty_id, program_name) VALUES (1, 'Программная инженерия');
INSERT INTO public.programs (faculty_id, program_name) VALUES (1, 'Прикдадная информатика');
INSERT INTO public.programs (faculty_id, program_name) VALUES (2, 'Прикладная математика');
INSERT INTO public.programs (faculty_id, program_name) VALUES (3, 'Другое');

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS public.users;

