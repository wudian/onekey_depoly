package redis

import (
	"gitlab.zhonganinfo.com/tech_bighealth/za-delos/midapi/api"

	"fmt"

	"time"

	r "github.com/garyburd/redigo/redis"
)

type Client struct {
	pool *r.Pool
}

func NewClient(args ...interface{}) api.DataBase {
	return &Client{
		pool: &r.Pool{
			MaxIdle:     args[1].(int),
			IdleTimeout: time.Duration(args[2].(int)) * time.Second,
			Dial: func() (r.Conn, error) {
				server := args[0].(string)
				c, err := r.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				if true {
					// 地址：r-tj7748e7ace28aa4.redis.rds.aliyuncs.com
					// 端口：6379
					// 密码：H043sHKO4OGPhAqM
					if _, err := c.Do("AUTH", "H043sHKO4OGPhAqM"); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
		},
	}

}

func (c *Client) OnClose() {
	c.pool.Close()
}

func (c *Client) Insert(args ...interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if len(args) < 2 {
		err = fmt.Errorf("Redis Insert Params nums Error")
		return
	}
	_, err = conn.Do("SET", args[0].(string), args[1])
	if err != nil {
		return
	}

	if len(args) > 2 {
		conn.Do("EXPIRE", args[0].(string), args[2].(int))
	}

	return
}

func (c *Client) Get(args ...interface{}) (value []byte, err error) {

	conn := c.pool.Get()
	defer conn.Close()

	return r.Bytes(conn.Do("GET", args[0].(string)))

}

func (c *Client) IsExist(args ...interface{}) (ok bool, err error) {
	conn := c.pool.Get()
	defer conn.Close()
	ok, err = r.Bool(conn.Do("EXISTS", args[0].(string)))
	if err != nil {
		return
	}

	return
}

func (c *Client) Del(args ...interface{}) (err error) {

	_, err = c.pool.Get().Do("DEL", args[0].(string))

	return
}

func (c *Client) HInsert(args ...interface{}) (err error) {
	conn := c.pool.Get()
	defer conn.Close()

	if len(args) < 3 {
		err = fmt.Errorf("Redis HInsert Params nums Error")
		return
	}
	_, err = conn.Do("HSET", args[0].(string), args[1].(string), args[2])
	if err != nil {
		return
	}

	if len(args) > 3 {
		conn.Do("EXPIRE", args[0].(string), args[3].(int))
	}

	return
}
func (c *Client) HGet(args ...interface{}) (value []byte, err error) {

	conn := c.pool.Get()
	defer conn.Close()

	return r.Bytes(conn.Do("HGET", args[0].(string), args[1].(string)))

}

func (c *Client) HIsExist(args ...interface{}) (ok bool, err error) {
	conn := c.pool.Get()
	defer conn.Close()
	ok, err = r.Bool(conn.Do("HEXISTS", args[0].(string), args[1].(string)))
	if err != nil {
		return
	}

	return
}

func (c *Client) HDel(args ...interface{}) (err error) {

	_, err = c.pool.Get().Do("HDEL", args[0].(string), args[1].(string))

	return
}
