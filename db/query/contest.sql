-- name: CreateContest :one
INSERT INTO sf_contest (
    user_id, 
    subject_id, 
    num_question, 
    time_exam, 
    time_start_exam, 
    state, 
    questions
) 
VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;

-- name: GetContest :one
SELECT 
    id, 
    subject_id, 
    num_question, 
    time_exam, 
    state
FROM 
    sf_contest
WHERE 
    id = $1;

-- name: UpdateContest :one
UPDATE sf_contest
SET 
    num_question = COALESCE($2, num_question),
    time_exam = COALESCE($3, time_exam),
    time_start_exam = COALESCE($4, time_start_exam),
    state = COALESCE($5, state),
    questions = COALESCE($6, questions)
WHERE 
    id = $1
RETURNING 
    id, 
    user_id, 
    subject_id, 
    num_question, 
    time_exam, 
    time_start_exam, 
    state, 
    questions;

-- name: GetContestByState :many 
SELECT 
    id, 
    subject_id, 
    num_question, 
    time_exam, 
    state
FROM 
    sf_contest
WHERE 
    state = $1;