package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pjebs/restgate"
	"github.com/unrolled/secure"
	"net/http"
)

func main() {
	//r := gin.Default()
	////r.Use(httpsHandler()) // https对应的中间件
	////r.Use(authMiddlewareStatic())
	//r.Use(authMiddlewareDatabase())
	//r.GET("/auth1", func(c *gin.Context) {
	//	fmt.Println(c.Request.Host)
	//	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "result": "验证通过"})
	//})
	////r.Run(":8080")
	//path := "/home/wangzhitao/go_workspace/media-server/CA/"
	//r.RunTLS(":8080", path+"com-tomsky-media-0919111751_chain.crt", path+"com-tomsky-media-0919111751_key.key")

	path := "/home/wangzhitao/Documents/music"
	r := gin.Default()
	r.StaticFS("/music", http.Dir(path))
	r.Run(":8080")
	//service.ParseMP3(path)
}

func httpsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true, // 只允许https请求
			//SSLHost:              "",      // http到https重定向
			STSSeconds:           1536000, // Strict-Transport-Security header的时效:1年
			STSIncludeSubdomains: true,    //includeSubdomains will be appended to the Strict-Transport-Security header
			STSPreload:           true,    //STS Preload(预加载)
			FrameDeny:            true,    //X-Frame-Options 有三个值:DENY（表示该页面不允许在 frame 中展示，即便是在相同域名的页面中嵌套也不允许）、SAMEORIGIN、ALLOW-FROM uri
			ContentTypeNosniff:   true,    //禁用浏览器的类型猜测行为,防止基于 MIME 类型混淆的攻击
			BrowserXssFilter:     true,    //启用XSS保护,并在检查到XSS攻击时，停止渲染页面
			// IsDevelopment: true, // 开发模式
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, "数据不安全")
			return
		}
		if status := context.Writer.Status(); status > 300 && status < 399 {
			context.Abort()
			return
		}
		context.Next()
	}
}

func authMiddlewareStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		gate := restgate.New("X-Auth-Key",
			"X-Auth-Secret",
			restgate.Static,
			restgate.Config{
				Key:                []string{"admin", "tomsky"},
				Secret:             []string{"admin", "123456"},
				HTTPSProtectionOff: false, // 如果是http服务，这里使用true
			})
		nextCalled := false
		nextAdapter := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Next()
		}
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if nextCalled == false {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func authMiddlewareDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		gate := restgate.New("X-Auth-Key",
			"X-Auth-Secret",
			restgate.Database,
			restgate.Config{
				DB:                 db,
				TableName:          "auth",             //数据表名
				Key:                []string{"key"},    //key对应的字段名
				Secret:             []string{"secret"}, //secret对应的字段名
				HTTPSProtectionOff: false,              // 如果是http服务，这里使用true
			})
		nextCalled := false
		nextAdapter := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Next()
		}
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if nextCalled == false {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

var db *sql.DB

func init() {
	db, _ = SqlDB()
}

func SqlDB() (*sql.DB, error) {
	DB_TYPE := "mysql"
	DB_HOST := "localhost"
	DB_PORT := "3306"
	DB_USER := "root"
	DB_NAME := "media"
	DB_PASSWORD := "123456"
	openString := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open(DB_TYPE, openString)
	return db, err

}
