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
	var errChan = make(chan error, 5)
	var countChan = make(chan int, 5)
	var errors []error
	client := http.DefaultClient

	// 启动一个 goroutine 来收集完成的任务数量和错误信息
	jobs.Add(1)
	go func() {
		defer func() {
			jobs.Done()
			var ok bool
			if _, ok = <-countChan; ok {
				close(countChan)
			}
			if _, ok = <-errChan; ok {
				close(errChan)
			}
		}()
		sum := 0
		for {
			if sum == len(providers) {
				return
			}
			select {
			case count, ok := <-countChan:
				if !ok {
					return
				}
				sum += count
			case err, ok := <-errChan:
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
				countChan <- 1
			}()

			// 根据 provider 是否是远程的, 选择不同的获取方式
			if provider.Remote {
				outbounds, err := fetchFromRemote(provider.Name, provider.Path, client, logger)
				if err != nil {
					errChan <- err
					return
				}
				if err := updateNodes(provider.Name, provider.UUID, outbounds, ent_client); err != nil {
					errChan <- err
					return
				}

			} else {
				outbounds, err := fetchFromLocal(provider.Name, provider.Path, logger)
				if err != nil {
					errChan <- err
					return
				}
				if err := updateNodes(provider.Name, provider.UUID, outbounds, ent_client); err != nil {
					errChan <- err
					return
				}

			}
		}()
	}

	// 等待所有的 goroutine 完成
	jobs.Wait()
	return errors
}
