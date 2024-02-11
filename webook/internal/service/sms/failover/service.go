package failover

import (
	"context"
	"errors"
	"gin_learn/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs ...sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		log.Panicln(err)
	}
	return errors.New("全部服务商都失败了")
}

// SendV2 相对于Send的改进是
// 每次不用从第一个svc开始，相对的负载均衡一点
// 区别了错误 case context.DeadlineExceeded, context.Canceled 跟用户体验密切相关
func (f FailoverSMSService) SendV2(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+uint64(length); i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 输出日志
			log.Println("")
		}
	}
	return errors.New("全部服务商都失败了")
}
