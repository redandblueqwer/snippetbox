package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) router() http.Handler {
	///////////////////////////////////////////////////////
	// 使用http默认的路由
	// mux := http.NewServeMux()

	// // 静态文件服务器
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)
	/////////////////////////////////////////////////////////////////

	// 使用httprouter的路由

	router := httprouter.New()

	// 自定义router的 NotFound 形式
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// 创建session建立的中间件
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// 创建路由和handler处理
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	// 用户登陆handler函数注册
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// 授权保护
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	// return app.logRequest(secureHeader(mux))
}
