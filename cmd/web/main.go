package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/snippetbox/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	debug          bool
}

func main() {
	// 定义一个新的命令行标志
	addr := flag.String("addr", ":4000", "HTTP network address")
	// 数据库DSN地址
	dsn := flag.String("dsn", "web:0725@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL data source name")
	debug := flag.Bool("debug", false, "Debug mod of application")
	// 解析命令行标志
	flag.Parse()

	// 创建两个消息日志 声明输出format

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 打开数据库
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// 初始化模板
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// 示例化解码器
	formDecoder := form.NewDecoder()

	// 初始化session管理器
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	// 同时设置客户端浏览器Cookie和服务端session的有效期
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true // session仅响应https

	// 初始化app
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		debug:          *debug,
	}

	// tls证书的设置
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// 创建路由
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.router(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// log.Print("Starting server on :4000")
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
