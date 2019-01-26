package common

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"
)

type handlerFun interface{}

// HandleRequest get a handler function to process the request by request url
func HandleRequest(context *gin.Context, fun handlerFun) {
	args, err := buildHandleFuncArgs(fun, context)
	if err != nil {
		RespondErrorResp(context, err)
		return
	}

	result := callHandleFunc(fun, args...)
	if err := result[len(result)-1]; err != nil {
		RespondErrorResp(context, err.(error))
		return
	}

	if exist := processPaginationIfPresent(fun, args, result, context); exist {
		return
	}

	if len(result) == 1 {
		RespondSuccessResp(context, nil)
		return
	}

	RespondSuccessResp(context, result[0])
}

func callHandleFunc(fun handlerFun, args ...interface{}) []interface{} {
	fv := reflect.ValueOf(fun)

	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}

	rs := fv.Call(params)
	result := make([]interface{}, len(rs))
	for i, r := range rs {
		result[i] = r.Interface()
	}
	return result
}

func processPaginationIfPresent(fun handlerFun, args []interface{}, result []interface{}, context *gin.Context) bool {
	ft := reflect.TypeOf(fun)
	if ft.NumIn() != 3 {
		return false
	}

	list := result[0]
	size := reflect.ValueOf(list).Len()
	query := args[2].(*PaginationQuery)

	paginationInfo := &PaginationInfo{Start: query.Start, Limit: query.Limit, HasNext: size == int(query.Limit)}
	RespondSuccessPaginationResp(context, list, paginationInfo)
	return true
}

func buildHandleFuncArgs(fun handlerFun, context *gin.Context) ([]interface{}, error) {
	args := []interface{}{context}

	req, err := createHandleReqArg(fun, context)
	if err != nil {
		return nil, errors.Wrap(err, "createHandleReqArg")
	}

	if err := checkDisplayOrder(req); err != nil {
		return nil, err
	}

	if req != nil {
		args = append(args, req)
	}

	ft := reflect.TypeOf(fun)

	// not exist pagination
	if ft.NumIn() != 3 {
		return args, nil
	}

	query, err := ParsePagination(context)
	if err != nil {
		return nil, errors.Wrap(err, "ParsePagination")
	}

	args = append(args, query)
	return args, nil
}

func checkDisplayOrder(req interface{}) error {
	if req == nil {
		return nil
	}

	reqType := reflect.TypeOf(req)
	reqVal := reflect.ValueOf(req)

	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
		reqVal = reqVal.Elem()
	}
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		if field.Type != reflect.TypeOf(Display{}) {
			continue
		}
		display := reqVal.Field(i).Interface().(Display)
		if strings.Trim(display.Sorter.By, "") == "" {
			return nil
		}

		order := strings.Trim(display.Sorter.Order, "")
		if order != "desc" && order != "asc" {
			display.Sorter.Order = "desc"
		}
	}
	return nil
}

func createHandleReqArg(fun handlerFun, context *gin.Context) (interface{}, error) {
	ft := reflect.TypeOf(fun)
	if ft.NumIn() == 1 {
		return nil, nil
	}
	argType := ft.In(1)
	argKind := argType.Kind()

	// point type must dereference once
	if argKind == reflect.Ptr {
		argType = argType.Elem()
	}

	reqArg := reflect.New(argType).Interface()
	if err := context.ShouldBindJSON(reqArg); err != nil {
		return nil, errors.Wrap(err, "bind reqArg")
	}

	b, err := json.Marshal(reqArg)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}

	context.Set(ReqBodyLabel, string(b))

	if argKind == reflect.Ptr {
		return reqArg, nil
	}

	return reflect.ValueOf(reqArg).Elem().Interface(), nil
}

var (
	errorType           = reflect.TypeOf((*error)(nil)).Elem()
	contextType         = reflect.TypeOf((*gin.Context)(nil))
	paginationQueryType = reflect.TypeOf((*PaginationQuery)(nil))
)

func ValidateFuncType(fun handlerFun) error {
	ft := reflect.TypeOf(fun)
	if ft.Kind() != reflect.Func || ft.IsVariadic() {
		return errors.New("need nonvariadic func in " + ft.String())
	}

	if ft.NumIn() < 1 || ft.NumIn() > 3 {
		return errors.New("need one or two or three parameters in " + ft.String())
	}

	if ft.In(0) != contextType {
		return errors.New("the first parameter must point of context in " + ft.String())
	}

	if ft.NumIn() == 2 && ft.In(1).Kind() != reflect.Struct && ft.In(1).Kind() != reflect.Ptr {
		return errors.New("the second parameter must struct or point in " + ft.String())
	}

	if ft.NumIn() == 3 && ft.In(2) != paginationQueryType {
		return errors.New("the third parameter of pagination must point of paginationQuery in " + ft.String())
	}

	if ft.NumOut() < 1 || ft.NumOut() > 2 {
		return errors.New("the size of return value must one or two in " + ft.String())
	}

	// if has pagination, the first return value must slice or array
	if ft.NumIn() == 3 && ft.Out(0).Kind() != reflect.Slice && ft.Out(0).Kind() != reflect.Array {
		return errors.New("the first return value of pagination must slice of array in " + ft.String())
	}

	if !ft.Out(ft.NumOut() - 1).Implements(errorType) {
		return errors.New("the last return value must error in " + ft.String())
	}
	return nil
}
