package regex

import (
	"errors"
	"github.com/dubbogo/gost/log/logger"
	"regexp"
)

const (
	FILENAME_PATTERN = "^(\\S+)_(\\S+)_(\\S+)" + "_REQ_" + "(\\S+)_(\\S+)\\.(json|xml|zip)$"
)

type FilenameReq struct {
	AppId       string `json:"appId"`
	FilenamePre string `json:"filenamePre"`
	Filename    string `json:"filename"`
	Timestamp   string `json:"timestamp"`
}

func MatchFileName(filename string) (*FilenameReq, error) {
	pattern := regexp.MustCompile(FILENAME_PATTERN)
	// 执行匹配
	matches := pattern.FindStringSubmatch(filename)
	var appid string
	var filenamePre string
	var timestamp string
	// 检查是否匹配
	if len(matches) > 6 {
		logger.Infof("Match found: %s", matches)
		// 输出匹配到的每个部分（模拟元组）
		appid = matches[4]
		filenamePre = matches[1] + "_" + matches[2] + "_" + matches[3]
		timestamp = matches[5]
	} else {
		return nil, errors.New("File name not match")
	}
	req := &FilenameReq{AppId: appid, FilenamePre: filenamePre, Timestamp: timestamp, Filename: filename}
	return req, nil
}
