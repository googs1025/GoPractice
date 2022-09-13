package model

import (
	"fmt"
	"time"
)


// 控制反转的作用。
// https://www.cyhone.com/articles/facebookgo-inject/

func init()  {
	LoadConf()
	InitDb()
}

type Server struct {
	UserApi *UserController `inject:""`
	PostApi *PostController `inject:""`
}

type UserController struct {
	UserService *UserService `inject:""`
	Config 		*Conf `inject:""`
}

type PostController struct {
	UserService *UserService `inject:""`
	PostService *PostService `inject:""`
	Conf *Conf
}

type UserService struct {
	Db *Db `inject:""`
	Conf *Conf `inject:""`
}

type PostService struct {
	Db *Db `inject:""`
}


type Db struct {

}

type Conf struct {

}

func LoadConf() *Conf {
	fmt.Println("正在载入config配置文件")
	time.Sleep(time.Second * 3)
	return &Conf{}
}

func InitDb() *Db {
	fmt.Println("正在初始化db！")
	time.Sleep(time.Second * 3)
	return &Db{}
}


func (s *Server) Run() {

	fmt.Println("服务器已完成初始化，正在运行！")
}
