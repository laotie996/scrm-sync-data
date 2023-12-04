package console

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SyncDataConsole struct {
}

type StreamObject struct {
	Id                       *WatchId `bson:"_id"`
	OperationType            string
	FullDocument             map[string]interface{}
	FullDocumentBeforeChange map[string]interface{}
	Ns                       NS
	UpdateDescription        map[string]interface{}
	DocumentKey              map[string]interface{}
}
type NS struct {
	Database   string `bson:"db"`
	Collection string `bson:"coll"`
}
type WatchId struct {
	Data string `bson:"_data"`
}

const (
	OperationTypeInsert  = "insert"
	OperationTypeDelete  = "delete"
	OperationTypeUpdate  = "update"
	OperationTypeReplace = "replace"
)

var resumeToken bson.Raw

func (c *SyncDataConsole) insertAccount(data map[string]interface{}) {
	sql := `insert into account(_id, root_cst_id, cst_id,cst_name,task_type,worker_id,worker_account,worker_user_name,
                                    worker_phone,user_account,user_name,user_phone,user_first_name,user_last_name,user_status,user_status_code,
                                    user_status_remark,source_type,source_uuid,request_uuid,pull_status,pull_at,created_at,updated_at) 
    values ('%s', %d, %d, '%s', %d , %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', %d, '%s','%s', %d, %d, %d, %d);`
	sql = fmt.Sprintf(sql, data["_id"].(primitive.ObjectID).Hex(),
		this.Format.ToInt64(data["root_cst_id"]), this.Format.ToInt64(data["cst_id"]),
		this.Format.ToString(data["cst_name"]), this.Format.ToInt32(data["task_type"]),
		this.Format.ToInt64(data["worker_id"]), this.Format.ToString(data["worker_account"]),
		this.Format.ToString(data["worker_user_name"]), this.Format.ToString(data["worker_phone"]),
		this.Format.ToString(data["user_account"]), this.Format.ToString(data["user_name"]),
		this.Format.ToString(data["user_phone"]), this.Format.ToString(data["user_first_name"]),
		this.Format.ToString(data["user_last_name"]), this.Format.ToInt32(data["user_status"]),
		this.Format.ToInt32(data["user_status_code"]), this.Format.ToString(data["user_status_remark"]),
		this.Format.ToInt32(data["source_type"]), this.Format.ToString(data["source_uuid"]),
		this.Format.ToString(data["request_uuid"]), this.Format.ToInt32(data["pull_status"]),
		this.Format.ToInt64(data["pull_at"]), this.Format.ToInt64(data["created_at"]),
		this.Format.ToInt64(data["updated_at"]),
	)
	this.Logger.Debugf("%+v", sql)
	err := this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("insert failed,sql:%+v", sql)
		this.Logger.Errorf("insert into account failed,%+v", err.Error())
	}
	return
}

func (c *SyncDataConsole) deleteAccount(data map[string]interface{}) {
	if data == nil {
		return
	}
	sql := fmt.Sprintf(`alter table account delete where _id='%s'`, data["_id"].(primitive.ObjectID).Hex())
	this.Logger.Debugf("%+v", sql)
	err := this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("delete failed,sql:%+v", sql)
		this.Logger.Errorf("delete from account failed,%+v", err.Error())
	}
	return
}

