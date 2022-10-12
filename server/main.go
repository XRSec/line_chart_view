package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type (
	ViewOpts struct {
		DB *gorm.DB
	}
	App struct {
		Num  int
		Time string `gorm:"primarykey"`
	}
	Datas struct {
		Code int `json:"code"`
		Data []struct {
			Count   int     `json:"count"`
			Res     string  `json:"res"`
			Traffic float64 `json:"traffic"`
			UpRecv  float64 `json:"up_recv"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
)

var (
	view      ViewOpts
	DBName    = time.Now().Format("2006-01")
	Debug     bool
	Url       string
	ApiKey    string
	ApiSecret string
	err       error
)

func init() {
	flag.BoolVar(&Debug, "d", false, "填充数据库")
	flag.StringVar(&Url, "u", "http://localhost:8080", "API地址")
	flag.StringVar(&ApiKey, "ak", "123456", "API Key")
	flag.StringVar(&ApiSecret, "as", "123456", "API Secret")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	cwd, _ := os.Getwd()
	if _, err := os.Stat(cwd + "/view.db"); err != nil {
		if _, err = os.Create(cwd + "/view.db"); err != nil {
			log.Errorf("创建数据库失败: %v", err)
		}
	}
	view.DB, err = gorm.Open(sqlite.Open("view.db"), &gorm.Config{})
	if err = view.DB.Table(DBName).AutoMigrate(&App{}); err != nil {
		log.Errorf("初始化数据库失败: %v", err)
		return
	}

	if Debug {
		log.Info("开始填充数据库")
		populateDB()
		os.Exit(0)
	}
}
func main() {
	//test()
	//os.Exit(0)
	c := cron.New()
	if _, err := c.AddFunc("*/12 * * * *", Collect); err != nil {
		log.Errorf("添加定时任务失败: %v", err)
		return
	}
	log.Info("Collect Start")
	c.Start()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"sbcnm": "woainishabi",
	}))
	cwd, _ := os.Getwd()
	authorized.Static("/home", cwd+"/static")
	authorized.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/home")
		c.Abort()
	})
	authorized.GET("/api/top", func(c *gin.Context) {
		var apps []App
		timeV1 := time.Now().Add(-time.Hour * 24).Format("2006-01-02 15:04")
		timeV2 := time.Now().Format("2006-01-02 15:04")
		numV1 := c.Query("num")

		// 获取数量
		numV2, err := strconv.ParseInt(numV1, 10, 64)
		if err != nil {
			goto end
		}

		if int(numV2) == 119 {

			view.DB.Table(DBName).Where("time <= ? AND time > ? ", timeV2, timeV1).Order("time DESC").Limit(119).Find(&apps)
			apps = append(apps[1:]) // 删除第一个数据，然后前端会重新获取
			c.JSON(http.StatusOK, gin.H{
				"data": apps,
				"code": http.StatusOK,
			})
			return
		} else if numV2 == 1 {
			view.DB.Table(DBName).Where("time <= ? ", timeV2).Order("time DESC").Limit(int(numV2)).Find(&apps)
			c.JSON(http.StatusOK, gin.H{
				"data": apps,
				"code": http.StatusOK,
			})
			return
		} else {
			goto end
		}

	end:
		log.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "曹尼玛,别瞎几把测试",
		})
	})
	log.Infoln("Server Start, Port: 8081")
	if err := r.Run(":8081"); err != nil {
		return
	}
}

func Collect() {
	client := &http.Client{}
	os.Getenv("API-URL")
	req, _ := http.NewRequest("GET", Url, nil)
	req.Header.Set("api-key", ApiKey)
	req.Header.Set("api-secret", ApiSecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("请求接口失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Errorf("关闭数据包失败: %v", err)
		}
	}(resp.Body)

	var dates Datas
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("读取数据包失败: %v", err)
		return
	}

	if err = json.Unmarshal(body, &dates); err != nil {
		log.Errorf("解析数据包失败: %v, %v", body, err)
		return
	}

	if dates.Code != 0 || resp.StatusCode != http.StatusOK || body == nil {
		log.Errorln("请求异常: %v", resp.Body)
		return
	}

	num := 0
	for i := 0; i < len(dates.Data); i++ {
		num += dates.Data[i].Count
	}
	view.DB.Table(DBName).Create(&App{Num: num, Time: time.Now().Format("2006-01-02 15:04")})
	log.Infoln("成功保存数据")
}

// test 生成数据
func populateDB() {
	for l := 1; l < time.Now().Day(); l++ {
		for i := 0; i < 24; i++ {
			for _, k := range []int{0, 12, 24, 36, 48} {
				view.DB.Table(DBName).Create(&App{Num: rand.Intn(1000) + 1000, Time: fmt.Sprintf("2022-10-%02d %02d:%02d", l, i, k)})
			}
		}
	}
}
