package demo

import "gitee.com/y19941115mx/ygo/framework"

type Service struct {
	IService
	container framework.Container
}

func NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &Service{container: container}, nil
}

func (s *Service) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}
}
