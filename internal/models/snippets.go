package models

import (
	"database/sql"
	"errors"
	"time"
)

// 定义结构体接受数据

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// 数据库连接池
type SnippetModel struct {
	DB *sql.DB
}

// 插入数据
// 输入: title\content\expires time
// 输出: id\error
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	//  定义sql查询语句
	stmt := `INSERT INTO snippets (title,content,created,expires) 
	         VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY) )`

	// DB.Exec() 用于插入或者删除，不用返回查询结果
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil

}

// 依据id查询数据
// 输入: id
// 输出: Snippet数据结构\erro
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	//  定义sql查询语句
	stmt := `SELECT id,title,content,created,expires FROM snippets
	         WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// DB.QueryRow() 返回单行结果
	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// 查询10条最新的数据
// 输出: Snippet数据结构的list
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id,title,content,created, expires FROM 
	         snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// DB.Query() 返回多行结果
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippet := []*Snippet{}

	for rows.Next() {
		// 临时变量
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippet = append(snippet, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippet, nil
}
