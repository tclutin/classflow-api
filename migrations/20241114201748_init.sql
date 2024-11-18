-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
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

CREATE TABLE IF NOT EXISTS public.type_of_subject (
    type_of_subject_id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.buildings (
    buildings_id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    latitude DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.users (
    user_id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher', 'leader')),
    fullname TEXT,
    telegram TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
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

CREATE TABLE IF NOT EXISTS public.schedule (
    schedule_id BIGINT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    buildings_id BIGINT NOT NULL,
    type_of_subject_id BIGINT NOT NULL,
    subject_name TEXT NOT NULL,
    room TEXT NOT NULL,
    is_even BOOLEAN NOT NULL,
    day_of_week INT NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (group_id) REFERENCES public.groups (group_id),
    FOREIGN KEY (type_of_subject_id) REFERENCES public.type_of_subject(type_of_subject_id),
    FOREIGN KEY (buildings_id) REFERENCES public.buildings(buildings_id)
);

INSERT INTO public.faculties (faculty_name) VALUES ('ИИТ');
INSERT INTO public.faculties (faculty_name) VALUES ('Математический факультет');
INSERT INTO public.faculties (faculty_name) VALUES ('Другое');

INSERT INTO public.programs (faculty_id, program_name) VALUES (1, 'Программная инженерия');
INSERT INTO public.programs (faculty_id, program_name) VALUES (1, 'Прикдадная информатика');
INSERT INTO public.programs (faculty_id, program_name) VALUES (2, 'Прикладная математика');
INSERT INTO public.programs (faculty_id, program_name) VALUES (3, 'Другое');

INSERT INTO public.type_of_subject (name) VALUES ('Лекция');
INSERT INTO public.type_of_subject (name) VALUES ('Практика');
INSERT INTO public.type_of_subject (name) VALUES ('Лабораторная работа');
INSERT INTO public.type_of_subject (name) VALUES ('Другое');

INSERT INTO public.buildings (name, latitude, longitude, address) VALUES ('1 корпус', 33.3, 33.3, 'хуй знает');

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS public.faculties;
DROP TABLE IF EXISTS public.programs;
DROP TABLE IF EXISTS public.type_of_subject;
DROP TABLE IF EXISTS public.buildings;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.groups;
DROP TABLE IF EXISTS public.members;
DROP TABLE IF EXISTS public.schedule;
