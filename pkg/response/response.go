package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Response is a formatted struct for api results
type Response struct {
	// Code is the status code by default, but also can be
	// a custom code
	Code int `json:"code,omitempty"`
	// Message shows detail thing back to caller
	Message string `json:"message,omitempty"`
	// RequestID needs to be used with middleware
	RequestID string `json:"request_id,omitempty"`
	// Data accepts any thing as the response data
	Data interface{} `json:"data,omitempty"`
} //@name Response

// Resp returns the custom response
func Resp(c *fiber.Ctx, statusCode int, res Response) error {
	if res.Code == 0 {
		res.Code = statusCode
	}

	if id := c.Response().Header.Peek(fiber.HeaderXRequestID); len(id) > 0 && res.RequestID == "" {
		res.RequestID = utils.GetString(id)
	}

	return c.Status(statusCode).JSON(res)
}

// Data returns data with status code OK by default
func Data(c *fiber.Ctx, data interface{}) error {
	return Resp(c, fiber.StatusOK, Response{Data: data})
}

// Message responses with 200 and specific message
func Message(c *fiber.Ctx, msg string) error {
	return CodeMessage(c, fiber.StatusOK, msg)
}

// Messagef responses with 200 and specific formatted message
func Messagef(c *fiber.Ctx, format string, args ...interface{}) error {
	return CodeMessagef(c, fiber.StatusOK, fmt.Sprintf(format, args...))
}

// CodeMessage responses with specific status code and message
func CodeMessage(c *fiber.Ctx, code int, msg string) error {
	return RespCommon(c, code, msg)
}

// CodeMessagef responses with specific status code and formatted message
func CodeMessagef(c *fiber.Ctx, code int, format string, args ...interface{}) error {
	return CodeMessage(c, code, fmt.Sprintf(format, args...))
}

// RespCommon
func RespCommon(c *fiber.Ctx, code int, msg ...string) error {
	res := Response{
		// Message: utils.StatusMessage(code),
	}

	// if len(msg) > 0 {
	res.Message = msg[0]
	// }

	if strings.Contains(res.Message, "no rows in result set") {
		code = http.StatusBadRequest
		res.Message = "Not found"
	}

	if strings.Contains(res.Message, "mismatched param and argument count") ||
		strings.Contains(res.Message, "number of field descriptions must equal number") ||
		strings.Contains(res.Message, "SQLSTATE") || strings.Contains(res.Message, "can't scan dest") {

		code = http.StatusInternalServerError
		res.Message = "unexpected error occurred! Please contact administrator"
	}
	return Resp(c, code, res)
}

// RespOK responses with status code 200 RFC 7231, 6.3.1
func RespOK(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusOK, msg...)
}

// RespCreated responses with status code 201 RFC 7231, 6.3.2
func RespCreated(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusCreated, msg...)
}

// RespAccepted responses with status code 202 RFC 7231, 6.3.3
func RespAccepted(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusAccepted, msg...)
}

// RespNonAuthoritativeInformation responses with status code 203 RFC 7231, 6.3.4
func RespNonAuthoritativeInformation(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusNonAuthoritativeInformation, msg...)
}

// RespNoContent responses with status code 204 RFC 7231, 6.3.5
func RespNoContent(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusNoContent, msg...)
}

// RespResetContent responses with status code 205 RFC 7231, 6.3.6
func RespResetContent(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusResetContent, msg...)
}

// RespPartialContent responses with status code 206 RFC 7233, 4.1
func RespPartialContent(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusPartialContent, msg...)
}

// RespMultiStatus responses with status code 207 RFC 4918, 11.1
func RespMultiStatus(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusMultiStatus, msg...)
}

// RespAlreadyReported responses with status code 208 RFC 5842, 7.1
func RespAlreadyReported(c *fiber.Ctx, msg ...string) error {
	return RespCommon(c, fiber.StatusAlreadyReported, msg...)
}

func New(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}
