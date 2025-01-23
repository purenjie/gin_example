package middleware

// import (
// 	"bytes"
// 	"crypto/md5"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"gin.example.com/entity"
// 	"gin.example.com/entity/config"
// 	"gin.example.com/middleware/log"
// 	"github.com/gin-gonic/gin"
// )

// // VerifySign 请求验签
// func VerifySign() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		app := c.GetHeader("app")
// 		ts := c.GetHeader("ts")
// 		nonce := c.GetHeader("nonce")
// 		sign := c.GetHeader("sign")
// 		// check
// 		if app == "" {
// 			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
// 				"code": entity.ErrCodeRequest,
// 				"msg":  entity.ErrApp.Error(),
// 			})
// 			return
// 		}
// 		now := int(time.Now().Unix())
// 		timestamp, _ := strconv.Atoi(ts)
// 		// 5min 有效期
// 		if timestamp > now+300 || timestamp < now-300 {
// 			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
// 				"code": entity.ErrCodeRequest,
// 				"msg":  entity.ErrTS.Error(),
// 			})
// 			return
// 		}
// 		secret := config.GetAppSecret(app)
// 		if secret == "" {
// 			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
// 				"code": entity.ErrCodeRequest,
// 				"msg":  entity.ErrApp.Error(),
// 			})
// 			return
// 		}
// 		params, err := getParams(c.Request)
// 		if err != nil {
// 			log.Errorf("getParams|err: %s", err.Error())
// 			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
// 				"code": entity.ErrCodeSystem,
// 				"msg":  entity.ErrSystem.Error(),
// 			})
// 			return
// 		}
// 		validSignStr := getValidSignStr(params, app, ts, secret, nonce)

// 		validSign := fmt.Sprintf("%x", md5.Sum([]byte(validSignStr)))
// 		if validSign != sign {
// 			log.Errorf("sign: %s, validSign: %s, validSignStr: %s", sign, validSign, validSignStr)
// 			c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
// 				"code": entity.ErrCodeRequest,
// 				"msg":  entity.ErrSign.Error(),
// 			})
// 			return
// 		}
// 		c.Next()
// 	}
// }

// func getParams(r *http.Request) (url.Values, error) {
// 	var params url.Values
// 	if r.Method == http.MethodGet {
// 		if err := r.ParseForm(); err != nil {
// 			return nil, fmt.Errorf("GET|parse form error: %s", err.Error())
// 		}
// 		params = r.Form
// 	}
// 	if r.Method == http.MethodPost {
// 		if r.Header.Get("content-type") == "application/json" {
// 			form := make(map[string]interface{})
// 			decoder := json.NewDecoder(r.Body)
// 			if err := decoder.Decode(&form); err != nil {
// 				return nil, fmt.Errorf("POST|parse json data error: %s", err.Error())
// 			}
// 			params = make(url.Values)
// 			for k, v := range form {
// 				switch v := v.(type) {
// 				case string:
// 					params[k] = []string{v}
// 				default:
// 					val, _ := json.Marshal(v)
// 					params[k] = []string{string(val)}
// 				}
// 			}

// 			body, _ := json.Marshal(form)
// 			r.Body = io.NopCloser(bytes.NewReader(body))
// 		} else {
// 			if err := r.ParseForm(); err != nil {
// 				return nil, fmt.Errorf("POST|parse form error: %s", err.Error())
// 			}
// 			params = r.Form
// 		}
// 	}
// 	return params, nil
// }

// // getValidSignStr 获取待加密的字符串
// func getValidSignStr(params url.Values, app, timestamp, secret, nonce string) string {
// 	var keys []string
// 	for k := range params {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	var kv []string
// 	for _, key := range keys {
// 		var value = params[key][0]
// 		kv = append(kv, key+"="+value)
// 	}
// 	src := strings.Join(kv, "&")
// 	return app + secret + timestamp + src + nonce
// }
