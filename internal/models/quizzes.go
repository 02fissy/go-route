package models

import (
	"database/sql"
	"errors"
	"time"
)

type Quiz struct{
	ID int
	Skill string
	Quiz string
	Created time.Time
}

type QuizModeler interface{
    Insert(skill string, quiz string ) (int, error)
    Get(id int) (Quiz, error)
    Latest() ([]Quiz, error)
}

type QuizModel struct{
	DB *sql.DB
}
func (m *QuizModel) Insert(skill string, quiz string ) (int, error){
	  stmt := `INSERT INTO quizzes (skill, quiz, created)
    VALUES(?, ?, CURRENT_TIMESTAMP)`
    result, err := m.DB.Exec(stmt, skill, quiz)
    if err != nil {
        return 0, err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
   
    return int(id), nil
}
func(m *QuizModel) Get(id int) (Quiz, error){
	stmt := `SELECT id, skill, quiz, created FROM quizzes
         WHERE id = ?`
	   
    row := m.DB.QueryRow(stmt, id)
	
    var s Quiz
    
    err := row.Scan(&s.ID, &s.Skill, &s.Quiz, &s.Created)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Quiz{}, ErrNoRecord
        } else {
            return Quiz{}, err
        }
    }
    return s, nil
}
func(m *QuizModel) Latest() ([]Quiz, error){
	stmt := `SELECT id, skill, quiz, created FROM quizzes
    ORDER BY id DESC LIMIT 10`
    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

     var quizzes []Quiz
    for rows.Next() {
        var s Quiz
        err := rows.Scan(&s.ID, &s.Skill, &s.Quiz, &s.Created)
		 if err != nil {
            return nil, err
        }
        quizzes = append(quizzes, s)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return quizzes, nil
}