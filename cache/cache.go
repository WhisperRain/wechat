package cache

import "time"

//Cache interface
type Cache interface {
	Get(key string)  interface{}
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) (bool, error)
	Delete(key string) error

    HGet(key,field string ,reply interface{})error
	//GetWithErrorBack(key string,reply interface{}) error
	HSetWxUser(ip, agentKey string, user interface{}) error
	GetWithErrorBack(key string, reply interface{}) error
}
