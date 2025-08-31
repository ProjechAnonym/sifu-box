package nodes

import (
	"net/http"
	"sifu-box/ent"
	"sync"

	"go.uber.org/zap"
)

// Merge 并发地从多个 provider 获取数据并更新到数据库中
// 参数:
//   - providers: 要处理的 provider 列表
//   - ent_client: 用于操作数据库的 ent 客户端
//   - logger: 用于记录日志的 zap logger 实例
//
// 返回值:
//   - []error: 在处理过程中发生的错误列表
func Fetch(providers []*ent.Provider, ent_client *ent.Client, logger *zap.Logger) []error {
	var jobs sync.WaitGroup
	var err_chan = make(chan error, 5)
	var count_chan = make(chan int, 5)
	var errors []error
	client := http.DefaultClient

	// 启动一个 goroutine 来收集完成的任务数量和错误信息
	jobs.Add(1)
	go func() {
		defer func() {
			jobs.Done()
			var ok bool
			if _, ok = <-count_chan; ok {
				close(count_chan)
			}
			if _, ok = <-err_chan; ok {
				close(err_chan)
			}
		}()
		sum := 0
		for {
			if sum == len(providers) {
				return
			}
			select {
			case count, ok := <-count_chan:
				if !ok {
					return
				}
				sum += count
			case err, ok := <-err_chan:
				if !ok {
					return
				}
				errors = append(errors, err)
			}
		}
	}()

	// 遍历所有 providers, 并为每个 provider 启动一个 goroutine 进行处理
	for _, provider := range providers {
		jobs.Add(1)
		go func() {
			defer func() {
				jobs.Done()
				count_chan <- 1
			}()

			// 根据 provider 是否是远程的, 选择不同的获取方式
			if provider.Remote {
				outbounds, err := fetchFromRemote(provider.Name, provider.Path, client, logger)
				if err != nil {
					err_chan <- err
					return
				}
				if err := updateNodes(provider.Name, provider.UUID, outbounds, ent_client); err != nil {
					err_chan <- err
					return
				}

			} else {
				outbounds, err := fetchFromLocal(provider.Name, provider.Path, logger)
				if err != nil {
					err_chan <- err
					return
				}
				if err := updateNodes(provider.Name, provider.UUID, outbounds, ent_client); err != nil {
					err_chan <- err
					return
				}

			}
		}()
	}

	// 等待所有的 goroutine 完成
	jobs.Wait()
	return errors
}
