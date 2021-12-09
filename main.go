package main

import (
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Idx  string `json:"idx" xml:"idx" form:"idx" query:"idx"`
	Pw   string `json:"pw" xml:"pw" form:"pw" query:"pw"`
}

func (User) TableName() string {
	return "user"
}
func main() {
	dsn := "admin:diary1234@tcp(teamdiary.ciszsojgfruj.ap-northeast-2.rds.amazonaws.com)/hello?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Db연결에 실패")
	}

	db.AutoMigrate(&User{})

	e := echo.New()
	e.Static("/static", "page")

	e.File("/join", "page/join.htm")
	e.File("/complate", "page/complate.htm")

	//가입
	e.POST("/join", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		if len(u.Name) == 0 {
			return c.HTML(http.StatusOK, "<script>window.location.href='http://localhost:1323/join';alert('네임을 입력하세요.')</script>")
		}
		if len(u.Idx) < 5 {
			return c.HTML(http.StatusOK, "<script>window.location.href='http://localhost:1323/join';alert('아이디를 입력하세요.')</script>")
		}
		if len(u.Pw) < 8 {
			return c.HTML(http.StatusOK, "<script>window.location.href='http://localhost:1323/join';alert('패스워드를 입력하세요.')</script>")
		}
		ur := &User{
			Idx:  u.Idx,
			Pw:   u.Pw,
			Name: u.Name,
		}
		db.Create(&ur)
		return c.HTML(http.StatusOK, "<script>window.location.href='http://localhost:1323/complate';</script>")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
