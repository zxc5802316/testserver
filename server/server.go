package server

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出
// Server http 服务  返回 errgroup 第一次产生错误
func Server(addr string) error {
	g, ctx := errgroup.WithContext(context.Background())
	srv := http.Server{Addr: addr}
	// g1
	// 开启服务端监听
	// 当 g1 产出 error ,errgroup 会调用 g.cancel()
	// g2 ，g3 也会退出
	g.Go(func() error {
		return srv.ListenAndServe()
	})

	//g2
	// select 阻塞 ，等待 ctx.done()
	// 有 ctx.done() 产生  g3 会推出
	// Shutdown 关闭 服务器
	// g2 ，g1 退出
	g.Go(func() error {

		select {
		case <-ctx.Done():
		}
		tctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*3)
		defer cancelFunc()
		return srv.Shutdown(tctx)
	})

	//g3
	// 接收到信号 返回 error g3 退出
	// g2 退出 关闭服务器
	// g1 退出
	g.Go(func() error {
		stop := make(chan os.Signal,1)
		//syscall.SIGINT == ctrl + c ?
		signal.Notify(stop, syscall.SIGINT)
		select {
		case sig := <-stop:
			close(stop)
			// 也可以 使用 switch 根据signal 不同 做不同的处理
			return errors.New(fmt.Sprintf("get os.Signal %v", sig))
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	return g.Wait()
}
