package service

import (
	"errors"
	"github.com/xddzb/dousheng/middleware"
	"github.com/xddzb/dousheng/model"
)

// QueryUserLogin 注册用户并得到token和id
func QueryUserLogin(username, password string) (*RegisterResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

type QueryUserLoginFlow struct {
	//传入的参数
	username string
	password string
	//中间参数
	userid int64
	token  string
	//传给controller层的参数
	data *RegisterResponse
}

func (f *QueryUserLoginFlow) Do() (*RegisterResponse, error) {
	//对参数进行合法性验证
	if err := f.checkNum(); err != nil {
		return nil, err
	}

	//查询用户是否存在并返回token和id
	if err := f.prepareDate(); err != nil {
		return nil, err
	}

	//打包response
	if err := f.packResponse(); err != nil {
		return nil, err
	}
	return f.data, nil
}

func (f *QueryUserLoginFlow) checkNum() error {
	if f.username == "" {
		return errors.New("用户名为空")
	}
	if len(f.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if f.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (f *QueryUserLoginFlow) prepareDate() error {
	userLoginDAO := model.NewUserLoginDao()
	var login model.UserLogin
	//准备好userid
	err := userLoginDAO.QueryUserLogin(f.username, f.password, &login)
	if err != nil {
		return err
	}
	f.userid = login.UserInfoId

	//颁发token
	token, err := middleware.GenerateToken(login)
	if err != nil {
		return err
	}
	f.token = token
	return nil
}

func (f *QueryUserLoginFlow) packResponse() error {
	f.data = &RegisterResponse{
		UserId: f.userid,
		Token:  f.token,
	}
	return nil
}
