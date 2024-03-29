package controllers

import (

	"dtools/models/env"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type EnvController struct {
	BaseController
}

func (self *EnvController) List() {
	self.Data["pageTitle"] = "环境设置"
	self.Display()
}

func (self *EnvController) Add() {
	self.Data["pageTitle"] = "新增环境"
	self.Display()
}

func (self *EnvController) Edit() {
	self.Data["pageTitle"] = "编辑环境"

	id, _ := self.GetInt("id", 0)
	Env, _ := env.EnvGetById(id)
	row := make(map[string]interface{})
	row["id"] = Env.Id
	row["env_name"] = Env.EnvName
	row["env_host"] = Env.EnvHost
	row["detail"] = Env.Detail
	self.Data["env"] = row
	self.Display()
}

func (self *EnvController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}
	envName := strings.TrimSpace(self.GetString("envName"))

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if envName != "" {
		filters = append(filters, "env_name__icontains", envName)
	}
	result, count := env.EnvGetList(page, self.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["env_name"] = v.EnvName
		row["detail"] = v.Detail
		row["env_host"] = v.EnvHost
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	self.AjaxList("成功", MSG_OK, count, list)
}

func (self *EnvController) AjaxSave() {
	Env_id, _ := self.GetInt("id")
	if Env_id == 0 {
		Env := new(env.Env)

		Env.EnvName = strings.TrimSpace(self.GetString("env_name"))
		Env.EnvHost = strings.TrimSpace(self.GetString("env_host"))
		Env.Detail = strings.TrimSpace(self.GetString("detail"))
		Env.CreateId = self.userId
		Env.UpdateId = self.userId
		Env.CreateTime = time.Now().Unix()
		Env.UpdateTime = time.Now().Unix()
		Env.Status = 1

		_, err := env.EnvGetByName(Env.EnvName)

		if err == nil {
			self.ajaxMsg("环境名称已经存在", MSG_ERR)
		}

		if _, err := env.EnvAdd(Env); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg("", MSG_OK)
	}

	EnvUpdate, _ := env.EnvGetById(Env_id)
	// 修改
	EnvUpdate.EnvName = strings.TrimSpace(self.GetString("env_name"))
	EnvUpdate.EnvHost = strings.TrimSpace(self.GetString("env_host"))
	EnvUpdate.Detail = strings.TrimSpace(self.GetString("detail"))
	EnvUpdate.UpdateId = self.userId
	EnvUpdate.UpdateTime = time.Now().Unix()
	EnvUpdate.Status = 1

	if err := EnvUpdate.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}

func (self *EnvController) AjaxDel() {

	Env_id, _ := self.GetInt("id")
	Env, _ := env.EnvGetById(Env_id)
	Env.UpdateTime = time.Now().Unix()
	Env.UpdateId = self.userId
	Env.Status = 0
	Env.Id = Env_id

	if err := Env.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("", MSG_OK)
}
