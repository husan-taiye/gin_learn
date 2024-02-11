package failover

import (
	"context"
	"gin_learn/webook/internal/service/sms"
	"sync/atomic"
)

type TimeFailoverSMSService struct {
	// 服务商
	svcs []sms.Service
	idx  int32
	// 连续超时的个数
	cnt int32

	// 阈值
	threshold int32
}

func (t TimeFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		// 连续超时
		// 切换新的下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 成功切换了
			// 则重置超时个数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)

	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
		// 不知道什么错误
		// 可以考虑换下一个
		// - 超时错误 可能偶发，尽量再试试
		// - 非超时。直接下一个
		//return err
	}
	return err
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeFailoverSMSService{}
}
