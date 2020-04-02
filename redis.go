package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

func execCmd(conn redis.Conn, cmd string) string {
	ss := strings.Split(strings.TrimSpace(cmd), " ")
	tmp := make([]string, 0, len(ss))
	for _, v := range ss {
		if t := strings.TrimSpace(v); len(t) > 0 {
			tmp = append(tmp, t)
		}
	}
	fmt.Println(tmp, "length:", len(tmp))
	if len(tmp) == 0 {
		return ""
	}
	args := make([]interface{}, 0, len(tmp))
	for _, v := range tmp[1:] {
		args = append(args, v)
	}
	resp, err := conn.Do(tmp[0], args...)
	if err != nil {
		return err.Error()
	}
	return stringfyResponse(resp)
}

func stringfyResponse(resp interface{}) string {
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
		fmt.Println("resp：", resp)
		if resp == nil {
			return "<nil>"
		} else {
			fmt.Println("resp type:", reflect.TypeOf(resp).String())
		}
		return "unknown response type"
	}
}

func connectToRedis(host string, port int, password string) (redis.Conn, error) {
	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Second),
		redis.DialWriteTimeout(time.Second),
		redis.DialReadTimeout(time.Second),
		redis.DialKeepAlive(time.Second * 30),
	}
	if password != "" {
		options = append(options, redis.DialPassword(password))
	}
	return redis.Dial("tcp4", fmt.Sprintf("%s:%d", host, port), options...)
}
