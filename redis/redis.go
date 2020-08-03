package redis

import (
	"github.com/garyburd/redigo/redis"
    "strconv"
)

var DBconn redis.Conn
var err error


func Line(host string)error{
	DBconn, err = redis.Dial("tcp", host)
    if err != nil {
        return err
    }
    return nil
}

func SetTokenExpire(token string,time int)(err error){
    _, err = DBconn.Do("EXPIRE", token, time)  
    return
}

func SetToken(token string, key []byte)(err error){
    _, err = DBconn.Do("SET", token, key)
    return
}

func SetTokenEX(token string, key []byte,second int)(err error){
    sec := strconv.Itoa(second)
    _, err = DBconn.Do("SET",token, key, "EX", sec)
	return
}

func GetToken(token string)(data []byte ,err error){
    data, err = redis.Bytes(DBconn.Do("GET", token))
    return
}

func CheckToken(token string)(exist bool,err error) {
    exist, err = redis.Bool(DBconn.Do("EXISTS", token))  
    return 
}

func DelToken(token string)(err error){
    _,err = DBconn.Do("DEL",token)
    return
}