package controllers

import (
	"dtools/models/admin"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dtools/libs"
	"dtools/models/role"
	"github.com/astaxie/beego"
)

type AdminController struct {
	BaseController
}

func (self *AdminController) List() {
	self.Data["pageTitle"] = "管理员管理"
	self.Display()
	//self.TplName = "admin/list.html"
}

func (self *AdminController) Add() {
	self.Data["pageTitle"] = "新增管理员"

	// 角色
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, _ := role.RoleGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		list[k] = row
	}

	self.Data["role"] = list

	self.Display()
}

func (self *AdminController) Edit() {
	self.Data["pageTitle"] = "编辑管理员"

	id, _ := self.GetInt("id", 0)
	Admin, _ := admin.AdminGetById(id)
	row := make(map[string]interface{})
	row["id"] = Admin.Id
	row["login_name"] = Admin.LoginName
	row["real_name"] = Admin.RealName
	row["phone"] = Admin.Phone
	row["email"] = Admin.Email
	row["role_ids"] = Admin.RoleIds
	self.Data["admin"] = row

	role_ids := strings.Split(Admin.RoleIds, ",")

	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, _ := role.RoleGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["checked"] = 0
		for i := 0; i < len(role_ids); i++ {
			role_id, _ := strconv.Atoi(role_ids[i])
			if role_id == v.Id {
				row["checked"] = 1
			}
			fmt.Println(role_ids[i])
		}
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		list[k] = row
	}
	self.Data["role"] = list
	self.Display()
}

func (self *AdminController) AjaxSave() {
	Admin_id, _ := self.GetInt("id")
	if Admin_id == 0 {
		Admin := new(admin.Admin)
		Admin.LoginName = strings.TrimSpace(self.GetString("login_name"))
		Admin.RealName = strings.TrimSpace(self.GetString("real_name"))
		Admin.Phone = strings.TrimSpace(self.GetString("phone"))
		Admin.Email = strings.TrimSpace(self.GetString("email"))
		Admin.RoleIds = strings.TrimSpace(self.GetString("roleids"))
		Admin.UpdateTime = time.Now().Unix()
		Admin.UpdateId = self.userId
		Admin.Status = 1

		// 检查登录名是否已经存在
		_, err := admin.AdminGetByName(Admin.LoginName)

		if err == nil {
			self.ajaxMsg("登录名已经存在", MSG_ERR)
		}
		//新增
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
		Admin.CreateTime = time.Now().Unix()
		Admin.CreateId = self.userId
		if _, err := admin.AdminAdd(Admin); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
		self.ajaxMsg("", MSG_OK)
	}

	Admin, _ := admin.AdminGetById(Admin_id)
	//修改
	Admin.Id = Admin_id
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.LoginName = strings.TrimSpace(self.GetString("login_name"))
	Admin.RealName = strings.TrimSpace(self.GetString("real_name"))
	Admin.Phone = strings.TrimSpace(self.GetString("phone"))
	Admin.Email = strings.TrimSpace(self.GetString("email"))
	Admin.RoleIds = strings.TrimSpace(self.GetString("roleids"))
	Admin.UpdateTime = time.Now().Unix()
	Admin.UpdateId = self.userId
	Admin.Status = 1

	resetPwd, _ := self.GetInt("reset_pwd")
	if resetPwd == 1 {
		pwd, salt := libs.Password(4, "")
		Admin.Password = pwd
		Admin.Salt = salt
	}
	if err := Admin.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg(strconv.Itoa(resetPwd), MSG_OK)
}

func (self *AdminController) AjaxDel() {

	Admin_id, _ := self.GetInt("id")
	status := strings.TrimSpace(self.GetString("status"))
	if Admin_id == 1 {
		self.ajaxMsg("超级管理员不允许操作", MSG_ERR)
	}

	Admin_status := 0
	if status == "enable" {
		Admin_status = 1
	}
	Admin, _ := admin.AdminGetById(Admin_id)
	Admin.UpdateTime = time.Now().Unix()
	Admin.Status = Admin_status
	Admin.Id = Admin_id

	if err := Admin.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	self.ajaxMsg("操作成功", MSG_OK)
}

func (self *AdminController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	realName := strings.TrimSpace(self.GetString("realName"))

	StatusText := make(map[int]string)
	StatusText[0] = "<font color='red'>禁用</font>"
	StatusText[1] = "正常"

	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	//
	if realName != "" {
		filters = append(filters, "real_name__icontains", realName)
	}
	result, count := admin.AdminGetList(page, self.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["login_name"] = v.LoginName
		row["real_name"] = v.RealName
		row["phone"] = v.Phone
		row["email"] = v.Email
		row["role_ids"] = v.RoleIds
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		row["status"] = v.Status
		row["status_text"] = StatusText[v.Status]
		list[k] = row
	}
	self.AjaxList("成功", MSG_OK, count, list)
}
