package main

import (
	"crawler/collect"
	"crawler/parse"
	"crawler/proxy"
	"fmt"
	"time"
	// "go.uber.org/zap"
	// "gorm.io/gorm/logger"
)

func main() {
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	p, _ := proxy.RoundRobinProxySwitcher(proxyURLs...)
	cookie := "bid=-UXUw--yL5g; dbcl2=\"214281202:q0BBm9YC2Yg\"; __yadk_uid=jigAbrEOKiwgbAaLUt0G3yPsvehXcvrs; push_noty_num=0; push_doumail_num=0; __utmz=30149280.1665849857.1.1.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.21428; ck=SAvm; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1665925405%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=*; __utma=30149280.2072705865.1665849857.1665849857.1665925407.2; __utmc=30149280; __utmt=1; __utmb=30149280.23.5.1665925419338; _pk_id.100001.8cb4=fc1581490bf2b70c.1665849856.2.1665925421.1665849856."

	var worklist []*collect.Request
	for i := 0; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		fmt.Printf("str:%v\n", str)
		worklist = append(worklist, &collect.Request{
			Url:       str,
			Cookie:    cookie,
			ParseFunc: parse.ParseURL,
		})
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				// logger.Error("read content failed",
				// 	zap.Error(err),
				// )
				continue
			}
			res := item.ParseFunc(body, item)
			for _, item := range res.Items {
				fmt.Printf("result:%v\n", item)
				// logger.Info("result",
				// 	zap.String("get url:", item.(string)))
			}
			worklist = append(worklist, res.Requests...)
		}
	}

}
