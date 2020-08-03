/*Copyright (C) 2020 EricWu <skynocover@gmail.com>
See COPYING for license details.*/
package main

import (
	"RedisDesktop/redis"
	"RedisDesktop/save"
	"fmt"
	"github.com/zserge/lorca"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

var (
	vieweight  int = 1000
	viewheight int = 780

	ui, err            = lorca.New("", "", vieweight, viewheight)
	storages tstorages = tstorages{}
)

type URL struct {
	HostPort string
}

type tstorage struct {
	key   string
	value string
}

type tstorages struct {
	arr []tstorage
}

func (this *tstorages) Push(key, value string) {
	storage := tstorage{
		key:   key,
		value: value,
	}
	if len(this.arr) < 6 {
		this.arr = append(this.arr, storage)
	} else {
		this.arr = append(this.arr[1:], storage)
	}
	for i := range this.arr {
		ui.Eval(`document.getElementById("key` + strconv.Itoa(len(this.arr)-i-1) + `").textContent = "` + this.arr[i].key + `"`)
		ui.Eval(`document.getElementById("value` + strconv.Itoa(len(this.arr)-i-1) + `").textContent = "` + this.arr[i].value + `"`)
	}
}

func main() {
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	ui.Load("data:text/html," + url.PathEscape(read("./file/index.html")))
	ThisURL := URL{}
	save.Load("file/saveurl", &ThisURL)
	ui.Eval(`document.getElementById("input").value = "` + ThisURL.HostPort + `"`)

	//點選連線按鈕
	ui.Bind("connect", func(url string) {
		err = redis.Line(url)
		if err != nil {
			ui.Eval(`window.alert("Redis連線失敗` + fmt.Sprintf("%v", err) + `")`)
		} else {
			ui.Eval(`window.alert("Redis連線成功")`)
			ThisURL.HostPort = url
			save.Save("file/saveurl", &ThisURL)
		}
	})

	//點選查詢按鈕
	ui.Bind("srh", func(key string) {
		value, err := redis.GetToken(key)
		if err != nil {
			ui.Eval(`window.alert("Redis查詢失敗: ` + fmt.Sprintf("%v", err) + `")`)
		} else {
			storages.Push(key, fmt.Sprintf("%s", value))
		}
	})

	//點選刪除按鈕
	ui.Bind("del", func(selected string) {
		log.Println(selected)
	})

	<-ui.Done()
}

//Tool

//讀取檔案
func read(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
