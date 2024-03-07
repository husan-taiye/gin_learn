package article

import (
	"context"
	"gin_learn/webook/internal/domain"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDao struct {
	//
	client *mongo.Client
	// 代表 webook 的
	db *mongo.Database
	// 代表制作库
	col *mongo.Collection
	//代表线上库
	liveCol *mongo.Collection

	node *snowflake.Node
}

func (m *MongoDBDao) Insert(ctx context.Context, art Article) (int64, error) {
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	if err != nil {
		return 0, err
	}
	// 没有自增逐渐
	// GUID Globally Unique ID
	return id, nil
}

func (m *MongoDBDao) UpdateById(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDao) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	// upsert语意
	now := time.Now().UnixMilli()
	update := bson.E{"$set", PublishArticle{art}}
	upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	filter := bson.M{"id": art.Id}

	_, err = m.liveCol.UpdateOne(ctx, filter, bson.D{update, upsert},
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBDao) Upsert(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDao) SyncStatus(ctx context.Context, art domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDao{
		col:     db.Collection("articles"),
		liveCol: db.Collection("publish_articles"),
		node:    node,
	}
}
