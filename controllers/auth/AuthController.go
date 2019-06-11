package auth

import (
	"dtools/controllers"
	"dtools/models/auth"
)

type AuthController struct {
	controllers.BaseController
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
	self.AjaxList("成功", controllers.MSG_OK, count, list)
}
