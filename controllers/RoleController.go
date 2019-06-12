package controllers

import (
	"dtools/models/role"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

type RoleController struct {
	BaseController
}

func (self *RoleController) List() {
	self.Data["pageTitle"] = "角色管理"
	self.Display()
}

func (self *RoleController) Add() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "新增角色"
	self.Display()
}
func (self *RoleController) Edit() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "编辑角色"

	id, _ := self.GetInt("id", 0)
	role1, _ := role.RoleGetById(id)
	row := make(map[string]interface{})
	row["id"] = role1.Id
	row["role_name"] = role1.RoleName
	row["detail"] = role1.Detail
	self.Data["role"] = row

	//获取选择的树节点
	roleAuth, _ := role.RoleAuthGetById(id)
	authId := make([]int, 0)
	for _, v := range roleAuth {
		authId = append(authId, v.AuthId)
	}
	self.Data["auth"] = authId
	fmt.Println(authId)
	self.Display()
}

func (self *RoleController) AjaxSave() {
	role1 := new(role.Role)

	role1.RoleName = strings.TrimSpace(self.GetString("role_name"))
	role1.Detail = strings.TrimSpace(self.GetString("detail"))
	role1.CreateTime = time.Now().Unix()
	role1.UpdateTime = time.Now().Unix()
	role1.Status = 1


	auths := strings.TrimSpace(self.GetString("nodes_data"))
	role_id, _ := self.GetInt("id")
	if role_id == 0 {
		//新增
		role1.CreateTime = time.Now().Unix()
		role1.UpdateTime = time.Now().Unix()
		role1.CreateId = self.userId
		role1.UpdateId = self.userId

		if id, err := role.RoleAdd(role1); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		} else {
			ra := new(role.RoleAuth)
			authsSlice := strings.Split(auths, ",")
			for _, v := range authsSlice {
				aid, _ := strconv.Atoi(v)
				ra.AuthId = aid
				ra.RoleId = id
				role.RoleAuthAdd(ra)
			}
		}
		self.ajaxMsg("", MSG_OK)
	}
	//修改
	role1.Id = role_id
	role1.UpdateTime = time.Now().Unix()
	role1.UpdateId = self.userId
	if err := role1.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	} else {
		// 删除该角色权限
		role.RoleAuthDelete(role_id)
		ra := new(role.RoleAuth)
		authsSlice := strings.Split(auths, ",")
		for _, v := range authsSlice {
			aid, _ := strconv.Atoi(v)
			ra.AuthId = aid
			ra.RoleId = int64(role_id)
			role.RoleAuthAdd(ra)
		}

	}
	self.ajaxMsg("", MSG_OK)
}

func (self *RoleController) AjaxDel() {

	role_id, _ := self.GetInt("id")
	role, _ := role.RoleGetById(role_id)
	role.Status = 0
	role.Id = role_id
	role.UpdateTime = time.Now().Unix()

	if err := role.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	// 删除该角色权限
	//models.RoleAuthDelete(role_id)
	self.ajaxMsg("", MSG_OK)
}

func (self *RoleController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	roleName := strings.TrimSpace(self.GetString("roleName"))
	self.pageSize = limit
	//查询条件
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	if roleName != "" {
		filters = append(filters, "role_name__icontains", roleName)
	}
	result, count := role.RoleGetList(page, self.pageSize, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["role_name"] = v.RoleName
		row["detail"] = v.Detail
		row["create_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["update_time"] = beego.Date(time.Unix(v.UpdateTime, 0), "Y-m-d H:i:s")
		list[k] = row
	}
	self.AjaxList("成功", MSG_OK, count, list)
}
