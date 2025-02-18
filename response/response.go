package response

import (
	"github.com/gin-gonic/gin"

	"github.com/irsanrasyidin/complete_project_module/models"
	"github.com/irsanrasyidin/complete_project_module/utils"
)

type IResponse interface {
	GetStatusCode() int
	SetRequestID(c *gin.Context)
}

type ErrorResponse struct {
	StatusCode int `json:"status_code"`
	RequestID  any `json:"request_id"`
	Message    any `json:"message"`
	ErrorCode  any `json:"error_code"`
	Error      any `json:"error"`
}

func (r *ErrorResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r *ErrorResponse) SetRequestID(c *gin.Context) {
	reqID, _ := c.Get(utils.RequestID)
	r.RequestID = reqID
}

type SuccessResponse struct {
	StatusCode int `json:"status_code"`
	RequestID  any `json:"request_id"`
	Message    any `json:"message"`
}

func (r *SuccessResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r *SuccessResponse) SetRequestID(c *gin.Context) {
	reqID, _ := c.Get(utils.RequestID)
	r.RequestID = reqID
}

type DataResponse struct {
	StatusCode int `json:"status_code"`
	RequestID  any `json:"request_id"`
	Message    any `json:"message"`
	Data       any `json:"data"`
}

func (r *DataResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r *DataResponse) SetRequestID(c *gin.Context) {
	reqID, _ := c.Get(utils.RequestID)
	r.RequestID = reqID
}

type PaginationResponse struct {
	StatusCode int               `json:"status_code"`
	RequestID  any               `json:"request_id"`
	Message    any               `json:"message"`
	Pagination *models.Pagination `json:"pagination"`
	Data       any               `json:"data"`
}

func (r *PaginationResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r *PaginationResponse) SetRequestID(c *gin.Context) {
	reqID, _ := c.Get(utils.RequestID)
	r.RequestID = reqID
}