package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Runner struct {
	//Kind       string `json:"kind"`        //类型，可执行程序，so文件等等

	WorkPath        string `json:"work_path"`
	Command         string `json:"command"`
	RequestJsonPath string `json:"request_json_path"`
	Language        string `json:"language"`   //编程语言
	StoreRoot       string `json:"store_root"` //oss 存储的跟路径
	Name            string `json:"name"`       //应用名称（英文标识）
	ToolType        string `json:"tool_type"`  //工具类型
	Version         string `json:"version"`    //应用版本
	OssPath         string `json:"oss_path"`   //文件地址
	User            string `json:"user"`       //所属租户
}

func (r *Runner) GetBuildRunnerName() string {
	return fmt.Sprintf("%s_%s_%s", r.User, r.Name, r.GetNextVersion())
}

func (r *Runner) GetBuildPath(root string) string {
	return fmt.Sprintf("%s/%s/%s/bin", root, r.User, r.Name)
}

func (r *Runner) GetVersionNum() (int, error) {
	replace := strings.ReplaceAll(r.Version, "v", "")
	version, err := strconv.Atoi(replace)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return version, nil
}

func (r *Runner) GetNextVersion() string {
	num, err := r.GetVersionNum()
	if err != nil {
		fmt.Println("GetVersionNum err:" + err.Error())
	}
	return fmt.Sprintf("v%d", num+1)
}

func (r *Runner) GetInstallPath(rootPath string) string {
	return fmt.Sprintf("%s/%s/%s/version/%s", strings.TrimSuffix(rootPath, "/"), r.User, r.Name, r.Version)
}

func (r *Runner) GetToolPath(rootPath string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(rootPath, "/"), r.User, r.Name)
}

func (r *Runner) GetNextVersionInstallPath(rootPath string) (string, error) {
	nextVersion := r.GetNextVersion()
	return fmt.Sprintf("%s/%s/%s/version/%s", strings.TrimSuffix(rootPath, "/"), r.User, r.Name, nextVersion), nil
}

func (r *Runner) Check() error {
	if r.Name == "" {
		return fmt.Errorf("name 不能为空")
	}
	if r.ToolType == "" {
		return fmt.Errorf("ToolType 不能为空")
	}
	if r.Version == "" {
		return fmt.Errorf("version 不能为空")
	}

	if r.User == "" {
		return fmt.Errorf("user 不能为空")
	}
	return nil
}

type UpdateVersion struct {
	RunnerConf *Runner `json:"runner_conf"`
	OldVersion string  `json:"old_version"`
	//NewVersion        string  `json:"new_version"`
	//NewVersionOssPath string  `json:"new_version_oss_path"`
}
