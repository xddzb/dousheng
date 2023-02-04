package service

import (
	"errors"
	"github.com/xddzb/dousheng/middleware"
	"github.com/xddzb/dousheng/model"
)

// 定义register接口返回的数据
type RegisterResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

const (
	MaxUsernameLength = 100
)

// AddUser 注册用户并得到token和id
func AddUser(username, password string) (*RegisterResponse, error) {
	return NewUser(username, password).Do()
}

func NewUser(username, password string) *AddUserFlow {
	return &AddUserFlow{username: username, password: password}
}

type AddUserFlow struct {
	//传入的参数
	username string
	password string
	//中间参数
	userid int64
	token  string
	//传给controller层的参数
	data *RegisterResponse
}

func (f *AddUserFlow) Do() (*RegisterResponse, error) {
	//对参数进行合法性验证
	if err := f.checkNum(); err != nil {
		return nil, err
	}

	//更新数据到数据库
	if err := f.updateData(); err != nil {
		return nil, err
	}

	//打包response
	if err := f.packResponse(); err != nil {
		return nil, err
	}
	return f.data, nil
}

func (f *AddUserFlow) checkNum() error {
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

func (f *AddUserFlow) updateData() error {
	//判断用户名是否已经存在
	userLoginDAO := model.NewUserLoginDao()
	if userLoginDAO.IsUserExistByUsername(f.username) {
		return errors.New("用户名已存在")
	}

	//准备好userInfo,默认name为username
	userLogin := model.UserLogin{Username: f.username, Password: f.password}
	userinfo := model.UserInfo{User: &userLogin, Name: f.username}
	//添加新用户信息到数据库，由于userLogin属于userInfo，故更新userInfo即可
	userInfoDAO := model.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return err
	}

	//颁发token
	token, err := middleware.GenerateToken(userLogin)
	if err != nil {
		return err
	}
	f.token = token
	f.userid = userinfo.Id
	return nil
}

func (f *AddUserFlow) packResponse() error {
	f.data = &RegisterResponse{
		UserId: f.userid,
		Token:  f.token,
	}
	return nil
}
