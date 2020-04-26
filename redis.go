package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func execCmd(conn redis.Conn, cmd string, args ...interface{}) string {
	logrus.Debugln("exec cmd:", cmd, "args:", args)
	resp, err := conn.Do(cmd, args...)
	if err != nil {
		return err.Error()
	}
	return stringfyResponse(resp)
}

func stringfyResponse(resp interface{}) string {
	if resp == nil {
		return "<nil>"
	}
	switch t := resp.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.Itoa(int(t))
	case error:
		return t.Error()
	case []interface{}:
		buf := bytes.NewBuffer(nil)
		for k, v := range t {
			buf.WriteString(fmt.Sprintf("%d) ", k+1) + stringfyResponse(v) + "\r\n")
		}
		return buf.String()
	default:
		logrus.Errorln("unknown resp type:", reflect.TypeOf(resp).String())
		return "unknown response type"
	}
}

func DialRedis(host string, port int, password string) (redis.Conn, error) {
	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
		redis.DialReadTimeout(time.Second),
		redis.DialKeepAlive(time.Minute * 30),
	}
	if password != "" {
		options = append(options, redis.DialPassword(password))
	}
	return redis.Dial("tcp4", fmt.Sprintf("%s:%d", host, port), options...)
}
