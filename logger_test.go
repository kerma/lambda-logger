package logger

import (
	"bufio"
	"bytes"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestLogger(t *testing.T) {
	var expectedOut = "{\"key\":\"value\",\"message\":\"Message nr 1\"}\n"
	var out bytes.Buffer
	w := bufio.NewWriter(&out)
	log := New()
	log.w = w
	log.BindString("key", "value").SetMessageKey("message").Printf("%s nr %d", "Message", 1)

	_ = w.Flush()
	if expectedOut != out.String() {
		t.Fatalf("Received: %#v", out.String())
	}
}

func TestLoggerWithEvent(t *testing.T) {
	var expectedWithEvent = "{\"apiStage\":\"v1\",\"error\":\"HTTP ERROR: 400\",\"requestId\":\"f56d1423-cb54-41a1-84d2-f2c133d1819a\",\"requestMethod\":\"POST\",\"requestPath\":\"/path\"}\n"
	var out bytes.Buffer
	w := bufio.NewWriter(&out)
	e := events.APIGatewayProxyRequest{
		Resource:                        "",
		Path:                            "/path",
		HTTPMethod:                      "POST",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext: events.APIGatewayProxyRequestContext{
			AccountID:        "",
			ResourceID:       "",
			OperationName:    "",
			Stage:            "v1",
			DomainName:       "",
			DomainPrefix:     "",
			RequestID:        "f56d1423-cb54-41a1-84d2-f2c133d1819a",
			Protocol:         "",
			Identity:         events.APIGatewayRequestIdentity{},
			ResourcePath:     "",
			Authorizer:       nil,
			HTTPMethod:       "",
			RequestTime:      "",
			RequestTimeEpoch: 0,
			APIID:            "",
		},
		Body:            "",
		IsBase64Encoded: false,
	}

	logger := NewFromRequest(e)
	logger.w = w
	logger.Errorf("HTTP ERROR: %d", 400)

	_ = w.Flush()
	if expectedWithEvent != out.String() {
		t.Fatalf("Received: %#v", out.String())
	}
}
