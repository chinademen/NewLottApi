package models

import (
	"NewLottApi/dao"
	"common"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type tTransaction struct {
	Table
}

var TransactionList = &tTransaction{Table: Table{TableName: "transaction"}}

const (
	TransferIn  = 1 //额度转移,转入
	TransferOut = 2 //额度转移,转出

	ERRNO_CREATE_ERROR_DATA    int = -101
	ERRNO_CREATE_ERROR_SAVE    int = -102
	ERRNO_CREATE_SUCCESSFUL    int = -100
	ERRNO_CREATE_ERROR_BALANCE int = -103
)

var NoIssueLotteries = []string{"31", "85", "89", "87", "86", "91", "83", "78", "54", "94", "95", "96", "97", "98", "100", "101", "99", "58"}

/**
 * 增加新的账变
 * @param *dao.Users      oUser
 * @param *dao.Accounts   oAccount
 * @param *dao.Projects   oProject
 * @param int             iType
 * @param float64         fAmount
 */
func (m *tTransaction) AddProjectTransaction(o orm.Ormer, oUser *dao.Users, oAccount *dao.Accounts, oProject *dao.Projects, oSeriesWay *dao.SeriesWays, iType int, fAmount float64, mExt map[string]string) (int, error) {
	if fAmount <= 0 {
		return ERRNO_CREATE_ERROR_DATA, nil
	}

	//獲取賬變對象和賬戶變化對象
	oTransaction, oNewBalance, bSuccess := m.compileProjectData(oUser, oAccount, oProject, oSeriesWay, iType, oProject.Amount, mExt)
	if !bSuccess {
		return ERRNO_CREATE_ERROR_SAVE, nil
	}

	iTransactionId, err := dao.AddTransactions(o, oTransaction)

	if iTransactionId < 1 {
		return ERRNO_CREATE_ERROR_SAVE, err
	}

	err = dao.UpdateAccountsById(o, oNewBalance)
	if err != nil {
		return ERRNO_CREATE_ERROR_BALANCE, err
	}

	return ERRNO_CREATE_SUCCESSFUL, nil
}

/**
 * 增加新的账变
 * @param *dao.Users      oUser
 * @param *dao.Accounts   oAccount
 * @param *dao.Traces     oTrace
 * @param int             iType
 * @param float64         fAmount
 */
func (m *tTransaction) AddTraceTransaction(o orm.Ormer, oUser *dao.Users, oAccount *dao.Accounts, oTrace *dao.Traces, oSeriesWay *dao.SeriesWays, iType int, fAmount float64, mExt map[string]string) (int, *dao.Accounts) {

	if fAmount <= 0 {
		return ERRNO_CREATE_ERROR_DATA, oAccount
	}

	//獲取賬變對象和賬戶變化對象
	oTransaction, oNewBalance, bSuccess := m.compileTraceData(oUser, oAccount, oTrace, oSeriesWay, iType, oTrace.Amount, mExt)
	if !bSuccess {
		return ERRNO_CREATE_ERROR_SAVE, oNewBalance
	}
	iTransactionId, _ := dao.AddTransactions(o, oTransaction)
	if iTransactionId < 1 {
		return ERRNO_CREATE_ERROR_SAVE, oNewBalance
	}

	err := dao.UpdateAccountsById(o, oNewBalance)
	if err != nil {
		return ERRNO_CREATE_ERROR_BALANCE, oNewBalance
	}
	return ERRNO_CREATE_SUCCESSFUL, oNewBalance
}

/*
 * 返回用戶賬戶對象和賬變對象(針對追號)
 */
func (m *tTransaction) compileTraceData(oUser *dao.Users, oAccount *dao.Accounts, oTrace *dao.Traces, oSeriesWay *dao.SeriesWays, iType int, fAmount float64, mExt map[string]string) (*dao.Transactions, *dao.Accounts, bool) {
	oTransactionType, _ := dao.GetTransactionTypesById(iType)
	oTransaction := new(dao.Transactions)
	if oTransactionType.Id == 0 {
		return oTransaction, oAccount, false
	}

	iMechantId, _ := strconv.Atoi(mExt["merchant_id"])

	oTransaction.UserId = uint64(oUser.Id)
	oTransaction.IsTester = oTrace.IsTester
	oTransaction.Amount = fAmount
	oTransaction.TypeId = uint32(iType)
	oTransaction.IsIncome = uint8(oTransactionType.Credit)
	oTransaction.PreviousFrozen = oAccount.Frozen
	oTransaction.PreviousBalance = oAccount.Balance
	oTransaction.PreviousAvailable = oAccount.Available
	oTransaction.Frozen = oAccount.Frozen
	oTransaction.Balance = oAccount.Balance
	oTransaction.Available = oAccount.Available
	oTransaction.AccountId = uint64(oAccount.Id)
	oTransaction.Username = oUser.Username
	oTransaction.Description = oTransactionType.Description
	oTransaction.TerminalId = oTrace.TerminalId
	oTransaction.ProxyIp = oTrace.ProxyIp
	oTransaction.Issue = oTrace.StartIssue
	oTransaction.LotteryId = oTrace.LotteryId
	oTransaction.WayId = uint(oTrace.WayId)
	oTransaction.Coefficient = oTrace.Coefficient
	oTransaction.TraceId = uint64(oTrace.Id)
	oTransaction.Ip = mExt["clientIP"]
	oTransaction.ProxyIp = mExt["proxyIP"]
	oTransaction.MerchantId = uint(iMechantId)
	oTransaction.WayId = uint(oSeriesWay.Id)

	oTransaction.Note = ""
	//	oTransaction.AdminUserId =
	//	oTransaction.AdminUserId =
	if oTransactionType.ProjectLinked > 0 {
		//	oTransaction.ProjectId =
		//	oTransaction.ProjectNo
	}

	oTransaction = m.SaveKey(oTransaction)
	oTransaction = m.MakeSeriesNumber(oTransaction)

	fBalance := oAccount.Balance + float64(oTransactionType.Balance)*fAmount
	oTransaction.Balance = fBalance
	oAccount.Balance = fBalance

	fAvailable := oAccount.Available + float64(oTransactionType.Available)*fAmount
	oTransaction.Available = fAvailable
	oAccount.Available = fAvailable

	fFrozen := oAccount.Frozen + float64(oTransactionType.Frozen)*fAmount
	oTransaction.Frozen = fFrozen
	oAccount.Frozen = fFrozen

	//创建时间
	timeDateNow := common.GetNowDatetime(common.DATE_FORMAT_YMDHIS)
	oTransaction.CreatedAt = timeDateNow
	oTransaction.UpdatedAt = timeDateNow

	return oTransaction, oAccount, true
}

/*
 * 返回用戶賬戶對象和賬變對象(針對注單)
 */
func (m *tTransaction) compileProjectData(oUser *dao.Users, oAccount *dao.Accounts, oProject *dao.Projects, oSeriesWay *dao.SeriesWays, iType int, fAmount float64, mExt map[string]string) (*dao.Transactions, *dao.Accounts, bool) {
	oTransactionType, _ := dao.GetTransactionTypesById(iType)
	oTransaction := new(dao.Transactions)
	if oTransactionType.Id == 0 {
		return oTransaction, oAccount, false
	}

	iMechantId, _ := strconv.Atoi(mExt["merchant_id"])

	oTransaction.UserId = uint64(oUser.Id)
	oTransaction.IsTester = oProject.IsTester
	oTransaction.Amount = fAmount
	oTransaction.TypeId = uint32(iType)
	oTransaction.IsIncome = uint8(oTransactionType.Credit)
	oTransaction.PreviousFrozen = oAccount.Frozen
	oTransaction.PreviousBalance = oAccount.Balance
	oTransaction.PreviousAvailable = oAccount.Available
	oTransaction.Frozen = oAccount.Frozen
	oTransaction.Balance = oAccount.Balance
	oTransaction.Available = oAccount.Available
	oTransaction.AccountId = uint64(oAccount.Id)
	oTransaction.Username = oUser.Username
	oTransaction.Description = oTransactionType.Description
	oTransaction.TerminalId = uint8(oProject.TerminalId)
	oTransaction.ProxyIp = oProject.ProxyIp
	oTransaction.Issue = oProject.Issue
	oTransaction.LotteryId = oProject.LotteryId
	oTransaction.WayId = uint(oProject.WayId)
	oTransaction.Coefficient = oProject.Coefficient
	oTransaction.TraceId = uint64(oProject.Id)
	oTransaction.Ip = mExt["clientIP"]
	oTransaction.ProxyIp = mExt["proxyIP"]
	oTransaction.MerchantId = uint(iMechantId)
	oTransaction.WayId = uint(oSeriesWay.Id)
	oTransaction.ProjectId = uint64(oProject.Id)

	oTransaction.Note = ""
	if oTransactionType.ProjectLinked > 0 {
		//	oTransaction.ProjectId =
		//	oTransaction.ProjectNo
	}

	oTransaction = m.SaveKey(oTransaction)
	oTransaction = m.MakeSeriesNumber(oTransaction)

	//更新賬戶餘額
	fBalance := oAccount.Balance + float64(oTransactionType.Balance)*fAmount
	oTransaction.Balance = fBalance
	oAccount.Balance = fBalance

	//更新賬戶可用餘額
	fAvailable := oAccount.Available + float64(oTransactionType.Available)*fAmount
	oTransaction.Available = fAvailable
	oAccount.Available = fAvailable

	//更新賬戶凍結餘額
	fFrozen := oAccount.Frozen + float64(oTransactionType.Frozen)*fAmount
	oTransaction.Frozen = fFrozen
	oAccount.Frozen = fFrozen

	//创建时间
	timeDateNow := common.GetNowDatetime(common.DATE_FORMAT_YMDHIS)
	oTransaction.CreatedAt = timeDateNow
	oTransaction.UpdatedAt = timeDateNow

	return oTransaction, oAccount, true
}

/*
 * 生成賬變唯一值
 */
func (m *tTransaction) SaveKey(oTransactions *dao.Transactions) *dao.Transactions {

	aFields := []string{}
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.UserId))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.TypeId))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.AccountId))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.TraceId))
	aFields = append(aFields, fmt.Sprintf("%f", oTransactions.Amount))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.LotteryId))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.WayId))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.ProjectId))
	aFields = append(aFields, fmt.Sprintf("%f", oTransactions.Amount))
	aFields = append(aFields, fmt.Sprintf("%d", oTransactions.AdminUserId))
	aFields = append(aFields, oTransactions.Ip)
	aFields = append(aFields, oTransactions.ProxyIp)
	aFields = append(aFields, oTransactions.Description)
	aFields = append(aFields, oTransactions.Issue)

	sFields := strings.Join(aFields, "|")
	oTransactions.Safekey = common.GetMd5(sFields)
	return oTransactions
}

/*
 * 生成賬變標識符serialNumber
 */
func (m *tTransaction) MakeSeriesNumber(oTransactions *dao.Transactions) *dao.Transactions {
	oTransactions.SerialNumber = common.Uniqid(fmt.Sprintf("%d", oTransactions.UserId), true)
	return oTransactions
}
