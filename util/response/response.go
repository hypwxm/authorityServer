package response

import "worldbar/logger"

type Sender struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageData struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func NewSender() *Sender {
	return &Sender{}
}

func (self *Sender) Fail(msg string) {

	logger.Logger.Debug(msg)

	self.Code = 1
	self.Message = msg
	self.Data = nil
}

func (self *Sender) Success(data interface{}) {
	self.Code = 0
	self.Message = "请求成功"
	self.Data = data
}

func (self *Sender) SuccessList(data interface{}, total int) {
	self.Code = 0
	self.Message = "请求成功"
	self.Data = &PageData{
		Data:  data,
		Total: total,
	}
}
