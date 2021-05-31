package service

import (
	"babygrow/DB/redis"
	"babygrow/logger"
	"babygrow/service/weixin/model"
	"babygrow/service/weixin/utils/aes"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
)

// 获取sessionKey到
// 通常情况下需要前端小程序调用wx.checkSession来判断，sessionKey是否已经过期，过期的情况下才需要调用该方法
func Code2SessionKey(auth *model.MiniProgramUserAuth) (*model.MiniSessionKey, error) {
	// 请求session_key需要小程序的appid，appsecret，和小程序wx.login返回的code
	if auth.AppId == "" || auth.AppSecret == "" || auth.Code == "" {
		return nil, errors.New("参数错误")
	}
	var err error
	var params req.QueryParam = req.QueryParam{
		"appid":      auth.AppId,
		"secret":     auth.AppSecret,
		"js_code":    auth.Code,
		"grant_type": "authorization_code",
	}
	resp, err := req.Get("https://api.weixin.qq.com/sns/jscode2session", params)
	if err != nil {
		return nil, err
	}
	defer resp.Response().Body.Close()
	result := new(model.MiniSessionKey)
	if err = resp.ToJSON(&result); err != nil {
		return nil, err
	}
	if result.Errcode != 0 {
		return nil, errors.New(result.Errmsg)
	}
	return result, nil
}

// 设置存储到缓存中的session_key的key
func formatStorageKey(code string) string {
	return "sessionkey" + code
}

// 存储sessionKey到缓存
// 过期时间30分钟
func StorageSessionKey(code string, sessionKey string) error {
	if code == "" {
		return errors.New("参数错误")
	}
	return redis.NewClient().Set(formatStorageKey(code), sessionKey, time.Hour*2).Err()
}

// 根据code获取本地存储的sessionKey
func GetCachedSessionKey(auth *model.MiniProgramUserAuth) (string, error) {
	if auth.Code == "" {
		return "", errors.New("参数错误")
	}
	sessionKey, err := redis.NewClient().Get(formatStorageKey(auth.Code)).Result()
	if err == nil && sessionKey != "" {
		// 如果缓存中存在对应的session_ke，返回结果
		return sessionKey, nil
	}
	return "", errors.New("获取缓存session_key失败")
}

// 根据code获取本地存储的sessionKey
func GetSessionKey(auth *model.MiniProgramUserAuth) (string, error) {
	if auth.Code == "" {
		return "", errors.New("参数错误")
	}
	if sessionKey, err := GetCachedSessionKey(auth); err == nil {
		return sessionKey, nil
	}
	// 如果缓存中不存在对应的session_key，则通过code，重新去微信服务器请求一遍session_key
	sessionKeyMap, err := Code2SessionKey(auth)
	if err != nil {
		return "", err
	}
	// 顺带再去缓存一下，成功与否在这里不重要，不影响主流程
	go StorageSessionKey(auth.Code, sessionKeyMap.SessionKey)
	return sessionKeyMap.SessionKey, nil
}

// 获取微信授权用户信息
func GetMiniUserInfo(auth *model.MiniProgramUserAuth) ([]byte, error) {
	// 验证获取信息所需的基本参数是否满足
	if auth.AppId == "" || auth.AppSecret == "" || auth.EncryptedData == "" || auth.Iv == "" {
		return nil, errors.New("参数错误")
	}

	var aesKey, ivDecode, encryptedData []byte

	// 获取对应code的session_key
	sessionKey, err := GetSessionKey(auth)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"event": "获取session_key",
		}).Error(err)
		return nil, err
	}
	auth.SessionKey = sessionKey

	// 对称解密使用的算法为 AES-128-CBC，数据采用PKCS#7填充。
	// 对称解密的目标密文为 Base64_Decode(encryptedData)。
	// 对称解密秘钥 aeskey = Base64_Decode(session_key), aeskey 是16字节。
	// 对称解密算法初始向量 为Base64_Decode(iv)，其中iv由数据接口返回。
	if aesKey, err = base64.StdEncoding.DecodeString(auth.SessionKey); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"event": "解码session_key",
		}).Error(err)
		return nil, err
	}
	if ivDecode, err = base64.StdEncoding.DecodeString(auth.Iv); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"event": "解码iv",
		}).Error(err)
		return nil, err
	}
	if encryptedData, err = base64.StdEncoding.DecodeString(auth.EncryptedData); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"event": "解码encryptedData",
		}).Error(err)
		return nil, err
	}

	if data, err := aes.CBCAesDecrypt(encryptedData, aesKey, ivDecode); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"event": "微信小程序用户信息解密",
		}).Error(err)
		return nil, err
	} else {
		return data, nil
	}
}

// 存微信的用户信息，如果相同的openid已经存在，则去更新用户
func StoreUserInfo(userInfo *model.MiniProgramUserInfo, needUpdate bool) (string, error) {
	user := new(model2.GrUser)
	var err error
	if userInfo.OpenId != "" {
		if user, err = user.Get(&model2.GrUser{
			OpenId: userInfo.OpenId,
		}); err != nil && err != sql.ErrNoRows {
			return "", err
		}
	}
	var gender string
	if userInfo.Gender == 1 {
		gender = "男"
	} else if userInfo.Gender == 2 {
		gender = "女"
	} else {
		gender = "无"
	}
	// 去更新
	if user != nil && user.ID != "" {
		if needUpdate {
			updateQuery := &model2.UpdateByIDQuery{
				ID:       user.ID,
				Nickname: userInfo.NickName,
				Avatar:   userInfo.AvatarUrl,
				OpenId:   userInfo.OpenId,
				UnionId:  userInfo.UnionId,
				Gender:   gender,
			}
			if err = user.Update(updateQuery); err != nil {
				return "", err
			}
		}
		return user.ID, nil

	} else {
		// 去新增
		user = &model2.GrUser{
			Nickname: userInfo.NickName,
			OpenId:   userInfo.OpenId,
			UnionId:  userInfo.UnionId,
			Gender:   gender,
			Avatar:   userInfo.AvatarUrl,
		}
		return user.Insert()
	}
}
