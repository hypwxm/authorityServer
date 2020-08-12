package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// CommonError 微信返回的通用错误json
type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// DecodeWithCommonError 将返回值按照CommonError解析
func DecodeWithCommonError(response []byte, apiName string) (err error) {
	var commError CommonError
	err = json.Unmarshal(response, &commError)
	if err != nil {
		return
	}
	if commError.ErrCode != 0 {
		return fmt.Errorf("%s Error , errcode=%d , errmsg=%s", apiName, commError.ErrCode, commError.ErrMsg)
	}
	return nil
}

/**
 * 对错误进行转换
 * 比如sql的 no rows in sql -> 数据不存在
 */
func ErrorFormat(err error) string {
	if err == sql.ErrNoRows {
		return "数据不存在"
	}
	return err.Error()
}
