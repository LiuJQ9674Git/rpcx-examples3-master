package serviceprovider


import (
	"os"
	"fmt"
	"time"
	"flag"
	"syscall"
	"context"
	"os/signal"

	//"rpcxstudy/defines"

	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"

)

//rpcx开发服务提供方例子，以zookeeper作为服务注册中心
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

//noinspection ALL
func main(){
	//命令行输入端口
	//命令行参数,启动如./xxxx -port 10001
	port := flag.String("port", ":10001", "port")
	flag.Parse()

	server := server.NewServer()
	plugin := &serverplugin.serverplugin{
		ServiceAddress:   "tcp@localhost:"+*port,      //本机地址
		ZooKeeperServers: []string{"localhost:2181"},  //zookeeper地址
		BasePath:         "/rpcx_test",
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   time.Minute,
	}
	if err := plugin.Start();err != nil {
		server.Close()
		os.Exit(1)
	}else{
		server.Plugins.Add(plugin)
	}

	server.RegisterName("Hello", new(Hello), "")
	go server.Serve("tcp", "0.0.0.0:"+*port)
	StartSignal(server,plugin)
}

//noinspection GoUnresolvedReference
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

