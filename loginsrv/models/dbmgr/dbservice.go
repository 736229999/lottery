package dbmgr

import mgo "gopkg.in/mgo.v2"

//DbService Db服务针对每一个单独的服务器,Mango DB 不允许在同一台物理机上监听多个端口
type DbService struct {
	Session       *mgo.Session
	DbMap         map[string]*mgo.Database
	CollectionMap map[string]*mgo.Collection
}

//构造一个DbService 这种伪构造的方法相当之恶心
func newDbService() *DbService {
	dbService := &DbService{}
	dbService.DbMap = make(map[string]*mgo.Database)
	dbService.CollectionMap = make(map[string]*mgo.Collection)
	return dbService
}
