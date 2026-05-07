package main

import (
	"fmt"
	"net/http"
	"sync"
	//"io"
	"log"
)

func Fetch(urls []string) <- chan Article1 {
	var wg sync.WaitGroup
<<<<<<< HEAD
	ch := make(chan Article1, 5)
=======
	ch := make(chan Article1, len(urls))
>>>>>>> d3f970c82297a06d4880e55f2fb1e2d60a6f261c

	for _, url := range urls{
		wg.Add(1)

		go func (url string)  {

			defer wg.Done()
			
			
			resp, err := http.Get(url)
<<<<<<< HEAD
			if err != nil {
				log.Printf("请求失败:[URL]: %s, %s", url, err)
				return 
			}
=======
>>>>>>> d3f970c82297a06d4880e55f2fb1e2d60a6f261c
			defer resp.Body.Close()

			fmt.Println("开始抓取：",url)

<<<<<<< HEAD
			xmlData, err := Parse(resp.Body, url)

=======
			if err != nil {
			log.Printf("请求失败:%s，%s",url,err)
			}

			if resp.StatusCode != 200 {
			
				log.Printf("状态码错误: %d", resp.StatusCode)
			}
			fmt.Println("1")
			xmlData, err := Parse(resp.Body)
			fmt.Println("2")
>>>>>>> d3f970c82297a06d4880e55f2fb1e2d60a6f261c
			if err != nil {
				log.Print("解析失败：",err)
			}

			for _, article := range xmlData {
				ch <- article
			}
		}(url)
	}
	go func() {
    	wg.Wait()   // ① 等待所有 goroutine 完成
    	close(ch)   // ② 所有任务完成后关闭通道
	}()
	fmt.Println("3")
	return ch
}
	