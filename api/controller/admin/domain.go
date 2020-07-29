package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	"strings"
	"time"
)

type DomainInfoResource struct {
	Name    	string    	`form:"Name"`
	IsCert		int			`form:"IsCert"`
	CertName 	string    	`form:"CertName"`
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
	data.CertName = strings.TrimSpace(data.CertName)
	// 证书唯一性检查
	if data.IsCert == 1 {
		models.DB.Model(&models.DomainInfo{}).
			Where("cert_name = ?", data.CertName ).Find(&domain)

		if domain.ID > 0 {
			util.JsonRespond(500, "重复的证书检查，请检查！", "", c)
			return
		}
	}

	domain = models.DomainInfo{
		Name: data.Name,
		Status: data.Status,
		IsCert: data.IsCert,
		Desc: data.Desc}

	if data.IsCert == 1  {
		if data.CertName == "" {
			util.JsonRespond(500, "有证书情况必须设置检测证书有效性的二级域名！", "", c)
			return
		}

		domain.CertName = data.CertName
	}

	//e := models.DB.Save(&domain).Error
	//if e != nil {
	//	util.JsonRespond(500, e.Error(), "", c)
	//	return
	//}

	e := models.DB.Create(&domain).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}


	CheckDomainMaster(domain)
	if data.IsCert == 1  {
		CheckDomainCert(domain)
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
		util.JsonRespond(400, "Invalid Edit Domain Info Data", "", c)
		return
	}

	data.Name = strings.TrimSpace(data.Name)
	data.CertName = strings.TrimSpace(data.CertName)

	// 证书唯一性检查
	if data.IsCert == 1 {
		models.DB.Model(&models.DomainInfo{}).
			Where("cert_name = ?", data.CertName ).
			Where("id != ?", c.Param("id")).Find(&domain)

		if domain.ID > 0 {
			util.JsonRespond(500, "重复的证书检查，请检查！", "", c)
			return
		}
	}


	models.DB.Find(&domain, c.Param("id"))
	domainName 		:= domain.Name
	certName		:= domain.CertName

	domain.Name 	= data.Name
	domain.IsCert	= data.IsCert
	domain.Status 	= data.Status
	domain.Desc 	= data.Desc

	if data.IsCert == 1  {
		if data.CertName == "" {
			util.JsonRespond(500, "有证书情况必须设置检测证书有效性的二级域名！", "", c)
			return
		}
	}

	if domain.CertName != data.CertName {
		domain.CertName = data.CertName
	}

	e := models.DB.Save(&domain).Error
	if e != nil {
		util.JsonRespond(500, e.Error(), "", c)
		return
	}

	if data.Name != domainName {
		CheckDomainMaster(domain)
	}

	if data.CertName != certName {
		CheckDomainCert(domain)
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


// 定时任务 每天检查域名和证书是否到期
// 到期前一个月开始每天发送通知，提醒续费或者更新证书
func CheckDomainAndCret()   {
	var domains []models.DomainInfo

	models.DB.Model(&models.DomainInfo{}).Where("status=1").Find(&domains)
	if len(domains) == 0 {
		return
	}

	if len(domains) > 0 {
		for _,domain:= range domains {
			// 检验域名
			CheckDomainMaster(domain)

			// 检测证书
			if domain.IsCert == 1 && domain.CertName != "" {
				CheckDomainCert(domain)
			}
		}
	}
}

func CheckDomainCert(domain models.DomainInfo) {
	cert, e := util.ParseRemoteCertificate(domain.CertName + ":443",10)
	if e!=nil {
		models.MakeNotify(1, 2,"定时任务异常" + e.Error(),"域名" + domain.Name + "证书有效性信息检查异常，请检测！","" )
	}

	now 			:= time.Now()
	nowAddOneMonth 	:= now.AddDate(0, 1, 0)

	if (cert.NotAfter.Before(nowAddOneMonth)) {
		// 执行消息通知逻辑
		// text message
		e := models.DingtalkSentChannel(0, "@" + models.DingUser + "SSl证书" + domain.Name + "快要到期了，请及时处理", models.DingList, false)
		if e != nil {
			models.MakeNotify(1, 2,"发送通知异常","SSl证书" + domain.Name + "快要到期了，请及时处理","" )
		}
	}

	if domain.CertEndTime != cert.NotAfter {
		e := models.DB.Model(&domain).Updates(map[string]interface{}{"cert_end_time": cert.NotAfter}).Error
		if e != nil {
			models.MakeNotify(1, 2,"更新域名证书信息异常" + e.Error(),"域名" + domain.Name + "更新域名证书信息异常，请及时处理","" )
		}
	}
}

func CheckDomainMaster(domain models.DomainInfo)  {
	resWho, e 	:= whois.Whois(domain.Name)
	if e != nil {
		//// 启用备用检测方式
		//CheckDomainBackup(domain)
		models.MakeNotify(1, 2,"定时任务异常" + e.Error(),"域名" + domain.Name + "查询异常！","" )
		return
	}

	resParse, e := whoisparser.Parse(resWho)
	if e != nil {
		models.MakeNotify(1, 2,"定时任务异常" + e.Error(),"域名" + domain.Name + "有效性信息解析错误请检测！","" )

		//// 启用备用检测方式
		//CheckDomainBackup(domain)
		return
	}

	//if resParse.Domain.ExpirationDate == "0001-01-01 00:00:00 +0000 UTC" {
	//	// 启用备用检测方式
	//	CheckDomainBackup(domain)
	//	return
	//}

	handDomainParse(domain, resParse.Domain.ExpirationDate)
}

func handDomainParse(domain models.DomainInfo, timestr string)  {
	now 			:= time.Now()
	nowAddOneMonth 	:= now.AddDate(0, 1, 0)
	end, _  := time.Parse(time.RFC3339, timestr)
	if (end.Before(nowAddOneMonth)) {
		// 执行消息通知逻辑
		// text message
		e := models.DingtalkSentChannel(0, "@" + models.DingUser + "域名" + domain.Name + "快要到期了，请及时处理", models.DingList, false)
		if e != nil {
			models.MakeNotify(1, 2,"发送通知异常" + e.Error(),"域名" + domain.Name + "快要到期了，请及时处理","" )
		}
	}

	if domain.DomainEndTime != end {
		domain.DomainEndTime = end
		e := models.DB.Model(&domain).Updates(map[string]interface{}{"domain_end_time": end}).Error
		//e 	:= models.DB.Model(&models.DomainInfo{}).Where("id = ?", domain.ID).Updates(map[string]interface{}{"domain_end_time": end}).Error
		//e := models.DB.Table("domain_info").Where("id = ?", domain.ID).Updates(map[string]interface{}{"domain_end_time": end}).Error
		// 不要使用该方式更新，尤其在证书检查那边，不然会造成前后内容覆盖
		//e := models.DB.Save(&domain).Error
		if e != nil {
			models.MakeNotify(1, 2,"更新域名信息异常" + e.Error(),"域名" + domain.Name + "更新域名信息异常，请及时处理","" )
		}
	}
}