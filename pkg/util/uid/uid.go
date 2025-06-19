package uid

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type logger interface {
	Error(error)
}

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
var Logger logger

// ErrTimeOut 获取uid超时错误
var ErrTimeOut = errors.New("get uid timeout")

type Uid struct {
	db         *sql.DB    // 数据库连接
	businessID string     // 业务id
	ch         chan int64 // id缓冲池
	min, max   int64      // id段最小值，最大值
}

// NewUid 创建一个Uid;len：缓冲池大小()
// db:数据库连接
// businessID：业务id
// len：缓冲池大小(长度可控制缓存中剩下多少id时，去DB中加载)
func NewUid(db *sql.DB, businessID string, len int) (*Uid, error) {
	lid := Uid{
		db:         db,
		businessID: businessID,
		ch:         make(chan int64, len),
	}
	go lid.productID()
	return &lid, nil
}

// Get 获取自增id,当发生超时，返回错误，避免大量请求阻塞，服务器崩溃
func (u *Uid) Get() (int64, error) {
	select {
	case <-time.After(1 * time.Second):
		return 0, ErrTimeOut
	case uid := <-u.ch:
		return uid, nil
	}
}

// productID 生产id，当ch达到最大容量时，这个方法会阻塞，直到ch中的id被消费
func (u *Uid) productID() {
	_ = u.reLoad()

	for {
		if u.min >= u.max {
			_ = u.reLoad()
		}

		u.min++
		u.ch <- u.min
	}
}

// reLoad 在数据库获取id段，如果失败，会每隔一秒尝试一次
func (u *Uid) reLoad() error {
	var err error
	for {
		err = u.getFromDB()
		if err == nil {
			return nil
		}

		// 数据库发生异常，等待一秒之后再次进行尝试
		if Logger != nil {
			Logger.Error(err)
		} else {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

// getFromDB 从数据库获取id段
func (u *Uid) getFromDB() error {
	var (
		maxID int64
		step  int64
	)

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	row := tx.QueryRow("SELECT max_id,step FROM uid WHERE business_id = ? FOR UPDATE", u.businessID)
	err = row.Scan(&maxID, &step)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE uid SET max_id = ? WHERE business_id = ?", maxID+step, u.businessID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	u.min = maxID
	u.max = maxID + step
	return nil
}
