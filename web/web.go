package web

import (
	"fmt"
	"net/http"
	"strconv"
	"trojan/core"
	"trojan/util"
	"trojan/web/controller"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

func userRouter(router *gin.Engine) {
	user := router.Group("/xray/user")
	{
		user.GET("", func(c *gin.Context) {
			requestUser := RequestUsername(c)
			if requestUser == "admin" {
				c.JSON(200, controller.UserList(""))
			} else {
				c.JSON(200, controller.UserList(requestUser))
			}
		})
		user.GET("/page", func(c *gin.Context) {
			curPageStr := c.DefaultQuery("curPage", "1")
			pageSizeStr := c.DefaultQuery("pageSize", "10")
			curPage, _ := strconv.Atoi(curPageStr)
			pageSize, _ := strconv.Atoi(pageSizeStr)
			c.JSON(200, controller.PageUserList(curPage, pageSize))
		})
		user.POST("", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")
			c.JSON(200, controller.CreateUser(username, password))
		})
		user.POST("/update", func(c *gin.Context) {
			sid := c.PostForm("id")
			username := c.PostForm("username")
			password := c.PostForm("password")
			c.JSON(200, controller.UpdateUser(sid, username, password))
		})
		user.POST("/expire", func(c *gin.Context) {
			sid := c.PostForm("id")
			sDays := c.PostForm("useDays")
			useDays, _ := strconv.Atoi(sDays)
			c.JSON(200, controller.SetExpire(sid, uint(useDays)))
		})
		user.DELETE("/expire", func(c *gin.Context) {
			sid := c.Query("id")
			c.JSON(200, controller.CancelExpire(sid))
		})
		user.DELETE("", func(c *gin.Context) {
			sid := c.Query("id")
			c.JSON(200, controller.DelUser(sid))
		})
	}
}

func xrayRouter(router *gin.Engine) {
	router.POST("/xray/start", func(c *gin.Context) {
		c.JSON(200, controller.Start())
	})
	router.POST("/xray/stop", func(c *gin.Context) {
		c.JSON(200, controller.Stop())
	})
	router.POST("/xray/restart", func(c *gin.Context) {
		c.JSON(200, controller.Restart())
	})
	router.GET("/xray/loglevel", func(c *gin.Context) {
		c.JSON(200, controller.GetLogLevel())
	})
	router.POST("/xray/update", func(c *gin.Context) {
		c.JSON(200, controller.Update())
	})
	// router.POST("/xray/switch", func(c *gin.Context) {
	// 	tType := c.DefaultPostForm("type", "xray")
	// 	c.JSON(200, controller.SetTrojanType(tType))
	// })
	router.POST("/xray/loglevel", func(c *gin.Context) {
		slevel := c.DefaultPostForm("level", "1")
		// level, _ := strconv.Atoi(slevel)
		c.JSON(200, controller.SetLogLevel(slevel))
	})
	router.POST("/xray/domain", func(c *gin.Context) {
		c.JSON(200, controller.SetDomain(c.PostForm("domain")))
	})
	router.GET("/xray/log", func(c *gin.Context) {
		controller.Log(c)
	})
}

func dataRouter(router *gin.Engine) {
	data := router.Group("/xray/data")
	{
		data.POST("", func(c *gin.Context) {
			sID := c.PostForm("id")
			sQuota := c.PostForm("quota")
			quota, _ := strconv.Atoi(sQuota)
			c.JSON(200, controller.SetData(sID, quota))
		})
		data.DELETE("", func(c *gin.Context) {
			sID := c.Query("id")
			c.JSON(200, controller.CleanData(sID))
		})
		data.POST("/resetDay", func(c *gin.Context) {
			dayStr := c.DefaultPostForm("day", "1")
			day, _ := strconv.Atoi(dayStr)
			c.JSON(200, controller.UpdateResetDay(uint(day)))
		})
		data.GET("/resetDay", func(c *gin.Context) {
			c.JSON(200, controller.GetResetDay())
		})
	}
}

func commonRouter(router *gin.Engine) {
	common := router.Group("/common")
	{
		common.GET("/version", func(c *gin.Context) {
			c.JSON(200, controller.Version())
		})
		common.GET("/serverInfo", func(c *gin.Context) {
			c.JSON(200, controller.ServerInfo())
		})
		common.POST("/loginInfo", func(c *gin.Context) {
			c.JSON(200, controller.SetLoginInfo(c.PostForm("title")))
		})
	}
}

func staticRouter(router *gin.Engine) {
	box := packr.New("trojanBox", "./templates")
	router.Use(func(c *gin.Context) {
		requestUrl := c.Request.URL.Path
		if box.Has(requestUrl) || requestUrl == "/" {
			http.FileServer(box).ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}

// Start web启动入口
func Start(host string, port int, isSSL bool) {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	staticRouter(router)
	router.Use(Auth(router).MiddlewareFunc())
	xrayRouter(router)
	userRouter(router)
	dataRouter(router)
	commonRouter(router)
	controller.SheduleTask()
	controller.CollectTask()
	util.OpenPort(port)
	if isSSL {
		config := core.Load("")
		ssl := &config.Inbounds[0].StreamSettings.XtlsSettings.Certificates[0]
		router.RunTLS(fmt.Sprintf("%s:%d", host, port), ssl.CertificateFile, ssl.KeyFile)
	} else {
		router.Run(fmt.Sprintf("%s:%d", host, port))
	}
}
