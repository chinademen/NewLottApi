package controllers

import (
	"NewLottApi/dao"
	"NewLottApi/log"
	"NewLottApi/models"
	"common"
	"common/ext/redisClient"
	"fmt"
	lotteryJobsModels "lotteryJobs/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//UserController 用户类
type UserController struct {
	MainController
}

func (u *UserController) Finish() {
	apilogObj := u.Ctx.Input.GetData("apiLog")
	if apilog, ok := apilogObj.(*dao.ApiLogs); ok == true {
		apilog.EndTime = time.Now()
		dao.AddApiLogs(apilog)
	}
	u.Ctx.Input.SetData("paramMap", nil)
	u.Ctx.Input.SetData("apiLog", nil)
	u.Ctx.Input.SetData("merchantID", nil)
}

// @Title InitUser
// @Description 初始化用户
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数; username, password, prize_group, ip"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /init.do [post]
func (c *UserController) InitUser() {

	reStatus := 101
	reMsg := "注册失败"
	d := map[string]interface{}{}

	filterUserMap := c.FilterUser()

	clientIp := filterUserMap["client_ip"] //商户服务器ip

	//锁定，每个ip5秒内只运行一个注册请求
	rKey := fmt.Sprintf("user_reg:%s", clientIp)
	rLock := redisClient.Redis.StringWrite(rKey, clientIp, 5)
	if rLock <= 0 {
		reStatus = 1000
		reMsg = "请求太频繁"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	merchantId := filterUserMap["merchant_id"] //商户id
	sIdentity := filterUserMap["identity"]     //商户唯一标识

	username := filterUserMap["username"]      //接入方平台用户名
	password := filterUserMap["password"]      //接入方平台用户密码
	prizeGroup := filterUserMap["prize_group"] //对应用户的奖金组设定
	userIP := filterUserMap["ip"]              //用户IP

	if len(prizeGroup) < 1 || len(userIP) < 1 {
		if len(prizeGroup) < 1 {
			reStatus = 102
			reMsg = "奖金组不能为空"
		} else {

			reStatus = 103
			reMsg = "客户IP不能为空"
		}

		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//验证帐号密码规则
	reStatus, reMsg = models.ChkUserAndPwd(username, password)
	if reStatus != 200 {
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//帐号是否存在
	reStatus, reMsg = models.IsKyAccount(sIdentity, username)
	if reStatus != 200 {
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//判断奖金组
	prizeGroupsRow := models.PrizeGroups.RGetOneByName(prizeGroup)
	if len(prizeGroupsRow) == 0 {
		c.RenderJson(105, "奖金组不存在", d) //将数据装载到json返回值
	}

	//注册帐号 start
	newDatetime := common.GetNowDatetime(common.DATE_FORMAT_YMDHIS)
	newDateUnixStr := common.InterfaceToString(time.Now().Unix())

	//用户主表
	userMap := map[string]string{}
	userMap["username"] = username
	userMap["password"] = models.EncodeString(password)
	userMap["fund_password"] = ""                   //资金密码,注册不需要
	userMap["account_id"] = newDateUnixStr          //账户id, 暂时用时间戳 之后update
	userMap["merchant_id"] = merchantId             //商户ID
	userMap["prize_group"] = prizeGroup             //奖金组
	userMap["blocked"] = "0"                        //1=冻结
	userMap["realname"] = ""                        //真实姓名
	userMap["nickname"] = ""                        //昵称
	userMap["email"] = ""                           //邮件
	userMap["mobile"] = ""                          //电话号码
	userMap["is_tester"] = "0"                      //1=测试
	userMap["bet_multiple"] = "1"                   //投注倍数
	userMap["bet_coefficient"] = "1"                //投注模式
	userMap["login_ip"] = ""                        //最后登录ip
	userMap["register_ip"] = userIP                 //注册ip
	userMap["token"] = ""                           //用户token
	userMap["signin_at"] = newDatetime              //登录时间
	userMap["activated_at"] = "0000-00-00 00:00:00" //活跃时间，投注
	userMap["register_at"] = newDatetime            //注册时间
	userMap["deleted_at"] = "0000-00-00 00:00:00"   //删除时间
	userMap["created_at"] = newDatetime             //数据库创建时间
	userMap["updated_at"] = newDatetime             //更新时间

	//1 开启事务
	mDbBeg, mDbBegErr := models.Mdb.Begin()
	if mDbBegErr != nil {

		mDbBeg.Rollback() // 回滚事务

		reStatus = 701
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//2 插入用户主表
	usQuery, usQueryErr := mDbBeg.Exec(models.Users.GetAddOnlySql(userMap))
	if usQueryErr != nil {

		mDbBeg.Rollback() // 回滚事务

		reStatus = 702
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//返回影响的行数
	usInt := 0
	usRowA, usRowAErr := usQuery.RowsAffected()
	if usRowAErr == nil {
		usInt = int(usRowA)
	}

	if usInt < 1 {
		mDbBeg.Rollback() // 回滚事务

		reStatus = 703
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//返回插入id
	usInId, usInIdErr := usQuery.LastInsertId()
	if usInIdErr != nil {
		mDbBeg.Rollback() // 回滚事务

		reStatus = 704
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	usId := strconv.FormatInt(usInId, 10)

	//用户账户表
	accountMap := map[string]string{}
	accountMap["merchant_id"] = userMap["merchant_id"] //商户ID
	accountMap["user_id"] = usId                       //用户id
	accountMap["username"] = userMap["username"]       //用户名
	accountMap["is_tester"] = userMap["is_tester"]     //1=测试
	accountMap["balance"] = "0"                        //总额度
	accountMap["frozen"] = "0"                         //冻结额度
	accountMap["available"] = "0"                      //可用额度
	accountMap["status"] = "1"                         //状态:1=正常，-1=删除
	accountMap["locked"] = "0"                         //1=冻结
	accountMap["created_at"] = newDatetime             //创建时间
	accountMap["updated_at"] = newDatetime             //更新时间
	accountMap["backup_made_at"] = newDatetime         //数据库更新时间

	//3 插入用户账户表
	accQuery, accQueryErr := mDbBeg.Exec(models.Accounts.GetAddOnlySql(accountMap))
	if accQueryErr != nil {
		mDbBeg.Rollback()

		reStatus = 705
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//返回影响的行数
	accInt := 0
	accRowA, accRowAErr := accQuery.RowsAffected()
	if accRowAErr == nil {
		accInt = int(accRowA)
	}

	if accInt < 1 {
		mDbBeg.Rollback()

		reStatus = 706
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//返回插入id
	accInId, accInIdErr := accQuery.LastInsertId()
	if accInIdErr != nil {
		mDbBeg.Rollback()

		reStatus = 707
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	accountsId := strconv.FormatInt(accInId, 10)
	usUpdateMap := map[string]string{
		"account_id": accountsId,
	}
	usUpdWhere := fmt.Sprintf("id='%s'", usId)

	//4 更新用户表账户id
	usUPQue, usUPQueErr := mDbBeg.Exec(models.Users.DbGetUpdateSql(usUpdateMap, usUpdWhere))
	if usUPQueErr != nil {
		mDbBeg.Rollback()

		reStatus = 708
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//返回影响的行数
	usUpdInt := 0
	usUPRowA, usUPRowAErr := usUPQue.RowsAffected()
	if usUPRowAErr == nil {
		usUpdInt = int(usUPRowA)
	}

	if usUpdInt < 1 {
		mDbBeg.Rollback()

		reStatus = 709
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	bSucc := true
	//得到 type=1 or type =2 的彩种信息
	oLotteries := models.GetGroupPrizeLottery()
	for _, lottV := range oLotteries {

		iGroupId := models.PrizeGroups.GetGroupId(lottV["series_id"], prizeGroup)
		if len(iGroupId) == 0 {
			continue
		}

		uPrizeMap := map[string]string{}
		uPrizeMap["merchant_id"] = userMap["merchant_id"] //商户id
		uPrizeMap["user_id"] = usId                       //用户ID
		uPrizeMap["username"] = userMap["username"]       //用户名
		uPrizeMap["series_id"] = lottV["series_id"]       //系列id
		uPrizeMap["lottery_id"] = lottV["id"]             //彩种
		uPrizeMap["group_id"] = iGroupId                  //组
		uPrizeMap["prize_group"] = prizeGroup             //奖金组
		uPrizeMap["classic_prize"] = prizeGroup           //经典奖金
		uPrizeMap["valid"] = "1"                          //有效
		uPrizeMap["is_agent"] = "0"                       //是否代理0: 普通用户, 1: 代理
		uPrizeMap["created_at"] = newDatetime
		uPrizeMap["updated_at"] = newDatetime

		//5 插入用户奖金组设置
		uPrizeSql := models.Mdb.GetInsertTrueSql(models.UserPrizeSets.Table.TableName, uPrizeMap)
		uPrizeQuery, uPrizeErr := mDbBeg.Exec(uPrizeSql)
		if uPrizeErr != nil {
			bSucc = false
			break
		}

		//返回影响的行数
		uPrizeInt, uPrizeRowErr := uPrizeQuery.RowsAffected()
		if uPrizeRowErr != nil || uPrizeInt != 1 {

			bSucc = false
			break
		}
	}

	if !bSucc {
		mDbBeg.Rollback() // 回滚事务

		reStatus = 710
		reMsg = "注册失败"
		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	//6 提交事务
	DBComtErr := mDbBeg.Commit()
	if DBComtErr == nil {

		reStatus = 200
		reMsg = "注册成功"
		d["userid"] = usId
		d["usename"] = username

		//记录日志
		//go hook.WebLog(c.Ctx, "user", "reg", reMsg)

		c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值
	}

	mDbBeg.Rollback()

	reStatus = 711
	reMsg = "注册失败"
	c.RenderJson(reStatus, reMsg, d) //将数据装载到json返回值

}

// @Title Login
// @Description 用户登录
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /login.do [post]
func (c *UserController) Login() {

	d := map[string]interface{}{}

	filterUserMap := c.FilterUser()

	sMerchantId := filterUserMap["merchant_id"] //商户id
	sIdentity := filterUserMap["identity"]      //商户唯一标识
	sUsername := filterUserMap["username"]      //用户名
	sPassword := filterUserMap["password"]      //用户密码
	sDevice := filterUserMap["device"]          //用户终端设备标志： 1-PC 2-手机
	sUserIP := filterUserMap["ip"]              //用户IP
	sBrowser := filterUserMap["browser"]        //浏览器类型

	//判斷終端&Ip
	if len(sDevice) < 1 {
		c.RenderJson(301, "设备标志不能为空", d)
	}

	if sDevice != "1" && sDevice != "2" {
		c.RenderJson(302, "用户终端有误", d)
	}

	if len(sUserIP) < 1 {
		c.RenderJson(303, "客户IP不能为空", d)
	}

	if len(sUsername) < 1 {
		c.RenderJson(304, "用戶名不能爲空", d)
	}

	if len(sPassword) < 1 {
		c.RenderJson(305, "用戶密碼不能爲空", d)
	}

	/////////////////////
	//验证帐号密码是否正确//
	/////////////////////
	reStatus, reMsg, userRow := models.ChkLogAccount(sIdentity, sUsername, sPassword)
	if reStatus != 200 {
		log.LogsWithFileName("", "Login-Err", reMsg+"|"+sUsername+":"+sPassword, log.USER)
		c.RenderJson(reStatus, reMsg, d)
	}

	//token
	iInt, sToken := models.UserToken.CreateToken(sMerchantId, userRow["id"], sUserIP, sDevice, sBrowser)
	if iInt < 1 {
		c.RenderJson(306, "获取token失败", d)
	}

	newDatetime := common.GetNowDatetime(common.DATE_FORMAT_YMDHIS)

	userUpdate := map[string]string{}
	userUpdate["login_ip"] = sUserIP       //最后登录ip
	userUpdate["signin_at"] = newDatetime  //登录时间
	userUpdate["updated_at"] = newDatetime //更新时间

	if models.Users.DbUpdate(userUpdate, fmt.Sprintf("id = '%s' ", userRow["id"])) < 1 {
		reStatus = 307
		reMsg = "登录失败"
		c.RenderJson(reStatus, reMsg, d)
	}

	d["token"] = sToken
	d["game_lobby_url"] = models.GameLobbyUrl
	c.RenderJson(200, "登录成功", d)
}

// @Title Balance
// @Description 用户额度查询,author(leon)
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /balance.do [post]
func (u *UserController) Balance() {
	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			u.RenderJson(500, "系统错误", err)
		}

	}()

	iStatus := 200
	sMsg := "success"
	var res interface{}

	filterUserMap := u.FilterUser()

	sUsername := filterUserMap["username"]     //用户名
	sPassword := filterUserMap["password"]     //用户密码
	sMechantId := filterUserMap["merchant_id"] //用户密码

	if len(sUsername) < 1 || len(sPassword) < 1 || len(sMechantId) < 1 {
		iStatus = 7001
		sMsg = "参数为空或者参数错误"
		u.RenderJson(iStatus, sMsg, res)
	}

	///////////
	//判斷用戶//
	//////////
	reStatus, reMsg, reResult, oUser := u.CheckUser(sUsername, sPassword, sMechantId)
	if reStatus != 200 {
		u.RenderJson(reStatus, reMsg, reResult)
	}

	//判斷用戶賬戶和餘額
	oAccount, err := dao.GetAccountsById(int(oUser.AccountId))
	if oAccount.Id == 0 || err != nil {
		iStatus = 7002
		sMsg = "用户账户出错"
		u.RenderJson(iStatus, sMsg, res)
	}

	if oAccount.Locked == 1 {
		iStatus = 7003
		sMsg = "用户账户已冻结"
		u.RenderJson(iStatus, sMsg, res)
	}

	res = map[string]float64{
		"balance":   oAccount.Balance,
		"available": oAccount.Available,
		"fronzen":   oAccount.Frozen,
	}
	u.RenderJson(iStatus, sMsg, res)
}

// @Title Transfer
// @Description 转移额度订单,author(leon)
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /transfer.do [post]
func (u *UserController) Transfer() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			u.RenderJson(500, "系统错误", err)
		}

	}()

	iStatus := 200
	sMsg := "success"
	var res interface{}

	filterUserMap := u.FilterUser()

	sUsername := filterUserMap["username"] //用户名
	sPassword := filterUserMap["password"] //用户密码
	sAmount := filterUserMap["amount"]     //轉移金額
	sType := filterUserMap["type"]         //操作類型
	sMerchantId := filterUserMap["merchant_id"]
	iMerchantId, _ := strconv.Atoi(sMerchantId)

	if len(sUsername) < 1 || len(sPassword) < 1 || len(sAmount) < 1 || len(sType) < 1 {
		iStatus = 7004
		sMsg = "参数为空或者参数错误"
		u.RenderJson(iStatus, sMsg, res)
	}

	if debug {
		fmt.Println("-------------传递的参数---------------")
		fmt.Println("sUsername-->", sUsername)
		fmt.Println("sPassword-->", sPassword)
		fmt.Println("sAmount-->", sAmount)
		fmt.Println("sType-->", sType)
		fmt.Println("sMerchantId-->", sMerchantId)
	}

	///////////
	//判斷用戶//
	//////////
	reStatus, reMsg, reResult, oUser := u.CheckUser(sUsername, sPassword, sMerchantId)
	if reStatus != 200 {
		u.RenderJson(reStatus, reMsg, reResult)
	}

	if debug {
		fmt.Println("-------------oUser---------------")
		fmt.Println("oUser.Id-->", oUser.Id)
		fmt.Println("oUser.Username-->", oUser.Username)
	}

	if oUser.Blocked == 1 {
		iStatus = 7005
		sMsg = "用戶账户已冻结"
		u.RenderJson(iStatus, sMsg, res)
	}

	///////////
	//判斷金額//
	///////////
	sAmount = strings.Trim(sAmount, " ")
	fAmount, _ := strconv.ParseFloat(sAmount, 64)
	if fAmount <= 0 {
		iStatus = 7006
		sMsg = "金额错误"
		u.RenderJson(iStatus, sMsg, res)
	}

	//判斷用戶賬戶和餘額
	oAccount, err := dao.GetAccountsById(int(oUser.AccountId))
	if oAccount.Id == 0 || err != nil {
		iStatus = 7007
		sMsg = "用戶账户出错"
		u.RenderJson(iStatus, sMsg, res)
	}

	if debug {
		fmt.Println("-------------oAccount---------------")
		fmt.Println("oAccount.Id-->", oAccount.Id)
		fmt.Println("oAccount.Username-->", oAccount.Username)
	}

	if oAccount.Locked == 1 {
		iStatus = 7008
		sMsg = "用戶賬戶已凍結"
		u.RenderJson(iStatus, sMsg, res)
	}

	//////////////
	//判斷操作類型//
	//////////////
	iType, _ := strconv.Atoi(sType)
	if iType != models.TransferIn && iType != models.TransferOut {
		iStatus = 7010
		sMsg = "操作类型有误"
		u.RenderJson(iStatus, sMsg, res)
	}

	////////
	//轉帳//
	///////
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		iStatus = 7011
		sMsg = "系统异常"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	//初始化張變記錄
	var oTransactions = new(dao.Transactions)
	if iType == models.TransferIn {
		oTransactions.IsIncome = 0
		oTransactions.TypeId = 1
	} else {
		oTransactions.IsIncome = 1
		oTransactions.TypeId = 2

		if oAccount.Available < fAmount {
			iStatus = 7009
			sMsg = "用户账户余额不足"
			u.RenderJson(iStatus, sMsg, res)
		}
	}

	//獲取賬變類型
	oTransactionsType := models.TransactionTypes.GetInfo(int(oTransactions.TypeId))
	if oTransactionsType.Available == 0 || oTransactionsType.Balance == 0 || oTransactionsType.Frozen == 0 {
		oTransactionsType, _ = dao.GetTransactionTypesById(int(oTransactions.TypeId))
	}
	if oTransactionsType.Id == 0 {
		iStatus = 7012
		sMsg = "帐变异常"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	oTransactions.AccountId = uint64(oAccount.Id)
	oTransactions.PreviousAvailable = oAccount.Available
	oTransactions.PreviousBalance = oAccount.Balance
	oTransactions.PreviousFrozen = oAccount.Frozen

	if debug {
		fmt.Println("oTransactionsType-->", oTransactionsType)
		fmt.Println("oTransactionsType.Id-->", oTransactionsType.Id)
		fmt.Println("oTransactionsType.Available-->", oTransactionsType.Available)
		fmt.Println("oTransactionsType.Balance-->", oTransactionsType.Balance)
		fmt.Println("oTransactionsType.Frozen-->", oTransactionsType.Frozen)
	}

	oTransactions.Available = oAccount.Available + (fAmount * float64(oTransactionsType.Available))
	oTransactions.Balance = oAccount.Balance + (fAmount * float64(oTransactionsType.Balance))
	oTransactions.Frozen = oAccount.Frozen + (fAmount * float64(oTransactionsType.Frozen))

	oTransactions.Amount = fAmount
	oTransactions.IsTester = oUser.IsTester
	oTransactions.UserId = uint64(oUser.Id)
	oTransactions.Username = oUser.Username
	oTransactions.MerchantId = uint(iMerchantId)

	oTransactions = models.TransactionList.MakeSeriesNumber(oTransactions)
	oTransactions = models.TransactionList.SaveKey(oTransactions)

	oTransactions.CreatedAt = time.Now().Format(common.DATE_FORMAT_YMDHIS)
	oTransactions.UpdatedAt = time.Now().Format(common.DATE_FORMAT_YMDHIS)
	iTransactionId, err := dao.AddTransactions(o, oTransactions)
	if iTransactionId < 1 || err != nil {
		o.Rollback()
		iStatus = 7013
		sMsg = "生成系统异常失敗"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	if debug {
		fmt.Println("available-->", oTransactions.Available)
		fmt.Println("balance-->", oTransactions.Balance)
		fmt.Println("frozen-->", oTransactions.Frozen)
	}

	//更新用戶賬戶餘額
	oAccount.Available = oTransactions.Available
	oAccount.Balance = oTransactions.Balance
	oAccount.Frozen = oTransactions.Frozen

	if debug {
		fmt.Println("oAccount.Available-->", oAccount.Available)
		fmt.Println("oAccount.Balance-->", oAccount.Balance)
		fmt.Println("oAccount.Frozen-->", oAccount.Frozen)
	}
	err = dao.UpdateAccountsById(o, oAccount)
	if err != nil {
		o.Rollback()
		iStatus = 7014
		sMsg = "用戶账户变更失敗"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	//生成轉帳信息
	var oTransferInfo = new(dao.TransferInfo)
	oTransferInfo.UserId = uint(oUser.Id)
	oTransferInfo.BillNo = oTransactions.SerialNumber
	oTransferInfo.Amount = fAmount
	oTransferInfo.Status = 1
	oTransferInfo.TypeId = uint8(oTransactions.TypeId)
	oTransferInfo.MerchantId = uint(iMerchantId)
	oTransferInfo.OrderNumber = strconv.Itoa(oTransactions.Id)
	oTransferInfo.AcceptedAt = time.Now()
	iTransferInfoId, err := dao.AddTransferInfo(oTransferInfo)
	if iTransferInfoId < 1 || err != nil {
		o.Rollback()
		iStatus = 7015
		sMsg = "生成转账记录失敗"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	err = o.Commit()
	if err != nil {
		iStatus = 7016
		sMsg = "系统异常"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	res = map[string]interface{}{
		"balance":      oAccount.Balance,
		"available":    oAccount.Available,
		"frozen":       oAccount.Frozen,
		"order_number": oTransactions.Id,
	}
	u.RenderJson(iStatus, sMsg, res)

}

// @Title Transferinfo
// @Description 根据商户传递过来的用户信息从彩票中心查询额度转入订单的状态信息,editor(leon)
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /transferinfo.do [post]
func (u *UserController) Transferinfo() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			u.RenderJson(500, "系统错误", err)
		}

	}()

	iStatus := 200
	sMsg := "success"
	var res interface{}

	filterUserMap := u.FilterUser()

	sUsername := filterUserMap["username"]        //用户名
	sPassword := filterUserMap["password"]        //用户密码
	sOrderNumber := filterUserMap["order_number"] //轉移金額
	sMerchantId := filterUserMap["merchant_id"]

	if len(sUsername) < 1 || len(sPassword) < 1 || len(sOrderNumber) < 1 || len(sMerchantId) < 1 {
		iStatus = 7017
		sMsg = "参数为空或者参数错误"
		u.RenderJson(iStatus, sMsg, res)
	}

	///////////
	//判斷用戶//
	//////////
	reStatus, reMsg, reResult, oUser := u.CheckUser(sUsername, sPassword, sMerchantId)
	if reStatus != 200 {
		u.RenderJson(reStatus, reMsg, reResult)
	}

	//根據訂單號獲取訂單信息
	oTransferInfo := models.TransferInfo.GetInfo(fmt.Sprintf("%s", oUser.Id), sOrderNumber)
	if oTransferInfo.Id == 0 {
		iStatus = 7018
		sMsg = "找不到该订单:" + sOrderNumber
		u.RenderJson(iStatus, sMsg, res)
	}

	res = map[string]interface{}{
		"status":       oTransferInfo.Status,
		"bill_no":      oTransferInfo.BillNo,
		"order_number": sOrderNumber,
		"amount":       oTransferInfo.Amount,
	}
	u.RenderJson(iStatus, sMsg, res)

}

// @Title Edit
// @Description 根据商户传递过来的用户信息变更用户在彩票中心的密码或奖金组,editor(leon)
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /edit.do [post]
func (u *UserController) Edit() {

	//捕捉錯誤
	defer func() {
		if err := recover(); err != nil {
			u.RenderJson(500, "系统错误", err)
		}

	}()

	iStatus := 200
	sMsg := "success"
	var res interface{}

	filterUserMap := u.FilterUser()

	sUsername := filterUserMap["username"]      //用户名
	sPassword := filterUserMap["password"]      //用户密码
	sIp := filterUserMap["ip"]                  //接口调用方服务器
	sMerchantId := filterUserMap["merchant_id"] //接入商id

	if debug {
		fmt.Println("-------------传入的参数-------------")
		fmt.Println("sUsername-->", sUsername)
		fmt.Println("sPassword-->", sPassword)
		fmt.Println("sIp-->", sIp)
		fmt.Println("sMerchantId-->", sMerchantId)
	}

	if len(sUsername) < 1 || len(sPassword) < 1 || len(sIp) < 1 || len(sMerchantId) < 1 {
		iStatus = 7019
		sMsg = "参数为空或者参数错误"
		u.RenderJson(iStatus, sMsg, res)
	}

	var sNewPassword, sNewPrizeGroup string
	if sValue, ok := filterUserMap["new_password"]; ok { //新密码
		sNewPassword = sValue
	}
	if sValue, ok := filterUserMap["new_prize_group"]; ok { //新奖金组
		sNewPrizeGroup = sValue
	}

	///////////
	//判斷用戶//
	//////////
	reStatus, reMsg, reResult, oUser := u.CheckUser(sUsername, sPassword, sMerchantId)
	if reStatus != 200 {
		u.RenderJson(reStatus, reMsg, reResult)
	}

	if debug {
		fmt.Println("sNewPassword-->", sNewPassword)
		fmt.Println("sNewPrizeGroup-->", sNewPrizeGroup)
	}

	var bUpdatePrize, bUpdatePwd = false, false
	if len(sNewPassword) > 5 {
		bUpdatePwd = true
	}
	if len(sNewPrizeGroup) > 3 {
		iNewPrizeGroup, err := strconv.Atoi(sNewPrizeGroup)
		if err != nil && iNewPrizeGroup > 0 && iNewPrizeGroup < 2000 && sNewPrizeGroup != oUser.PrizeGroup {
			bUpdatePrize = true
		}
	}

	if debug {
		fmt.Println("---------------更新内容---------------")
		fmt.Println("bUpdatePrize-->", bUpdatePrize)
		fmt.Println("bUpdatePwd-->", bUpdatePwd)
	}
	if !bUpdatePrize && !bUpdatePwd {
		iStatus = 7020
		sMsg = "没有找到可更新内容"
		u.RenderJson(iStatus, sMsg, res)
	}

	//////////////////
	//開始編輯用戶資料//
	/////////////////
	o := orm.NewOrm()
	o.Begin()
	if bUpdatePwd {
		oUser.Password = models.EncodeString(sNewPassword)
	}

	if bUpdatePrize {
		oUser.PrizeGroup = sNewPrizeGroup

		//获取所有奖金组
		aConditions := map[string]string{
			"user_id": strconv.Itoa(oUser.Id),
		}
		aUserPrizeSets, _ := dao.GetAllUserPrizeSets(aConditions, nil, nil, nil, 0, 100000)
		for _, oUserPrizeSets := range aUserPrizeSets {
			oUserPrizeSets.PrizeGroup = sNewPrizeGroup
			err := dao.UpdateUserPrizeSetsById(oUserPrizeSets, o)
			if err != nil {
				o.Rollback()
				iStatus = 7021
				sMsg = "更新用户奖金组失败"
				res = err
				u.RenderJson(iStatus, sMsg, res)
			}
		}

	}

	err := dao.UpdateUsersById(oUser, o)
	if err != nil {
		o.Rollback()
		iStatus = 7022
		sMsg = "用户信息更新失败"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}
	err = o.Commit()
	if err != nil {
		iStatus = 7023
		sMsg = "模型错误"
		res = err
		u.RenderJson(iStatus, sMsg, res)
	}

	models.Users.FlushUserCache(sMerchantId, sUsername)
	u.RenderJson(iStatus, sMsg, res)
}

/*
 * 檢查用戶數據(公用)
 */
func (u *UserController) CheckUser(sUsername, sPassword, sMerchantId string) (int, string, interface{}, *dao.Users) {
	iStatus := 200
	sMsg := "success"
	var res interface{}
	var empty *dao.Users
	mConditions := map[string]string{
		"username":    sUsername,
		"merchant_id": sMerchantId,
	}
	aUsers, err := dao.GetAllUsers(mConditions, nil, nil, nil, 0, 1)
	if len(aUsers) != 1 {
		iStatus = 7021
		sMsg = "用戶信息錯誤"
		res = err
		return iStatus, sMsg, res, empty
	}

	filterUserMap := u.FilterUser()

	oUser := aUsers[0]
	sIdentity := filterUserMap["identity"]
	iMerchantId, _ := strconv.Atoi(sMerchantId)
	if iMerchantId <= 0 {
		iStatus = 7022
		sMsg = "接入商信息錯誤"
		res = err
		return iStatus, sMsg, res, empty
	}

	//验证帐号密码是否正确
	reStatus, reMsg, userRow := models.ChkLogAccount(sIdentity, sUsername, sPassword)
	if reStatus != 200 {
		iStatus = reStatus
		sMsg = reMsg
		res = userRow
		return iStatus, sMsg, res, empty
	}
	return iStatus, sMsg, res, oUser
}

// @Title UserPrizeSet
// @Description 修改用户彩种奖金组
// @Param	merchant_identity		query 	string	true		"接入方标识代码 如:JMG "
// @Param	params				query 	string	true		"参数username=xxx&prize_group=1950"
// username		 	string			"商户平台用户名"
// prize_group		 	string			"新奖金组"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /user_prize_set.do [post]
func (u *UserController) UserPrizeSet() {
	Status := 200
	Msg := "success"
	var result interface{}

	filterUserMap := u.FilterUser()

	sPrizeGroup, ok := filterUserMap["prize_group"]
	if !ok || sPrizeGroup == "" {
		Status = 100
		Msg = "Prize Group Error!"
		u.RenderJson(Status, Msg, result)
	}

	sUserName, ok := filterUserMap["username"]
	if !ok || sUserName == "" {
		Status = 101
		Msg = "Username  Error!"
		u.RenderJson(Status, Msg, result)
	}

	mUser := lotteryJobsModels.User.GetByName(filterUserMap["merchant_id"], sUserName)
	if len(mUser) == 0 {
		Status = 102
		Msg = "User Is Not found!"
		u.RenderJson(Status, Msg, result)
	}

	if mUser["prize_group"] == sPrizeGroup {
		Status = 103
		Msg = "Old Is Prize Group  is the same New Prize Group !" + sPrizeGroup
		u.RenderJson(Status, Msg, result)
	}

	mPrizeGroups := lotteryJobsModels.PrizeGroup.GetPrizeGroupByName(sPrizeGroup)
	if len(mPrizeGroups) == 0 {
		Status = 104
		Msg = "Prize Group  Error !"
		u.RenderJson(Status, Msg, result)
	}

	mLotteriesSeries := lotteryJobsModels.Lottery.GetAllLotteryIdsGroupBySeries()

	mPrizeNewGroups := map[string]map[string]string{}
	for _, mPrizeGroup := range mPrizeGroups {
		mPrizeNewGroups[mPrizeGroup["series_id"]] = mPrizeGroup
	}

	db := new(lotteryJobsModels.Table)

	TX := db.BeginTransaction()

	//更新用户彩种奖金组
	mUserPrizeSets := lotteryJobsModels.UserPrizeSet.GetUserPrizeSets(mUser["id"], "")
	for _, mUserPrizeSet := range mUserPrizeSets {
		mPrizeGroup := mPrizeNewGroups[mLotteriesSeries[mUserPrizeSet["lottery_id"]]]
		sClassicPrize := mPrizeGroup["classic_prize"]
		sPrizeGroupId := mPrizeGroup["id"]

		mUpdateData := map[string]string{
			"group_id":      sPrizeGroupId,
			"prize_group":   sPrizeGroup,
			"classic_prize": sClassicPrize,
		}
		sWhere := fmt.Sprintf("user_id='%s' and lottery_id='%s'", mUser["id"], mUserPrizeSet["lottery_id"])
		sUserPrizeSetSql := lotteryJobsModels.UserPrizeSet.GetUpdateSql(mUpdateData, sWhere, "")
		if debug {
			fmt.Println("=====sUserPrizeSetSql=====", sUserPrizeSetSql)
		}

		uPrizeQuery, uPrizeErr := TX.Exec(sUserPrizeSetSql)
		if uPrizeErr != nil {
			TX.Rollback()
			Status = 105
			Msg = "Update User Prize Set Error!"
			u.RenderJson(Status, Msg, result)
		}

		//返回影响的行数
		uPrizeInt, uPrizeRowErr := uPrizeQuery.RowsAffected()
		if uPrizeRowErr != nil || uPrizeInt == 0 {
			TX.Rollback()
			Status = 105
			Msg = "Update User Prize Set  Error!"
			u.RenderJson(Status, Msg, result)
		}
	}

	//更新用户奖金组
	sUpdateUserSql := lotteryJobsModels.User.GetUpdateSql(map[string]string{"prize_group": sPrizeGroup}, fmt.Sprintf("id='%s'", mUser["id"]), "")

	if debug {
		fmt.Println("=====sUpdateUserSql=====", sUpdateUserSql)
	}

	uUserQuery, uUserErr := TX.Exec(sUpdateUserSql)
	if uUserErr != nil {
		TX.Rollback()
		Status = 106
		Msg = "Update User Prize Group  Error!"
		u.RenderJson(Status, Msg, result)
	}

	//返回影响的行数
	uUserInt, uUserRowErr := uUserQuery.RowsAffected()
	if uUserRowErr != nil || uUserInt == 0 {
		TX.Rollback()
		Status = 106
		Msg = "Update User Prize Group  Error!"
		u.RenderJson(Status, Msg, result)
	}

	TX.Commit()

	u.RenderJson(Status, Msg, result)
}

// @Title test
// @Description test
// @Param	merchant_identity		query 	controllers.IdentityInputType		true		"接入方标识代码 如:JMG "
// @Param	params		query 	controllers.ParamsInputType		true		"param参数"
// @Success 200 {JsonOut}  success!
// @Failure 500 error
// @router /test.do [post]
func (c *UserController) Test() {

}
