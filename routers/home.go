package routers

import (
	"strconv"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/middleware"
	"github.com/qiniu/log"
)

type RepoInfoForm struct {
	Name string `form:"repo-name" binding:"required"`
}

type RepositoryForm struct {
	Rid int64 `form:"rid" binding:"required"`
}

type TaskForm struct {
	Tid int64 `form:"tid" binding:"required"`
}

func NewRepo(rf RepoInfoForm, ctx *middleware.Context) {
	var err error
	rf.Name = base.SanitizedRepoPath(rf.Name)
	if _, err = models.CreateRepository(rf.Name); err != nil {
		log.Errorf("create repo error: %v", err)
	}
	ctx.Redirect(302, "/")
}

func NewBuild(rf RepositoryForm, ctx *middleware.Context) {
	task := new(models.Task)
	task.Arch = "amd64"
	task.Os = "darwin"
	task.Branch = "master"
	task.Sha = "xxxxxx"
	task.CgoEnable = false
	task.Rid = rf.Rid
	if _, err := models.CreateTask(task); err != nil {
		log.Errorf("create module error: %v", err)
	}
	ctx.Redirect(302, "/repo?id="+strconv.Itoa(int(rf.Rid)))
}

func ForceRebuild(tf TaskForm, ctx *middleware.Context) {
	if err := models.ResetTask(tf.Tid); err != nil {
		log.Errorf("reset task failed: %v", err)
	}
	ctx.Redirect(302, "/history?id="+strconv.Itoa(int(tf.Tid)))
}

func Home(ctx *middleware.Context) {
	pv := models.RefreshPageView("/")
	repos, err := models.GetAllRepos(50, 0)
	if err != nil {
		log.Errorf("get repos from db error: %v", err)
	}
	ctx.Data["Title"] = "home"
	ctx.Data["Repos"] = repos
	ctx.Data["PV"] = pv
	ctx.HTML(200, "home")
}

func Ruok() string {
	return "imok"
}

func NotFound(ctx *middleware.Context) {
	ctx.Handle(404, "Where you got this page", nil)
}
