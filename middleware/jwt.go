package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/xddzb/dousheng/model"
	"net/http"
	"time"
)

type MyCustomClaims struct {
	ID int64
	jwt.StandardClaims
}
type Response struct {
	StatusCode int64
	StatusMsg  string
}

var jwtkey = []byte("lucky_dai")

func GenerateToken(user model.UserLogin) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(300 * time.Minute)
	issuer := "frank"
	claims := MyCustomClaims{
		ID: user.UserInfoId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //token过期时间
			Issuer:    issuer,            //签发者
			IssuedAt:  time.Now().Unix(), //签发时间
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtkey)
	return token, err
}

func ParseToken(token string) (*MyCustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyCustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// GWT鉴权中间件
func JWTMidWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token") //获取post请求参数
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 403,
				StatusMsg:  "token不正确",
			})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, Response{
				StatusCode: 402,
				StatusMsg:  "token过期",
			})
			c.Abort()
			return
		}
		c.Set("user_id", tokenStruck.ID) //在请求上下文里面设置一些值，然后其他地方取值
		c.Next()
	}
}
