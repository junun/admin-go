package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type DomainInfoResource struct {
	Name    	string    	`form:"Name"`
	Channel    	string    	`form:"Channel"`
	StartTime   string    	`form:"StartTime"`
	EndTime     string    	`form:"EndTime"`
	Status		int			`form:"Status"`
	Desc 		string    	`form:"Desc"`
}

type DomainCertResource struct {
	Name    	string    	`form:"Name"`
	Did 		int 		`form:"Did"`
	Channel    	string    	`form:"Channel"`
	StartTime   string    	`form:"StartTime"`
	EndTime     string    	`form:"EndTime"`
	Status		int			`form:"Status"`
	Desc 		string    	`form:"Desc"`
}

// @Tags 域名管理
// @Description 域名列表
// @Summary  域名列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/info [get]
func GetDomainInfo(c *gin.Context)  {
	//if !middleware.PermissionCheckMiddleware(c,"domain-info-list") {
	//	util.JsonRespond(403, "请求资源被拒绝", "", c)
	//	return
	//}

	var domain []models.DomainInfo
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.DomainInfo{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&domain)
	models.DB.Model(&models.DomainInfo{}).Where(maps).Count(&count)

	data["lists"] = domain
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 域名管理
// @Description 域名新增
// @Summary  域名新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.DomainInfoResource true "域名信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/demo/info [post]
func AddDomainInfo(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-info-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DomainInfoResource
	var domain models.DomainInfo

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Domain Info Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 域名唯一性检查
	models.DB.Model(&models.DomainInfo{}).
		Where("name = ?", data.Name).Find(&domain)

	if domain.ID > 0 {
		util.JsonRespond(500, "重复的域名，请检查！", "", c)
		return
	}

	start, _	:= time.Parse(time.RFC3339, data.StartTime)
	end, _		:= time.Parse(time.RFC3339, data.EndTime)

	domain = models.DomainInfo{
		Name: data.Name,
		Channel: data.Channel,
		Status: data.Status,
		StartTime: start,
		EndTime: end,
		Desc: data.Desc}

	e := models.DB.Save(&domain).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加域名成功", "", c)
}

// @Tags 域名管理
// @Description 域名修改
// @Summary  域名修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "域名ID"
// @Param Data body admin.DomainInfoResource true "域名信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/info/{id} [put]
func PutDomainInfo(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-info-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DomainInfoResource
	var domain models.DomainInfo

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Domain Info Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 域名唯一性检查
	models.DB.Model(&models.DomainInfo{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&domain)

	if domain.ID > 0 {
		util.JsonRespond(500, "重复的域名，请检查！", "", c)
		return
	}

	models.DB.Find(&domain, c.Param("id"))

	start, _	:= time.Parse(time.RFC3339, data.StartTime)
	end, _		:= time.Parse(time.RFC3339, data.EndTime)

	domain.Name 	= data.Name
	domain.Channel 	= data.Channel
	domain.StartTime= start
	domain.EndTime 	= end
	domain.Status 	= data.Status
	domain.Desc 	= data.Desc

	e := models.DB.Save(&domain).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改域名成功", "", c)
}

// @Tags 域名管理
// @Description 域名删除
// @Summary  域名删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "域名ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/info/{id} [delete]
func DelDomainInfo(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-info-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	e := models.DB.Delete(models.DomainInfo{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除域名成功", "", c)
}


// @Tags 域名管理
// @Description 证书列表
// @Summary  证书列表
// @Produce  json
// @Param Authorization header string true "token"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/cert [get]
func GetDomainCret(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-cert-list") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}

	var cert []models.CertificateInfo
	var count int
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	models.DB.Model(&models.CertificateInfo{}).Where(maps).
		Offset(util.GetPage(c)).Limit(util.GetPageSize(c)).Find(&cert)
	models.DB.Model(&models.CertificateInfo{}).Where(maps).Count(&count)

	data["lists"] = cert
	data["count"] = count

	util.JsonRespond(200, "", data, c)
}

// @Tags 域名管理
// @Description 证书新增
// @Summary  证书新增
// @Produce  json
// @Param Authorization header string true "token"
// @Param Data body admin.DomainCertResource true "证书信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/demo/cert [post]
func AddDomainCret(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-cert-add") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DomainCertResource
	var cert models.CertificateInfo

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Add Cert Info Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 域名唯一性检查
	models.DB.Model(&models.DomainInfo{}).
		Where("name = ?", data.Name).Find(&cert)

	if cert.ID > 0 {
		util.JsonRespond(500, "重复的证书名，请检查！", "", c)
		return
	}

	start, _	:= time.Parse(time.RFC3339, data.StartTime)
	end, _		:= time.Parse(time.RFC3339, data.EndTime)

	cert = models.CertificateInfo{
		Name: data.Name,
		Did:data.Did,
		Channel: data.Channel,
		Status: data.Status,
		StartTime: start,
		EndTime: end,
		Desc: data.Desc}

	e := models.DB.Save(&cert).Error

	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "添加域名成功", "", c)
}

// @Tags 域名管理
// @Description 证书修改
// @Summary  证书修改
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "证书ID"
// @Param Data body admin.DomainCertResource true "证书信息"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/cert/{id} [put]
func PutDomainCret(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-cert-edit") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	var data DomainCertResource
	var cert models.CertificateInfo

	err := c.BindJSON(&data)
	if err != nil {
		util.JsonRespond(500, "Invalid Edit Domain Cert Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	// 域名唯一性检查
	models.DB.Model(&models.DomainInfo{}).
		Where("name = ?", data.Name).
		Where("id != ?", c.Param("id")).Find(&cert)

	if cert.ID > 0 {
		util.JsonRespond(500, "重复的证书名，请检查！", "", c)
		return
	}

	models.DB.Find(&cert, c.Param("id"))

	start, _	:= time.Parse(time.RFC3339, data.StartTime)
	end, _		:= time.Parse(time.RFC3339, data.EndTime)

	cert.Name 		= data.Name
	cert.Did		= data.Did
	cert.Channel 	= data.Channel
	cert.StartTime	= start
	cert.EndTime 	= end
	cert.Status 	= data.Status
	cert.Desc 		= data.Desc

	e := models.DB.Save(&cert).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "修改证书成功", "", c)
}

// @Tags 域名管理
// @Description 证书删除
// @Summary  证书删除
// @Produce  json
// @Param Authorization header string true "token"
// @Param id path int true "证书ID"
// @Success 200 {string} string {"code": 200, "message": "", "data": {}}
// @Failure 500 {string} string {"code": 500, "message": "", "data": {}}
// @Router /admin/domain/cert/{id} [delete]
func DelDomainCret(c *gin.Context)  {
	if !middleware.PermissionCheckMiddleware(c,"domain-cert-del") {
		util.JsonRespond(403, "请求资源被拒绝", "", c)
		return
	}
	e := models.DB.Delete(models.CertificateInfo{}, "id = ?", c.Param("id")).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	util.JsonRespond(200, "删除证书成功", "", c)
}