func (c *SyncDataConsole) updateAccount(data map[string]interface{}) {
	if data == nil {
		return
	}
	sql := `alter table account update root_cst_id=%d,cst_id=%d,cst_name='%s',task_type=%d,worker_id=%d,worker_account='%s', worker_user_name='%s',worker_phone='%s',user_account='%s',user_name='%s',user_phone='%s',user_first_name='%s',user_last_name='%s',user_status=%d,user_status_code=%d,user_status_remark='%s',source_type=%d,source_uuid='%s',request_uuid='%s',pull_status=%d,pull_at=%d,updated_at=%d where _id='%s';`
	sql = fmt.Sprintf(sql, this.Format.ToInt64(data["root_cst_id"]),
		this.Format.ToInt64(data["cst_id"]), this.Format.ToString(data["cst_name"]),
		this.Format.ToInt32(data["task_type"]), this.Format.ToInt64(data["worker_id"]),
		this.Format.ToString(data["worker_account"]), this.Format.ToString(data["worker_user_name"]),
		this.Format.ToString(data["worker_phone"]), this.Format.ToString(data["user_account"]),
		this.Format.ToString(data["user_name"]), this.Format.ToString(data["user_phone"]),
		this.Format.ToString(data["user_first_name"]), this.Format.ToString(data["user_last_name"]),
		this.Format.ToInt32(data["user_status"]), this.Format.ToInt32(data["user_status_code"]),
		this.Format.ToString(data["user_status_remark"]), this.Format.ToInt32(data["source_type"]),
		this.Format.ToString(data["source_uuid"]), this.Format.ToString(data["request_uuid"]),
		this.Format.ToInt32(data["pull_status"]), this.Format.ToInt64(data["pull_at"]),
		this.Format.ToInt64(data["updated_at"]),
		data["_id"].(primitive.ObjectID).Hex(),
	)
	this.Logger.Debugf("%+v", sql)
	err := this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("update failed,sql:%+v", sql)
		this.Logger.Errorf("update from account failed,%+v", err.Error())
	}
}

func (c *SyncDataConsole) insertMessage(data map[string]interface{}) (err error) {
	sql := `insert into tg_private_chat_message(_id, conv_id, tg_id, to_user_id, send_side, content, translated_content, msg_id,
                                    msg_type, msg_time,msg_timestamp,file_url,file_name,reply_to,created_at) values ('%s', %d, '%s','%s',%d,'%s','%s', %d, %d, '%s', %d, '%s','%s', %d, %d);`
	sql = fmt.Sprintf(sql,
		data["_id"].(primitive.ObjectID).Hex(), this.Format.ToInt64(data["conv_id"]),
		this.Format.ToString(data["tg_id"]), this.Format.ToString(data["to_user_id"]),
		this.Format.ToInt32(data["send_side"]), this.Format.ToString(data["content"]),
		this.Format.ToString(data["translated_content"]), this.Format.ToInt64(data["msg_id"]),
		this.Format.ToInt32(data["msg_type"]), this.Format.ToString(data["msg_time"].(primitive.DateTime).Time().Format(time.DateTime)),
		this.Format.ToInt64(data["msg_timestamp"]), this.Format.ToString(data["file_url"]),
		this.Format.ToString(data["file_name"]), this.Format.ToInt64(data["reply_to"]),
		this.Format.ToInt64(data["created_at"]))
	this.Logger.Debugf("%+v", sql)
	err = this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("insert failed,sql:%+v", sql)
		this.Logger.Errorf("insert into tg_private_chat_message failed,%+v", err.Error())
	}
	return
}

func (c *SyncDataConsole) deleteMessage(data map[string]interface{}) {
	if data == nil {
		return
	}
	sql := fmt.Sprintf(`alter table tg_private_chat_message delete where _id='%s'`, data["_id"].(primitive.ObjectID).Hex())
	this.Logger.Debugf("%+v", sql)
	err := this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("delete failed,sql:%+v", sql)
		this.Logger.Errorf("delete from tg_private_chat_message failed,%+v", err.Error())
	}
	return
}

