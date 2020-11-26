package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

const (
	defaultErrorKey   = "error"
	defaultMessageKey = "log" // both @message and @log are used by cloudwatch. Need a better key here to avoid confusion
)

var (
	environmentVars = map[string]string{
		"AWS_LAMBDA_FUNCTION_NAME":    "functionName",
		"AWS_LAMBDA_FUNCTION_VERSION": "functionVersion",
	}
)

type Output map[string]interface{}

type Serializer func(Output, io.Writer) error

func JsonSerializer() Serializer {
	return func(o Output, w io.Writer) error {
		return json.NewEncoder(w).Encode(o)
	}
}

type Logger struct {
	errorKey   string
	messageKey string
	params     params
	w          io.Writer
	serialize  Serializer
}

type params map[string]interface{}

// New creates a initiated Logger instance
func New() *Logger {
	return newWithSerializer(nil, JsonSerializer())
}

// NewFromRequest creates a new Logger instance using values from APIGatewayProxyRequest
func NewFromRequest(req events.APIGatewayProxyRequest) *Logger {
	return newWithSerializer(&req, JsonSerializer())
}

func newWithSerializer(req *events.APIGatewayProxyRequest, s Serializer) *Logger {
	l := &Logger{
		errorKey:   defaultErrorKey,
		messageKey: defaultMessageKey,
		params:     make(params, 4),
		w:          os.Stdout,
		serialize:  s,
	}
	if req != nil {
		l.BindRequest(*req)
	}
	return l
}

// SetErrorKey sets the error key in the output. Defaults to `error`
func (l *Logger) SetErrorKey(k string) *Logger {
	l.errorKey = k
	return l
}

// SetMessageKey sets the message key in the output. Defaults to `log`
func (l *Logger) SetMessageKey(k string) *Logger {
	l.messageKey = k
	return l
}

// BindString adds a string value to each log call output
func (l *Logger) BindString(k string, s string) *Logger {
	l.params[k] = s
	return l
}

// BindInt adds a int value to each log call output
func (l *Logger) BindInt(k string, i int) *Logger {
	l.params[k] = i
	return l
}

// BindNum adds a number value to each log call output
func (l *Logger) BindNum(k string, n float64) *Logger {
	l.params[k] = n
	return l
}

// BindEnv includes AWS lambda environment values to each log call output
func (l *Logger) BindEnv() *Logger {
	for e, k := range environmentVars {
		l.params[k] = os.Getenv(e)
	}
	return l
}

// BindRequest includes values from APIGatewayProxyRequest to each log call output
func (l *Logger) BindRequest(e events.APIGatewayProxyRequest) *Logger {
	l.params["requestMethod"] = e.HTTPMethod
	l.params["requestPath"] = e.Path
	l.params["requestId"] = e.RequestContext.RequestID
	l.params["apiStage"] = e.RequestContext.Stage
	if e.RequestContext.Identity.CognitoIdentityID != "" {
		l.params["cognitoIdentityId"] = e.RequestContext.Identity.CognitoIdentityID
	}
	return l
}

// Println prints a log message
func (l *Logger) Println(s string) {
	l.print(l.messageKey, s)
}

// Printf prints a formatted log message.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.print(l.messageKey, fmt.Sprintf(format, v...))
}

// Error prints a error message.
func (l *Logger) Error(s string) {
	l.print(l.errorKey, s)
}

// Errorf prints a formatted error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.print(l.errorKey, fmt.Sprintf(format, v...))
}

func (l *Logger) print(k, s string) {
	o := make(Output, len(l.params)+1)
	o[k] = s

	for k, v := range l.params {
		o[k] = v
	}

	err := l.serialize(o, l.w)
	if err != nil {
		fmt.Printf("ERROR: Could not serialize log message: %v", err)
	}
}
