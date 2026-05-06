package models

import(
	"time"
	"database/sql"
	"errors"
)

type Quiz struct{
	ID int
	Skill string
	Quiz string
	Created time.Time
}

type QuizModel struct{
	DB *sql.DB
}
func (m *QuizModel) Insert(skill string, quiz string ) (int, error){
	  stmt := `INSERT INTO quizzes (skill, quiz, created)
    VALUES(?, ?, UTC_TIMESTAMP())`
    result, err := m.DB.Exec(stmt, skill, quiz)
    if err != nil {
        return 0, err
    }
    // Use the LastInsertId() method on the result to get the ID of our
    // newly inserted record in the snippets table.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    // The ID returned has the type int64, so we convert it to an int type
    // before returning.
    return int(id), nil
}
func(m *QuizModel) Get(id int) (Quiz, error){
	stmt := `SELECT id, skill, quiz, created FROM quizzes
         WHERE id = ?`
	   // holds the result from the database.
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