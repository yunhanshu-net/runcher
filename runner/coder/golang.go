package coder

import (
	"fmt"
	"github.com/yunhanshu-net/runcher/model"
	"github.com/yunhanshu-net/runcher/model/dto/coder"
	"github.com/yunhanshu-net/runcher/pkg/codes"
	"github.com/yunhanshu-net/runcher/pkg/codex"
	"github.com/yunhanshu-net/runcher/pkg/osx"
	"github.com/yunhanshu-net/runcher/status"
	"os"
	"os/exec"
	"strings"
)

// AddBizPackage(codeBizPackage *coder.BizPackage) (*coder.BizPackageResp, error)
// AddApi(codeApi *coder.CodeApi) (*coder.AddApiResp,error)
// AddApis(codeApis []*coder.CodeApi) (resp *coder.AddApisResp, err error)
// CreateProject() (*coder.CreateProjectResp, error)
type Golang struct {
	runnerRoot string
	runner     *model.Runner
}

func (g *Golang) AddApi(codeApi *coder.CodeApi) (*coder.AddApiResp, error) {
	runner := g.runner
	pathInfo, err := g.addApi(codeApi)
	if err != nil {
		return nil, err
	}
	err = g.buildRunner(pathInfo.NextVersionPath, runner.GetBuildPath(g.runnerRoot), runner.GetBuildRunnerName())
	if err != nil {
		return nil, err
	}
	return &coder.AddApiResp{Version: runner.GetNextVersion()}, nil
}

func (g *Golang) createFile(filePath string, content string) error {
	codeFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer codeFile.Close()
	_, err = codeFile.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (g *Golang) buildRunner(workDir string, buildPath string, runnerName string) error {
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = workDir
	err := tidy.Run()
	if err != nil {
		return err
	}
	// 创建命令
	cmd := exec.Command("go", "build", "-o", "../../bin/"+runnerName)

	// 设置工作目录（可选）
	cmd.Dir = workDir // 当前目录，可以根据需要修改为项目路径

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("build failed: err:%v, output:%s", err, string(output))
	}
	fmt.Println(string(output))
	exists := osx.FileExists(buildPath + "/" + runnerName)
	if !exists {
		return status.ErrorCodeApiBuildError.WithMessage(workDir)
	}
	return nil
}

func (g *Golang) AddBizPackage(codeBizPackage *coder.BizPackage) (*coder.BizPackageResp, error) {
	runner := g.runner
	currentVersionPath := runner.GetInstallPath(g.runnerRoot)
	_, absPkgPath := codeBizPackage.GetPackageSaveFullPath(currentVersionPath)
	if osx.DirExists(absPkgPath) { //先判断Package是否存在
		return nil, status.ErrorCodeApiFileExist.WithMessage(absPkgPath)
	}
	//不存在才可以创建
	err := os.MkdirAll(absPkgPath, 0755)
	if err != nil {
		return nil, err
	}

	//在pkg下创建_init.go文件
	err = g.createFile(absPkgPath+"/init_.go", fmt.Sprintf(codes.InitCodeTemplate, codeBizPackage.EnName))
	if err != nil {
		return nil, err
	}
	manager := codex.NewGolangProjectManager(currentVersionPath)
	err = manager.AddPackages([]codex.PackageInfo{
		{
			Alias:      strings.ReplaceAll(codeBizPackage.AbsPackagePath, "/", "_"),
			ImportPath: fmt.Sprintf("git.yunhanshu.net/%s/%s/api/%s", runner.User, runner.Name, codeBizPackage.AbsPackagePath),
		},
	})
	if err != nil {
		return nil, err
	}

	return &coder.BizPackageResp{Version: runner.Version}, nil
}

