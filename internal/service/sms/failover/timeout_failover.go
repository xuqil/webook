package failover

import (
	"context"
	"github.com/xuqil/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service
	idx       int32
	cnt       int32
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		// 这里要切换，新的下标，往后挪了一个
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		return err
	}
}
