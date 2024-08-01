package apisvr

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/ahang7/go-IAM/internal/pkg/middleware"
	"github.com/ahang7/go-IAM/internal/pkg/middleware/auth"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	// APIServerAudience define the value of the audience field in the JWT token
	APIServerAudience = "iam.api.ch.com"
	// APIServerIssuer define the value of the issuer field in the JWT token
	APIServerIssuer = "iam-apiserver"
)

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(
		newBasicAuth().(auth.BasicStrategy),
		newJWTAuth().(auth.JWTStrategy),
	)
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username, password string) bool {
		return true
	})
}

func newJWTAuth() middleware.AuthStrategy {
	ginJWTMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.Key")),
		Timeout:          viper.GetDuration("jwt.Timeout"),
		MaxRefresh:       viper.GetDuration("jwt.MaxRefresh"),
		Authenticator:    authenticator(),
		Authorizator:     authorizator(),
		PayloadFunc:      payloadFunc(),
		Unauthorized:     Unauthorized(),
		LoginResponse:    loginResponse(),
		LogoutResponse:   logoutResponse(),
		RefreshResponse:  refreshResponse(),
		IdentityHandler:  identityHandler(),
		IdentityKey:      middleware.UserNameKey,
		TokenLookup:      "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:    "Bearer",
		TimeFunc:         time.Now,
		SendCookie:       true,
	})
	if err != nil {
		return nil
	}
	return auth.NewJWTStrategy(*ginJWTMiddleware)
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login loginInfo
		var err error

		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			return "", err
		}
		// TODO: get user information by the login username. db/cache
		// compare username and password
		// set user login time at
		// cache update

		return login.Username, err
	}
}

// 获取Authorization头的值，并调用strings.SplitN函数，获取一个切片变量auth，其值为 ["Basic","YWRtaW46QWRtaW5AMjAyMQ=="]
// 将YWRtaW46QWRtaW5AMjAyMQ==进行base64解码，得到admin:Admin@2021
// 调用strings.SplitN函数获取 admin:Admin@2021 ，得到用户名为admin，密码为Admin@2021
func parseWithHeader(c *gin.Context) (loginInfo, error) {
	authKey := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(authKey) != 2 || authKey[0] != "Basic" {
		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	payload, err := base64.StdEncoding.DecodeString(authKey[1])
	if err != nil {

		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

// 调用了Gin的ShouldBindJSON函数，来从Body中解析出用户名和密码
func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	if err := c.ShouldBind(&login); err != nil {
		return loginInfo{}, jwt.ErrFailedAuthentication
	}
	return login, nil
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if _, ok := data.(string); ok {
			return true
		}
		return false
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		// TODO: set claims[jwt.IdentityKey] = u.Name and ["sub"] = u.Name
		//if u, ok != data; ok {
		//	claims["username"] = u
		//}
		return claims
	}
}

func Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"message": message,
		})
	}
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func logoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "logout success",
		})
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func identityHandler() func(*gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return claims[jwt.IdentityKey]
	}
}
