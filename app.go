package goweb

import (
	"fmt"
	"github.com/liuyongshuai/thingo/controller"
	"github.com/liuyongshuai/thingo/router"
	"html/template"
	"net/http"
)

//新建一个APP
func NewThingoApp() *ThingoApp {
	app := &ThingoApp{
		Handlers: NewThingoHandler(),
	}
	return app
}

//APP结构体
type ThingoApp struct {
	Handlers *ThingoHandler //处理句柄
}

//开始运行
func (app *ThingoApp) Run() {
	app.Handlers.Tpl.SetRootPathDir(app.Handlers.TplDir).SetTplExt(app.Handlers.TplExt)
	err := http.ListenAndServe(":"+app.Handlers.Port, app.Handlers)
	if err != nil {
		fmt.Println(err)
	}
}

//设置监听端口
func (app *ThingoApp) SetPort(port string) *ThingoApp {
	app.Handlers.Port = port
	return app
}

//设置错误信息提示
func (app *ThingoApp) SetErrController(c controller.ThingoControllerInterface) *ThingoApp {
	c = c.(controller.ThingoControllerInterface)
	app.Handlers.SetErrController(c)
	return app
}

//设置POST最大内存
func (app *ThingoApp) SetMaxMemory(n int64) *ThingoApp {
	app.Handlers.SetMaxMemory(n)
	return app
}

//设置模板路径
func (app *ThingoApp) SetTplDir(dir string) *ThingoApp {
	app.Handlers.SetTplDir(dir)
	return app
}

//设置模板扩展名称
func (app *ThingoApp) SetTplExt(ext string) *ThingoApp {
	app.Handlers.SetTplExt(ext)
	return app
}

//设置给模板的公共参数
func (app *ThingoApp) SetTplCommonData(data map[interface{}]interface{}) *ThingoApp {
	app.Handlers.SetTplCommonData(data)
	return app
}

//设置给模板的公共参数
func (app *ThingoApp) AddTplCommonData(k interface{}, v interface{}) *ThingoApp {
	app.Handlers.AddTplCommonData(k, v)
	return app
}

//添加一个插件
func (app *ThingoApp) AddHooks(when int, hk HooksFunc) *ThingoApp {
	app.Handlers.AddHooks(when, hk)
	return app
}

//添加模板函数
func (app *ThingoApp) AddTplFuncMap(fm template.FuncMap) *ThingoApp {
	app.Handlers.Tpl.AddTplFuncs(fm)
	return app
}

//添加模板函数
func (app *ThingoApp) AddTplFunc(name string, fn interface{}) *ThingoApp {
	app.Handlers.Tpl.AddTplFunc(name, fn)
	return app
}

//添加一个路由
func (app *ThingoApp) AddRouter(r *router.ThingoRouterItem) *ThingoApp {
	app.Handlers.AddRouter(r)
	return app
}

//批量添加路由
func (app *ThingoApp) AddRouters(rs ...*router.ThingoRouterItem) *ThingoApp {
	app.Handlers.AddRouters(rs...)
	return app
}

//设置发生错误时的处理函数
func (app *ThingoApp) SetRecoverFunc(fn RecoverFunc) *ThingoApp {
	app.Handlers.SetRecoverFunc(fn)
	return app
}
