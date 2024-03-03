package startup

import (
	"gin_learn/webook/internal/service/sms"
	"gin_learn/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
