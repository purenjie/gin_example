package middleware

import (
	"net/http"
	"reflect"

	"gin.example.com/entity"
	"gin.example.com/middleware/log"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func invalidParam(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    entity.ErrCodeRequest,
		Message: msg,
	})
}

// HandleBindings handles the bindings of URL path parameters, query parameters, headers and JSON body.
func HandleBindings(fn interface{}) func(c *gin.Context) {

	fnTy := reflect.TypeOf(fn)
	argNum := fnTy.NumIn()
	retNum := fnTy.NumOut()
	if argNum != 2 || (retNum != 1 && retNum != 2) {
		panic("expected type `func(*gin.Context, *Req) (*Resp, error)` or `func(*gin.Context, *Req) error`")
	}

	ctxTy := fnTy.In(0)
	if ctxTy.Kind() != reflect.Ptr {
		panic("arg #1 expected type `*gin.Context`")
	}

	reqPtrTy := fnTy.In(1)
	if reqPtrTy.Kind() != reflect.Ptr {
		panic("arg #2 expected type `*Req`")
	}

	var errTy reflect.Type
	if retNum == 2 {
		errTy = fnTy.Out(1)
		respPtrTy := fnTy.Out(0)
		k := respPtrTy.Kind()
		if k != reflect.Ptr && k != reflect.Array && k != reflect.Interface {
			panic("ret #1 expected type pointer, array or interface")
		}
	} else {
		errTy = fnTy.Out(0)
	}

	if errTy.Kind() != reflect.Interface {
		panic("ret #2 expected type `error`")
	}

	return func(c *gin.Context) {
		fnVal := reflect.ValueOf(fn)
		ctxVal := reflect.ValueOf(c)
		reqPtrVal := reflect.New(reqPtrTy.Elem())
		reqPtr := reqPtrVal.Interface()

		if err := c.ShouldBindUri(reqPtr); err != nil {
			log.Error(c, "URI binding error: err=%v", err)
			invalidParam(c, err.Error())
			return
		}
		if err := c.ShouldBindQuery(reqPtr); err != nil {
			log.Error(c, "Query binding error: err=%v", err)
			invalidParam(c, err.Error())
			return
		}
		if err := c.ShouldBindHeader(reqPtr); err != nil {
			log.Error(c, "Header binding error: err=%v", err)
			invalidParam(c, err.Error())
			return
		}
		if c.Request.ContentLength > 0 {
			if err := c.ShouldBindJSON(reqPtr); err != nil {
				log.Error(c, "JSON binding error: err=%v", err)
				invalidParam(c, err.Error())
				return
			}
		}

		// generate trace id
		traceID, _ := uuid.NewV4()
		c.Set("trace_id", traceID.String())

		retVals := fnVal.Call([]reflect.Value{ctxVal, reqPtrVal})
		var errVal reflect.Value

		if len(retVals) == 1 {
			errVal = retVals[0]
		} else {
			errVal = retVals[1]
		}

		if ei := errVal.Interface(); ei != nil {
			err := ei.(error)
			msg := err.Error()

			log.Errorf(c, "Internal error on handler: err=%v", err)

			c.JSON(http.StatusOK, ResponseData{
				Code:    entity.ErrCodeSystem,
				Message: msg,
			})
			return
		}

		if len(retVals) == 2 {
			respVal := retVals[0]
			c.JSON(http.StatusOK, struct {
				ResponseData
				Data interface{} `json:"data"`
			}{
				ResponseData{Code: 0, Message: "ok"},
				respVal.Interface(),
			})
			return
		}

		c.JSON(http.StatusOK, ResponseData{Code: 0, Message: "ok"})
	}
}
