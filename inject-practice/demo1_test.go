package inject_practice

import (
	"github.com/facebookgo/inject"
	"golanglearning/new_project/for-gong-zhong-hao/Practical-practice/inject-practice/model"
	"testing"
)

func TestUseInject(t *testing.T) {

	conf := model.LoadConf()
	db := model.InitDb()


	server := model.Server{}
	graph := inject.Graph{}

	if err := graph.Provide(
		&inject.Object{
			Value: &server,
		},
		&inject.Object{
			Value: db,
		},
		&inject.Object{
			Value: conf,
		},
	); err != nil {
		panic(err)
	}

	if err := graph.Populate(); err != nil {
		panic(err)
	}

	server.Run()

}