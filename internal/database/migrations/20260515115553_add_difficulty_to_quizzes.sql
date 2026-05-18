-- +goose Up
ALTER TABLE quizzes ADD COLUMN difficulty TEXT NOT NULL DEFAULT 'easy';

-- +goose Down
ALTER TABLE quizzes DROP COLUMN difficulty;