func (c *SyncDataConsole) updateMessage(data map[string]interface{}) {
	if data == nil {
		return
	}
	sql := `alter table tg_private_chat_message update conv_id=%d,tg_id='%s',to_user_id='%s',
                                   send_side=%d,content='%s',translated_content='%s', msg_id=%d,msg_type=%d,
                                   msg_time='%s',msg_timestamp=%d,file_url='%s',file_name='%s',reply_to=%d where _id='%s'`
	sql = fmt.Sprintf(sql, this.Format.ToInt64(data["conv_id"]),
		this.Format.ToString(data["tg_id"]), this.Format.ToString(data["to_user_id"]),
		this.Format.ToInt32(data["send_side"]), this.Format.ToString(data["content"]),
		this.Format.ToString(data["translated_content"]), this.Format.ToInt64(data["msg_id"]),
		this.Format.ToInt32(data["msg_type"]), this.Format.ToString(data["msg_time"].(primitive.DateTime).Time().Format(time.DateTime)),
		this.Format.ToInt64(data["msg_timestamp"]), this.Format.ToString(data["file_url"]),
		this.Format.ToString(data["file_name"]), this.Format.ToInt64(data["reply_to"]),
		data["_id"].(primitive.ObjectID).Hex())
	this.Logger.Debugf("%+v", sql)
	err := this.Clickhouse.Exec(context.TODO(), sql)
	if err != nil {
		this.Logger.Infof("delete failed,sql:%+v", sql)
		this.Logger.Errorf("delete from tg_private_chat_message failed,%+v", err.Error())
	}
}

func (c *SyncDataConsole) syncAccountData() {
	scrmTaskListenDBSource := this.Config.MongoDBServiceConfig.DBMap["scrm-task"].ListenDBSource
	this.Logger.Infof(scrmTaskListenDBSource)
	clientOptions := options.Client().ApplyURI(scrmTaskListenDBSource).SetMaxPoolSize(1000)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		this.Logger.Errorf("connect to mongodb failed,%+v", err.Error())
		return
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		this.Logger.Errorf("ping mongodb failed,%+v", err.Error())
		return
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match",
			bson.M{"operationType": bson.M{"$in": bson.A{"insert", "delete", "replace", "update"}}},
		}},
	}
	getAccountResumeTokenCmd := this.Redis.Get(context.TODO(), "account_resume_token")
	resumeTokenJson := ""
	if getAccountResumeTokenCmd.Err() != nil {
		if getAccountResumeTokenCmd.Err() != redis.Nil {
			this.Logger.Errorf("get message resume token failed,%+v", getAccountResumeTokenCmd.Err())
			return
		}
	} else {
		resumeTokenJson = getAccountResumeTokenCmd.Val()
	}
	var opt *options.ChangeStreamOptions
	if len(resumeTokenJson) > 0 {
		err := jsoniter.UnmarshalFromString(resumeTokenJson, &resumeToken)
		if err != nil {
			this.Logger.Errorf("%+v", err.Error())
			return
		}
		opt = options.ChangeStream().SetFullDocument(options.UpdateLookup).SetFullDocumentBeforeChange(options.WhenAvailable).SetResumeAfter(resumeToken)
	} else {
		opt = options.ChangeStream().SetFullDocument(options.UpdateLookup).SetFullDocumentBeforeChange(options.WhenAvailable).SetStartAtOperationTime(&primitive.Timestamp{
			T: uint32(time.Now().Add(time.Duration(-60) * time.Second).Unix()),
			I: 0,
		})
	}
	watcher, err := client.Database("scrm-task").Collection("account").Watch(context.TODO(), pipeline, opt)
	if err != nil {
		this.Logger.Errorf("watch mongodb failed,%+v", err.Error())
		return
	}
	this.Logger.Infof("watch account collection successfully")
	for watcher.Next(context.TODO()) {
		var stream StreamObject
		err = watcher.Decode(&stream)
		if err != nil {
			this.Logger.Errorf("decode watch data failed,%+v", err.Error())
			continue
		}
		this.Logger.Infof("fullDocument:%+v", stream.FullDocument)
		this.Logger.Infof("fullDocumentBeforeChange:%+v", stream.FullDocumentBeforeChange)
		//保存现在resumeToken
		resumeToken = watcher.ResumeToken()
		switch stream.OperationType {
		case OperationTypeInsert:
			c.insertAccount(stream.FullDocument)
		case OperationTypeDelete:
			c.deleteAccount(stream.FullDocumentBeforeChange)
		case OperationTypeUpdate:
			c.updateAccount(stream.FullDocument)
		}
		resumeTokenJson, _ := jsoniter.MarshalToString(resumeToken)
		setAccountResumeTokenCmd := this.Redis.Set(context.TODO(), "account_resume_token", resumeTokenJson, 0)
		if setAccountResumeTokenCmd.Err() != nil {
			this.Logger.Errorf("%+v", setAccountResumeTokenCmd.Err())
			return
		}
	}
}

