package inject_practice

import (
	"golanglearning/new_project/for-gong-zhong-hao/Practical-practice/inject-practice/model"
	"testing"
)

func TestNoUseInject(t *testing.T) {
	conf := model.LoadConf()
	db := model.InitDb()

	userService := &model.UserService{
		Db: db,
		Conf: conf,
	}

	postService := &model.PostService{
		Db: db,
	}

	userHandler := &model.UserController{
		UserService: userService,
		Config: conf,
	}

	postHandler := &model.PostController{
		UserService: userService,
		PostService: postService,
		Conf: conf,
	}

	server := &model.Server{
		UserApi: userHandler,
		PostApi: postHandler,
	}

	server.Run()



}
