package goweb

import (
	"github.com/liuyongshuai/thingo/context"
	"github.com/liuyongshuai/thingo/controller"
	"github.com/liuyongshuai/thingo/router"
	"net/http"
	"reflect"
	"sync"
)

//插件类型
const (
	HooksBeforeRun = iota + 1
	HooksAfterRun
)

//插件函数
type HooksFunc func(ctx *context.ThingoContext)

//panic后的处理函数
type RecoverFunc func(*context.ThingoContext)

//控制器注册器
type ThingoHandler struct {
	Hooks         map[int][]HooksFunc                  //所有的插件列表
	Router        *router.ThingoRouterList             //路由列表
	RecoverFunc   RecoverFunc                          //panic后的处理函数
	pool          sync.Pool                            //context上下文池
	Tpl           *controller.TplBuilder               //模板对象类型
	TplExt        string                               //模板的扩展后缀，默认“tpl”
	TplDir        string                               //模板的根目录，默认“./tpl/”
	TplCommonData map[interface{}]interface{}          //模板的公共参数
	Port          string                               //监听的端口
	MaxMemory     int64                                //POST时的最大内存
	ErrController controller.ThingoControllerInterface //当匹配不上时的错误信息页面
}

func NewThingoHandler() *ThingoHandler {
	cr := &ThingoHandler{
		Hooks:         make(map[int][]HooksFunc),
		MaxMemory:     64 << 20,
		Tpl:           controller.NewTplBuilder(),
		TplExt:        "tpl",
		TplDir:        "./tpl",
		Router:        router.NewThingoRouterList(),
		TplCommonData: make(map[interface{}]interface{}),
	}
	cr.Hooks[HooksBeforeRun] = []HooksFunc{}
	cr.Hooks[HooksAfterRun] = []HooksFunc{}
	cr.pool.New = func() interface{} {
		return context.NewThingoContext()
	}
	return cr
}

//设置监听端口
func (cr *ThingoHandler) SetPort(port string) {
	cr.Port = port
}

//添加一个插件
func (cr *ThingoHandler) AddHooks(when int, hk HooksFunc) {
	plist, ok := cr.Hooks[when]
	if !ok {
		return
	}
	plist = append(plist, hk)
	cr.Hooks[when] = plist
}

//设置模板路径
func (cr *ThingoHandler) SetTplDir(dir string) {
	cr.TplDir = dir
}

//设置模板扩展名称
func (cr *ThingoHandler) SetTplExt(ext string) {
	cr.TplExt = ext
}

//设置给模板的公共参数
func (cr *ThingoHandler) SetTplCommonData(data map[interface{}]interface{}) {
	for k, v := range data {
		cr.TplCommonData[k] = v
	}
}

//设置给模板的公共参数
func (cr *ThingoHandler) AddTplCommonData(k interface{}, v interface{}) {
	cr.TplCommonData[k] = v
}

//添加一个路由
func (cr *ThingoHandler) AddRouter(r *router.ThingoRouterItem) {
	cr.Router.AddRouter(r)
}

//批量添加路由
func (cr *ThingoHandler) AddRouters(rs ...*router.ThingoRouterItem) {
	cr.Router.AddRouters(rs...)
}

//设置发生错误时的处理函数
func (cr *ThingoHandler) SetRecoverFunc(fn RecoverFunc) {
	cr.RecoverFunc = fn
}

//设置POST最大内存
func (cr *ThingoHandler) SetMaxMemory(n int64) {
	cr.MaxMemory = n
}

//设置错误信息提示
func (cr *ThingoHandler) SetErrController(c controller.ThingoControllerInterface) {
	cr.ErrController = c
}

//执行 http.Handler 接口
func (cr *ThingoHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//从池子里提取上下文实例
	ctx := cr.pool.Get().(*context.ThingoContext)
	if ctx == nil {
		panic("get context failed")
	}
	ctx.Reset(&rw, r)
	defer cr.pool.Put(ctx)

	//异常恢复函数设置
	if cr.RecoverFunc != nil {
		defer cr.RecoverFunc(ctx)
	}

	//解析表单提交上来的参数
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		ctx.Input.ParseFormOrMulitForm(cr.MaxMemory)
	}

	//控制层的类
	var controllerIface controller.ThingoControllerInterface
	var ok bool

	//开始匹配路由
	routerItem := cr.Router.Match(ctx, r)
	if routerItem == nil {
		ctx.Output.SetStatus(http.StatusNotFound)
		reflectVal := reflect.ValueOf(cr.ErrController)
		ct := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(ct)
		controllerIface, ok = vc.Interface().(controller.ThingoControllerInterface)
		if !ok {
			panic("invalid controller")
		}
	} else {
		//实例化一个控制层对象
		vc := reflect.New(routerItem.ControllerType)
		controllerIface, ok = vc.Interface().(controller.ThingoControllerInterface)
		if !ok {
			panic("invalid controller")
		}
	}

	//执行Before插件
	for _, hk := range cr.Hooks[HooksBeforeRun] {
		hk(ctx)
	}

	if ctx.Output.Started == true {
		return
	}

	//执行控制层
	controllerIface.Init(ctx, controllerIface, cr.Tpl, cr.TplCommonData)
	err := controllerIface.Prepare()
	if err == nil {
		controllerIface.Run()
	}
	//执行After插件
	for _, hk := range cr.Hooks[HooksAfterRun] {
		hk(ctx)
	}
	controllerIface.Finish()
	//刷新输出
	ctx.Output.Send()
}
