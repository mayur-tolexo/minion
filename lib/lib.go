package lib

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/tolexo/aero/conf"
)

//NewResp : Create New Response Objext
func NewResp() *Resp {
	return &Resp{}
}

//GetErrorCode : Get Error Code from *Error
func GetErrorCode(err error) (code int) {
	if err == nil {
		code = NO_ERROR
	} else {
		switch err.(type) {
		case *Error:
			code = err.(*Error).Code
		default:
			code = NO_ERROR
		}
	}
	return
}

//GetErrorHTTPCode : Get HTTP code from error code. Please check for httpCode = 0.
func GetErrorHTTPCode(err error) (httpCode int, status bool) {
	errCode := GetErrorCode(err)
	httpCode = ErrorHTTPCode[errCode]
	if err == nil {
		status = true
	}
	return
}

//Error : Implement Error method of error interface
func (e *Error) Error() string {
	return fmt.Sprintf("\nCode:\t\t[%d]\nMessage:\t[%v]\nStackTrace:\t[%v]\nDebugMsg:\t[%v]\n", e.Code, e.Msg, e.Trace, e.DebugMsg)
}

//newError : Create new *Error object
func newError(msg string, err error, code int, debugMsg ...string) *Error {
	stackDepth := conf.Int("error.stack_depth", 3)
	funcName, fileName, line := StackTrace(stackDepth)
	trace := fileName + " -> " + funcName + ":" + strconv.Itoa(line)
	errStr := strings.Join(debugMsg, " ")
	if err != nil {
		errStr += " " + err.Error()
	}

	return &Error{
		Msg:      msg,
		DebugMsg: errStr,
		Trace:    trace,
		Code:     code,
	}
}

//StackTrace : Get function name, file name and line no of the caller function
//Depth is the value from which it will start searching in the stack
func StackTrace(depth int) (funcName string, file string, line int) {
	var (
		ok bool
		pc uintptr
	)
	for i := depth; ; i++ {
		if pc, file, line, ok = runtime.Caller(i); ok {
			if trackAll := conf.Bool("error.track_all", false); trackAll == false &&
				strings.Contains(file, PACKAGE_NAME) {
				continue
			}
			fileName := strings.Split(file, "github.com")
			if len(fileName) > 1 {
				file = fileName[1]
			}
			_, funcName = packageFuncName(pc)
			break
		} else {
			break
		}
	}
	return
}

//packageFuncName : Package and function name from package counter
func packageFuncName(pc uintptr) (packageName string, funcName string) {
	if f := runtime.FuncForPC(pc); f != nil {
		funcName = f.Name()
		if ind := strings.LastIndex(funcName, "/"); ind > 0 {
			packageName += funcName[:ind+1]
			funcName = funcName[ind+1:]
		}
		if ind := strings.Index(funcName, "."); ind > 0 {
			packageName += funcName[:ind]
			funcName = funcName[ind+1:]
		}
	}
	return
}

//Set Response Values
func (r *Resp) Set(data interface{}, status bool, err error) {
	if err == nil {
		r.Data = data
	} else {
		r.Error = err
	}
	r.Status = status
}

//BuildResponse : creates the response of API
func BuildResponse(data interface{}, err error) (int, interface{}) {
	out := NewResp()
	httpCode, status := GetErrorHTTPCode(err)
	out.Set(data, status, err)
	return httpCode, out
}

//Unmarshal : unmarshal of request body and wrap error
func Unmarshal(reqBody string, structModel interface{}) (err error) {
	if err = jsoniter.Unmarshal([]byte(reqBody), structModel); err != nil {
		err = UnmarshalError(err)
	}
	return
}