func (g *Golang) CreateProject() (*coder.CreateProjectResp, error) {
	runner := g.runner
	err := os.MkdirAll(runner.GetToolPath(g.runnerRoot), 0755) //初始化项目目录
	if err != nil {
		return nil, err
	}

	//go.mod
	//main.go
	//bin
	//	data
	//	.request
	//	user_app_v1
	//  user_app_v2
	//version
	//	-v1
	//		-api
	//  -v2
	//		-api

	codePath := runner.GetToolPath(g.runnerRoot) + "/version/" + runner.Version
	err = os.MkdirAll(codePath, 0755) //初始版本目录
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(codePath+"/api", 0755) //初始api目录
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(runner.GetToolPath(g.runnerRoot)+"/bin/.request", 0755) //初始化可执行程序目录
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(runner.GetToolPath(g.runnerRoot)+"/bin/data", 0755) //初始化数据库目录
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(runner.GetToolPath(g.runnerRoot)+"/bin/conf", 0755) //初始化配置文件目录
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(runner.GetToolPath(g.runnerRoot)+"/bin/logs", 0755) //初始化logs目录
	if err != nil {
		return nil, err
	}

	//创建 go.mod 文件
	err = g.createFile(codePath+"/go.mod", fmt.Sprintf(`
module git.yunhanshu.net/%s/%s

go 1.23`, runner.User, runner.Name))
	if err != nil {
		return nil, err
	}
	//创建main文件
	manager := codex.NewGolangProjectManager(codePath)
	err = manager.CreateMain(nil)
	if err != nil {
		return nil, err
	}
	err = g.buildRunner(codePath, runner.GetBuildPath(g.runnerRoot), runner.GetBuildRunnerCurrentVersionName())
	if err != nil {
		return nil, err
	}
	return &coder.CreateProjectResp{Version: runner.Version}, nil
}

func (g *Golang) addApi(api *coder.CodeApi) (*model.RunnerPath, error) {
	runnerRoot := g.runnerRoot
	runner := g.runner
	if g.runner.Version == "" {
		version, err := runner.GetLatestVersion()
		if err != nil {
			return nil, err
		}
		runner.Version = version
	}
	pathInfo := runner.GetPaths(runnerRoot)
	api.Language = g.runner.Language

	_, addFileAbsFile := api.GetFileSaveFullPath(pathInfo.NextVersionPath)
	//if osx.DirExists(addFileSavePath) { //先判断package是否存在
	//	return nil, status.ErrorCodeApiFileExist.WithMessage(addFileSavePath)
	//}

	//先判断文件是否已经存在
	if osx.FileExists(addFileAbsFile) {
		return nil, status.ErrorCodeApiFileExist.WithMessage(addFileAbsFile)
	}

	exists := osx.DirExists(pathInfo.NextVersionPath)
	if !exists {
		//创建新版本工作目录，把当前项目代码复制到下一个版本
		err := os.MkdirAll(pathInfo.NextVersionPath, 0755)
		if err != nil {
			return nil, err
		}
		err = osx.CopyDirectory(pathInfo.CurrentVersionPath, pathInfo.NextVersionPath) //把当前项目代码保存一份复制到下一个版本
		if err != nil {
			return nil, err
		}
	}

	err := g.createFile(addFileAbsFile, api.Code) //创建新文件
	if err != nil {
		return nil, err
	}

	return &pathInfo, nil
}

func (g *Golang) AddApis(codeApis []*coder.CodeApi) (resp *coder.AddApisResp, err error) {

	resp = new(coder.AddApisResp)
	if len(codeApis) == 0 {
		return nil, fmt.Errorf("no api")
	}
	var errs []*coder.CodeApiCreateInfo
	var pathInfo *model.RunnerPath
	for _, codeApi := range codeApis {
		info, err := g.addApi(codeApi)
		if err != nil {
			errs = append(errs, &coder.CodeApiCreateInfo{
				Language:       codeApi.Language,
				Package:        codeApi.Package,
				AbsPackagePath: codeApi.AbsPackagePath,
				EnName:         codeApi.EnName,
				CnName:         codeApi.CnName,
				Msg:            err.Error(),
				Status:         "failed",
			})
		}
		if info != nil {
			pathInfo = info
		}
	}
	if pathInfo == nil {
		return nil, fmt.Errorf("pathInfo is nil")
	}
	resp.ErrList = errs
	resp.Version = g.runner.GetNextVersion()

	err = g.buildRunner(pathInfo.NextVersionPath, g.runner.GetBuildPath(g.runnerRoot), g.runner.GetBuildRunnerName())
	if err != nil {
		return nil, err
	}
	return resp, nil
}
