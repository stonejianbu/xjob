// Package xdb 本地键值存储
package xdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/stonejianbu/xjob/xmsg"
	"os"
	"path"
	"time"
)

// BucketName 存储桶名称
var BucketName = "msg_v.0"
var dbFile = path.Join(CurPath(), "msg.db")

// CurPath 获取当前路径
func CurPath() string {
	str, _ := os.Getwd()
	return str
}

// Set 持久化存储键值
func Set(k string, v xmsg.Msg) {
	db, _ := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	_ = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(BucketName))
		if err != nil {
			b = tx.Bucket([]byte(BucketName))
		}
		ret, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(k), ret); err != nil {
			return err
		}
		return nil
	})
}

// Get 获取持久化键值
func Get(k string) xmsg.Msg {
	msg := xmsg.Msg{}
	db, _ := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			return errors.New("no this bucket name")
		}
		ret := b.Get([]byte(k))
		if err := json.Unmarshal(ret, &msg); err != nil {
			return err
		}
		return nil
	})
	return msg
}

// Del 删除持久化
func Del(k string) {
	db, _ := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		_ = b.Delete([]byte(k))
		return nil
	})
}

// GetAll 获取所有的持久化键值
func GetAll() map[string]xmsg.Msg {
	msgs := make(map[string]xmsg.Msg)
	db, _ := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(BucketName))
		if b == nil {
			fmt.Println("no this bucket name")
			return errors.New("no this bucket name")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			msg := xmsg.Msg{}
			if err := json.Unmarshal(v, &msg); err != nil {
				return err
			}
			msgs[string(k)] = msg
		}

		return nil
	})
	return msgs
}