func (c *SyncDataConsole) syncMessageData() {
	tgCstAppListenDBSource := this.Config.MongoDBServiceConfig.DBMap["tg-cst-app"].ListenDBSource
	this.Logger.Infof(tgCstAppListenDBSource)
	clientOptions := options.Client().ApplyURI(tgCstAppListenDBSource).SetMaxPoolSize(1000)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		this.Logger.Errorf("connect to mongodb failed,%+v", err.Error())
		return
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		this.Logger.Errorf("ping mongodb failed,%+v", err.Error())
		return
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match",
			bson.M{"operationType": bson.M{"$in": bson.A{"insert", "delete", "replace", "update"}}},
		}},
	}
	getMessageResumeTokenCmd := this.Redis.Get(context.TODO(), "message_resume_token")
	resumeTokenJson := ""
	if getMessageResumeTokenCmd.Err() != nil {
		if getMessageResumeTokenCmd.Err() != redis.Nil {
			this.Logger.Errorf("get message resume token failed,%+v", getMessageResumeTokenCmd.Err())
			return
		}
	} else {
		resumeTokenJson = getMessageResumeTokenCmd.Val()
	}
	var opt *options.ChangeStreamOptions
	if len(resumeTokenJson) > 0 {
		err := jsoniter.UnmarshalFromString(resumeTokenJson, &resumeToken)
		if err != nil {
			this.Logger.Errorf("%+v", err.Error())
			return
		}
		opt = options.ChangeStream().SetFullDocument(options.UpdateLookup).SetFullDocumentBeforeChange(options.WhenAvailable).SetResumeAfter(resumeToken)
	} else {
		opt = options.ChangeStream().SetFullDocument(options.UpdateLookup).SetFullDocumentBeforeChange(options.WhenAvailable).SetStartAtOperationTime(&primitive.Timestamp{
			T: uint32(time.Now().Add(time.Duration(-10) * time.Second).Unix()),
			I: 0,
		})
	}
	watcher, err := client.Database("tg-cst-app").Collection("private_chat_message").Watch(context.TODO(), pipeline, opt)
	if err != nil {
		this.Logger.Errorf("watch mongodb failed,%+v", err.Error())
		return
	}
	this.Logger.Infof("watch private_chat_message  collection successfully")
	for watcher.Next(context.TODO()) {
		var stream StreamObject
		err = watcher.Decode(&stream)
		if err != nil {
			this.Logger.Errorf("decode watch data failed,%+v", err.Error())
			continue
		}
		this.Logger.Infof("fullDocument:%+v", stream.FullDocument)
		this.Logger.Infof("fullDocumentBeforeChange:%+v", stream.FullDocumentBeforeChange)
		//保存现在resumeToken
		resumeToken = watcher.ResumeToken()
		switch stream.OperationType {
		case OperationTypeInsert:
			c.insertMessage(stream.FullDocument)
		case OperationTypeDelete:
			c.deleteMessage(stream.FullDocumentBeforeChange)
		case OperationTypeUpdate:
			c.updateMessage(stream.FullDocument)
		}
		resumeTokenJson, _ := jsoniter.MarshalToString(resumeToken)
		setMessageResumeTokenCmd := this.Redis.Set(context.TODO(), "message_resume_token", resumeTokenJson, 0)
		if setMessageResumeTokenCmd.Err() != nil {
			this.Logger.Errorf("%+v", setMessageResumeTokenCmd.Err())
			return
		}
	}
}

func (c *SyncDataConsole) SyncData() {
	go c.syncAccountData()
	go c.syncMessageData()
}
