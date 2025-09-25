package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// 数据库的相关操作

// 创建用户数据，并写入数据库
func (m *UserModel) Insert(name, email, password string) error {

	// 通过加密算法，生成哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
			VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// 检查邮箱是否已经存在
		var mySQLError *mysql.MySQLError // mysql返回的错误
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail // model包的全局变量
			}
		}
		return err
	}
	return nil
}

// 输入邮箱，密码，返回用户id
func (m *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte
	// 通过邮箱查询数据库中存储的哈希密码
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// 判断数据库中的哈希密码与用户输入的明文密码是否匹配
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// 返回用户id
	return id, nil

}

// 输入id，返回用户是否存在
func (m *UserModel) Exist(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err

}
