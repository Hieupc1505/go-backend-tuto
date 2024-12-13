-- Tạo ký hóa contest_state
CREATE TYPE contest_state AS ENUM ('IDLE', 'RUNNING', 'FINISHED');

-- Tạo bảng sf_contest
CREATE TABLE sf_contest (
    id bigserial PRIMARY KEY NOT NULL,
    user_id bigint NOT NULL,
    subject_id bigint NOT NULL,
    num_question int NOT NULL,
    time_exam int NOT NULL,
    time_start_exam bigint NOT NULL,
    state contest_state NOT NULL DEFAULT 'IDLE', -- PostgreSQL không có ENUM mặc định, dùng VARCHAR thay thế
    questions text NOT NULL, -- Thay thế MEDIUMTEXT bằng TEXT
    created_time timestamptz NOT NULL DEFAULT (now()),
    updated_time timestamptz NOT NULL DEFAULT (now())
);

-- Tạo index cho cột user_id
CREATE INDEX idx_sf_contest_user_id ON sf_contest (user_id);

-- Thêm khóa ngoại user_id tham chiếu đến bảng sf_user
ALTER TABLE sf_contest
ADD CONSTRAINT sf_contest_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);