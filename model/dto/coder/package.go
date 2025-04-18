package coder

import (
	"github.com/yunhanshu-net/runcher/model"
	"strings"
)

type BizPackage struct {
	Runner *model.Runner `json:"runner"`

	AbsPackagePath string `json:"abs_package_path"`
	Language       string `json:"language"`
	EnName         string `json:"en_name"`
	CnName         string `json:"cn_name"`
	Desc           string `json:"desc"`
}

func (c *BizPackage) GetPackageSaveFullPath(sourceCodeDir string) (savePath string, absPackagePath string) {
	savePath = strings.TrimSuffix(sourceCodeDir, "/") + "/api"
	absPackagePath = savePath + "/" + c.AbsPackagePath
	return savePath, absPackagePath
}

func (c *BizPackage) GetPackageName() string {
	return c.EnName
}

type CreateProjectReq struct {
	model.Runner
}
type CreateProjectResp struct {
	Version string `json:"version"`
}

type AddApisResp struct {
	Version string               `json:"version"`
	ErrList []*CodeApiCreateInfo `json:"err_list"`
}

type AddApiResp struct {
	Version string `json:"version"`
}

type BizPackageResp struct {
	Version string `json:"version"`
}
