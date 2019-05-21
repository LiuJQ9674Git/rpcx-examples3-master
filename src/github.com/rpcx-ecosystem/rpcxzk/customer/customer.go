package main
import (
	"fmt"
	"time"
	"context"

	"github.com/rpcx-ecosystem/rpcxzk/defines"

	"github.com/smallnest/rpcx/share"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/client"
)

func main() {
	opt := client.Option{
		Retries:        1,
		RPCPath:        share.DefaultRPCPath,
		ConnectTimeout: 10 * time.Second,
		SerializeType:  protocol.JSON,
		CompressType:   protocol.None,
		BackupLatency:  10 * time.Second,
	}

	d := client.NewZookeeperDiscovery("/rpcx_test",
		"Hello",[]string{"localhost:2181"},nil)
	xclient := client.NewXClient("Hello", client.Failtry, client.RoundRobin, d, opt)
	defer xclient.Close()

	in := &defines.CmdIn{
		Param:"123",
	}

	out := &defines.CmdOut{
	}

	for i:=0;;i++{
		in.Param=fmt.Sprintf("%d",i)
		err := xclient.Call(context.Background(), "Test1", in, out)
		if err != nil {
			fmt.Println("failed to call:", err)
		}else{
			fmt.Println("Test1 Result:", out)
		}

		err = xclient.Call(context.Background(), "Test2", in, out)
		if err != nil {
			fmt.Println("failed to call:", err)
		}else{
			fmt.Println("Test2 Result:", out)
		}
	}
}

