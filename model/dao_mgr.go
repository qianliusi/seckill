package model

import (
	"github.com/jmoiron/sqlx"
	etcdClient "go.etcd.io/etcd/clientv3"
)

var (
	ActDao    = NewActivityDao()
	ProDao    = NewProductDao()
	ActProDao = NewActivityProductDao()
)
var (
	Db         *sqlx.DB
	EtcdClient *etcdClient.Client
)

func Init(db *sqlx.DB, etcdClient *etcdClient.Client) {
	Db = db
	EtcdClient = etcdClient
}
