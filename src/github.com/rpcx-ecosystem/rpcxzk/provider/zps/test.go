package main
import (
	"os"
	"fmt"
	"time"
	"flag"
	"syscall"
	"context"
	"os/signal"

	"github.com/smallnest/rpcx/server"
	"github.com/rpcx-ecosystem/rpcxzk/defines"
	"github.com/smallnest/rpcx/serverplugin"
)

func main(){
	//命令行输入端口
	//命令行参数,启动如./xxxx -port 10001
	port := flag.String("port", ":10001", "port")
	flag.Parse()

	server := server.NewServer()
	plugin := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@localhost:"+*port,      //本机地址
		ZooKeeperServers: []string{"localhost:2181"},  //zookeeper地址
		BasePath:         "/rpcx_test",
		//Metrics:          metrics.NewRegistry(),NewRegistry
		UpdateInterval:   time.Minute,
	}
	if err := plugin.Start();err != nil {
		server.Close()
		os.Exit(1)
	}else{
		server.Plugins.Add(plugin)
	}

	server.RegisterName("Hello", new(defines.Hello), "")
	go server.Serve("tcp", "0.0.0.0:"+*port)
	StartSignal(server,plugin)
}

func StartSignal(server *server.Server,zk *serverplugin.ZooKeeperRegisterPlugin) {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,syscall.SIGINT)
	for {
		s = <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			server.Shutdown(ctx)
			cancel()
			zk.Stop()
			fmt.Println("stop")
			return
		default:
			return
		}
	}
}

