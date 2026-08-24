package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Zapi-web/url-shortener/internal/domain"
)

type cache interface {
	Set(ctx context.Context, shortURL string, longURL string) error
	Get(ctx context.Context, shortURL string) (string, error)
}

type AsyncCache struct {
	next     cache
	tasks    chan domain.CacheTask
	wg       sync.WaitGroup
	isClosed bool
	mu       sync.RWMutex
	stopOnce sync.Once
}

func NewAsync(baseCache cache, bufferSize int) *AsyncCache {
	return &AsyncCache{next: baseCache, tasks: make(chan domain.CacheTask, bufferSize)}
}

func (a *AsyncCache) StartAsync(workersInt int) {
	a.wg.Add(workersInt)
	for i := 0; i < workersInt; i++ {
		go func() {
			defer a.wg.Done()
			a.work()
		}()
	}
}

func (a *AsyncCache) work() {
	for task := range a.tasks {
		a.doTask(task)
	}
}

func (a *AsyncCache) doTask(task domain.CacheTask) {
	cacheCtx, cancel := context.WithTimeout(
		context.WithoutCancel(task.RequestContext),
		300*time.Millisecond,
	)
	defer cancel()

	if err := a.next.Set(cacheCtx, task.Key, task.Value); err != nil {
		slog.ErrorContext(task.RequestContext, "failed to set key-value in cache", "key", task.Key, "err", err)
	}
}

func (a *AsyncCache) Set(ctx context.Context, shortURL string, longURL string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.isClosed {
		return domain.ErrQueueClosed
	}

	select {
	case a.tasks <- domain.CacheTask{RequestContext: ctx, Key: shortURL, Value: longURL}:
		return nil
	default:
		return domain.ErrQueueFull
	}
}

func (a *AsyncCache) Get(ctx context.Context, shortURL string) (string, error) {
	return a.next.Get(ctx, shortURL)
}

func (a *AsyncCache) Stop(ctx context.Context) error {
	a.stopOnce.Do(func() {
		a.mu.Lock()
		a.isClosed = true
		close(a.tasks)
		a.mu.Unlock()
	})

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
