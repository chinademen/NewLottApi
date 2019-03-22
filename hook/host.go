package hook

// "NewLottApi/models"
// "common"
// "strings"

/**
 * 检测域名是否合法域名(不包含www.)
 * @params map  请求的参数和host
 * @return
 * int 状态
 * string 信息
 */
func CheckHost(host string) (int, string) {
	status := 203
	msg := "域名不合法:" + host
	// if len(host) < 1 {
	// 	status = 204
	// 	msg = "域名不能为空"
	// 	return status, msg
	// }

	// //获取顶级域名
	// domain := common.GetTopDomain(host)
	// domainChk := models.Host.Get(domain)
	// if len(domainChk) > 0 || strings.Count(host, "127.0.0.1") > 0 {
	// 	status = 200
	// 	msg = "域名合法"
	// }

	return status, msg
}
