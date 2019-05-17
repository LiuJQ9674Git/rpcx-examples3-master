package defines

import (
	"fmt"
	"context"
)

//输入参数定义

type CmdIn struct {
	Param string
}
//输出参数定义

type CmdOut struct {
	Result int
	Info string
}
//定义一个Hello服务
type Hello struct {
	Param string
}
//定义Test1方法
func (this *Hello) Test1(ctx context.Context, in *CmdIn, out *CmdOut) error {
	fmt.Println("1 recv:",in)
	out.Result=12345678
	out.Info="set value"+in.Param
	return nil
}
//定义Test2方法
func (this *Hello) Test2(ctx context.Context, in *CmdIn, out *CmdOut) error {
	fmt.Println("2 recv:",in)
	out.Result=0
	out.Info=in.Param
	return nil
}

