package controllers

import (
	"dtools/models/auth"
	"dtools/utils"
	"fmt"
	"github.com/patrickmn/go-cache"
	"strconv"
	"strings"
	"time"
)

type AuthController struct {
	BaseController
}

func (self *AuthController) Index() {

	self.Data["pageTitle"] = "权限因子"
	self.Display()
}

func (self *AuthController) List() {
	self.Data["zTree"] = true //引入ztreecss
	self.Data["pageTitle"] = "权限因子"
	self.Display()
}

//获取全部节点
func (self *AuthController) GetNodes() {
	filters := make([]interface{}, 0)
	filters = append(filters, "status", 1)
	result, count := auth.AuthGetList(1, 1000, filters...)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["pId"] = v.Pid
		row["name"] = v.AuthName
		row["open"] = true
		list[k] = row
	}

	self.AjaxList("成功", MSG_OK, count, list)
}

//获取一个节点
func (self *AuthController) GetNode() {
	id, _ := self.GetInt("id")
	result, _ := auth.AuthGetById(id)
	// if err == nil {
	// 	self.ajaxMsg(err.Error(), MSG_ERR)
	// }
	row := make(map[string]interface{})
	row["id"] = result.Id
	row["pid"] = result.Pid
	row["auth_name"] = result.AuthName
	row["auth_url"] = result.AuthUrl
	row["sort"] = result.Sort
	row["is_show"] = result.IsShow
	row["icon"] = result.Icon

	fmt.Println(row)

	self.AjaxList("成功", MSG_OK, 0, row)
}

//新增或修改
func (self *AuthController) AjaxSave() {
	aut1 := new(auth.Auth)
	aut1.UserId = self.userId
	aut1.Pid, _ = self.GetInt("pid")
	aut1.AuthName = strings.TrimSpace(self.GetString("auth_name"))
	aut1.AuthUrl = strings.TrimSpace(self.GetString("auth_url"))
	aut1.Sort, _ = self.GetInt("sort")
	aut1.IsShow, _ = self.GetInt("is_show")
	aut1.Icon = strings.TrimSpace(self.GetString("icon"))
	aut1.UpdateTime = time.Now().Unix()

	aut1.Status = 1

	id, _ := self.GetInt("id")
	if id == 0 {
		//新增
		aut1.CreateTime = time.Now().Unix()
		aut1.CreateId = self.userId
		aut1.UpdateId = self.userId
		if _, err := auth.AuthAdd(aut1); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	} else {
		aut1.Id = id
		aut1.UpdateId = self.userId
		if err := aut1.Update(); err != nil {
			self.ajaxMsg(err.Error(), MSG_ERR)
		}
	}
	utils.Che.Set("menu"+strconv.Itoa(self.user.Id), nil, cache.DefaultExpiration)
	self.ajaxMsg("", MSG_OK)
}

//删除
func (self *AuthController) AjaxDel() {
	id, _ := self.GetInt("id")
	auth, _ := auth.AuthGetById(id)
	auth.Id = id
	auth.Status = 0
	if err := auth.Update(); err != nil {
		self.ajaxMsg(err.Error(), MSG_ERR)
	}
	utils.Che.Set("menu"+strconv.Itoa(self.user.Id), nil, cache.DefaultExpiration)
	self.ajaxMsg("", MSG_OK)
}
