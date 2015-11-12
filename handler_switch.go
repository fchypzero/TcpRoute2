package main
import (
	"net"
	"github.com/gamexg/preRead"
	"fmt"
)

// handler 切换器
// 会尝试使用各种类型的 Handler。
type switchHandlerNewer struct {
	handlerNewers []HandlerNewer // 子处理器列表
}

func NewSwitchHandlerNewer() *switchHandlerNewer {
	return &switchHandlerNewer{make([]HandlerNewer, 0, 10)}
}

func (sh *switchHandlerNewer)AppendHandlerNewer(h HandlerNewer) {
	sh.handlerNewers = append(sh.handlerNewers, h)
}

// 尝试创建处理
// 这里会循环尝试创建所有的处理器
func (sh *switchHandlerNewer)New(conn net.Conn) (h Handler, rePre bool, err error) {
	preConn := preread.NewPreConn(conn)
	preConn.NewPre()
	defer preConn.ClosePre()

	// 预先读一次数据
	b := make([]byte, 4096)
	if n, err := preConn.Read(b); err != nil {
		b = b[:n]
	} else {
		b = b[0:0]
	}
	preConn.ResetPreOffset()

	for _, hn := range sh.handlerNewers {
		if h, reset, _ := hn.New(preConn); h != nil {
			if reset {
				preConn.ResetPreOffset()
			}
			return h, false, nil
		} else {
			preConn.ResetPreOffset()
		}
	}

	return nil, true, NoHandleError(fmt.Sprintf("无法识别的协议：%v", b[:10]))
}
