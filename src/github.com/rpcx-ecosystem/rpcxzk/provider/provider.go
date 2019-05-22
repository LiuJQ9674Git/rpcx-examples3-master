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
	"github.com/rcrowley/go-metrics"
)

func main(){
	//命令行输入端口
	//命令行参数,启动如./xxxx -port 10001
	port := flag.String("port", "10001", "port")
	flag.Parse()
	//1.新建服务
	server := server.NewServer()
	//注册中心为ZK,把服务写入注册中心
	plugin := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@localhost:"+*port,      //服务地址
		ZooKeeperServers: []string{"localhost:2181"},  //zookeeper地址
		BasePath:         "/rpcx_test",//访问路径
		Metrics:          metrics.NewRegistry(),//NewRegistry
		UpdateInterval:   time.Minute,
	}
	//2.zk start
	if err := plugin.Start();err != nil {
		server.Close()
		os.Exit(1)
	}else{
		server.Plugins.Add(plugin)
	}

	//3. 注册服务
	server.RegisterName("Hello", new(defines.Hello), "")
	//4. 启动服务TCP
	go server.Serve("tcp", "0.0.0.0:"+*port)

	//5
	StartSignal(server,plugin)
	//
}

func StartSignal(server *server.Server,zk *serverplugin.ZooKeeperRegisterPlugin) {
	//
	var (
		c chan os.Signal
		s os.Signal
	)
	//一个缓冲的通道
	c = make(chan os.Signal, 1)
	//发送消息
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

