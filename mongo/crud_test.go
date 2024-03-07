package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	// 控制初始化超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println("执行命令: ", startedEvent.Command)
			//ctx = context.WithValue(ctx, "start_time", time.Now())
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	col := mdb.Collection("articles")

	// 插入
	res, err := col.InsertOne(ctx, Article{
		Id:      1,
		Title:   "mongo first 标题",
		Content: "内容",
		Ctime:   time.Now().UnixMilli(),
		Utime:   time.Now().UnixMilli(),
	})
	assert.NoError(t, err)
	// 文档id， mongodb 中的 _id 字段
	fmt.Printf("id %s", res.InsertedID)

	// 查找
	filter := bson.D{bson.E{Key: "id", Value: 1}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("%v \n", art)
	art = Article{}
	err = col.FindOne(ctx, Article{Id: 1}).Decode(&art)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("没有数据")
	}
	assert.NoError(t, err)
	fmt.Printf("%v \n", art)

	// 修改
	// ？？？？？todo
	// 正确的
	//sets := bson.D{bson.E{Key: "$set", Value: bson.M{"title": "新的标题"}}}
	// 还是不对
	//sets := bson.D{bson.E{Key: "$set", Value: bson.E{"title", "新的标题"}}}
	sets := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "title", Value: "新的标题"}}}
	upsetRes, err := col.UpdateMany(ctx, filter, sets)
	assert.NoError(t, err)
	fmt.Println("affected", upsetRes.ModifiedCount)
	upsetRes, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{Key: "$set", Value: Article{Title: "新的标题2", AuthorId: 12345}},
	})
	assert.NoError(t, err)
	fmt.Println("affected", upsetRes.ModifiedCount)

	// 删除
	//delRes, err := col.DeleteMany(ctx, filter)
	//assert.NoError(t, err)
	//fmt.Println("affected", delRes.DeletedCount)

	// or 查询
	or := bson.A{bson.D{bson.E{"id", 1}},
		bson.D{bson.E{"id", 2}}}
	orRes, err := col.Find(ctx, bson.D{bson.E{"$or", or}})
	assert.NoError(t, err)
	var ars []Article
	err = orRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("%v \n", ars)

	// and 查询
	and := bson.A{bson.D{bson.E{"id", 1}},
		bson.D{bson.E{"ctime", 1709696435983}}}
	andRes, err := col.Find(ctx, bson.D{bson.E{"$and", and}})
	assert.NoError(t, err)
	ars = []Article{}
	err = andRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("%v \n", ars)

	// in 查询
	in := bson.D{bson.E{"id", bson.D{bson.E{"$in", []any{1, 2}}}}}
	inRes, err := col.Find(ctx, in)
	ars = []Article{}
	err = inRes.All(ctx, &ars)
	assert.NoError(t, err)
	fmt.Printf("%v \n", ars)

}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   int8   `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
