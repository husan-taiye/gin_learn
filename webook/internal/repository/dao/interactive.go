package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GormInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// update a = a + 1
	// 数据库帮完成并发问题
	now := time.Now().UnixMilli()

	return dao.db.Clauses(clause.OnConflict{
		// Columns mysql不用写
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("`read_cnt + 1`"),
			"utime":    time.Now().UnixMilli(),
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	// 作业：就是直接在 LikeCnt 上创建一个索引
	// 1. 而后查询前 100 的，直接就命中索引，这样你前 100 最多 100 次回表
	// SELECT * FROM interactives ORDER BY like_cnt limit 0, 100
	// 还有一种优化思路是
	// SELECT * FROM interactives WHERE like_cnt > 1000 ORDER BY like_cnt limit 0, 100
	// 2. 如果你只需要 biz_id 和 biz_type，你就创建联合索引 <like_cnt, biz_id, biz>
	LikeCnt int64
	Ctime   int64
	Utime   int64
}
