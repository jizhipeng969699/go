package controllers

import (
	"bytes"
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

func zzinit() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Error("redis连接失败")
		return
	}
	defer conn.Close()

	//resp, err := conn.Do("set", "c1", "hello")
	//resp, err := conn.Do("del", "c1")
	//resp, err := conn.Do("get", "c1")
	//	resp, err := conn.Do("mset", "name", "xiaoming", "age", "18", "addr", "beijing")
	//n, _ := redis.String(resp, err)
	resp, err := conn.Do("mget", "name", "age", "addr")

	result, _ := redis.Values(resp, err)

	var age int
	var name, addr string

	redis.Scan(result, &name, &age, &addr) //这里的参数要和 do 的 一一对应

	beego.Info("mset：", name, age, addr)

	resp1, err := conn.Do("keys", "*")
	slic, _ := redis.Strings(resp1, err)
	beego.Info("slic", slic)

}

//三步骤  操作redis数据库
func RRfunc() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Info("redis数据库连接失败")
		return
	}
	defer conn.Close()

	//Send函数发出指令
	err = conn.Send("set", "name", "xiaoming")
	if err != nil {
		beego.Info("send函数发送指令失败")
		return
	}
	err = conn.Send("set", "age", "18")
	if err != nil {
		beego.Info("send函数发送指令失败")
		return
	}

	//flush将连接的输出缓冲区刷新到服务器
	err = conn.Flush()
	if err != nil {
		beego.Info("flush刷数据失败")
		return
	}

	//Receive接收服务器返回的数据
	result, err := conn.Receive()
	if err != nil {
		beego.Info("receive接受数据失败")
		return
	}

	beego.Info("结果---", result.(string))
}

//使用do  scan函数操作数据库
func DoRfunc() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Info("redis数据库连接失败")
		return
	}
	defer conn.Close()

	rep, err := conn.Do("mget", "age", "name")
	if err != nil {
		beego.Info("命令失败")
		return
	}

	//	strs, _ := redis.Strings(rep, err)

	strs, _ := redis.Values(rep, err)//按照数据的对应类型 设置 redis的类型

	var age int
	var name string
	str, _ := redis.Scan(strs, &age, &name)

	var buffer bytes.Buffer        //创建容器 相当于缓冲区
	enc := gob.NewEncoder(&buffer) //创建编码器

	//var a int
	enc.Encode(age) //这里进行编码  编码后放到了 buffer容器中
	beego.Info("jieguo-----------", strs)
	beego.Info("values", str)
	beego.Info("values", age, name)

}

//加密和解密
func aRfunc() {
	age := "15"

	conn, _ := redis.Dial("tcp", ":6379")
	defer conn.Close()

	var buffer bytes.Buffer        //创建容器   相当于缓冲区
	enc := gob.NewEncoder(&buffer) //创建加密  把加密好的数据放到  容器中

	enc.Encode(age) // 给数据加密

	conn.Do("set", "age", buffer.Bytes())
	beego.Info(buffer)

	data, _ := conn.Do("get", "age")//返回的是接口  去要进行类型断言
	beego.Info(data, "---------")
	var a string
	dec := gob.NewDecoder(bytes.NewReader(data.([]byte)))

	dec.Decode(&a)

	beego.Info("---------------", a)
}
