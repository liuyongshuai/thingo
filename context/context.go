package context

import (
	"fmt"
	"github.com/liuyongshuai/negoutils/snowflake"
	"net/http"
)

var snowFlake *snowflake.IdGenerator

func init() {
	snowFlake, _ = snowflake.
		NewIDGenerator().
		SetTimeBitSize(50).
		SetSequenceBitSize(8).
		SetWorkerIdBitSize(5).
		SetWorkerId(1).
		Init()
}

//返回请求上下文
func NewThingoContext() *ThingoContext {
	RuntofuCtx := &ThingoContext{
		Input:  NewThingoInput(),
		Output: NewThingoOutput(),
	}
	RuntofuCtx.Output.Context = RuntofuCtx
	RuntofuCtx.Input.Context = RuntofuCtx
	return RuntofuCtx
}

//上下文的定义
type ThingoContext struct {
	Input          *ThingoInput        //收到的请求里相关信息，包括参数、方法、上传文件等
	Output         *ThingoOuput        //要发送给端的暂存用的数据
	Request        *http.Request       //请求原始对象指针
	ResponseWriter http.ResponseWriter //响应原始对象
	UniqueKey      string              //本次请求的唯一标识符
}

//重置本次请求的上下文
func (RuntofuCtx *ThingoContext) Reset(rw *http.ResponseWriter, r *http.Request) {
	RuntofuCtx.Request = r
	RuntofuCtx.ResponseWriter = *rw
	RuntofuCtx.Input.Reset(RuntofuCtx)
	RuntofuCtx.Output.Reset(RuntofuCtx)
	var nextId int64 = 0
	var err error
	for {
		nextId, err = snowFlake.NextId()
		if err == nil {
			break
		}
	}
	RuntofuCtx.UniqueKey = fmt.Sprintf("%x", nextId)
}

//跳转，状态码可选，默认301
func (RuntofuCtx *ThingoContext) Redirect(locationUrl string, status ...int) {
	code := http.StatusTemporaryRedirect
	if len(status) > 0 {
		code = status[0]
	}
	http.Redirect(RuntofuCtx.ResponseWriter, RuntofuCtx.Request, locationUrl, code)
}

//刷新返回数据
func (RuntofuCtx *ThingoContext) Flush() {
	if f, ok := RuntofuCtx.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

//当客户端取消请求或连接断开时用
func (RuntofuCtx *ThingoContext) CloseNotify() <-chan bool {
	if cn, ok := RuntofuCtx.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
