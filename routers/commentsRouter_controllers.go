package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "CalculatePrize",
			Router: `/calculate_prize.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "CancelPrize",
			Router: `/cancel_prize.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "CancelIssue",
			Router: `/cancel_issue.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "CancelProject",
			Router: `/cancel_project.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "CancelTrace",
			Router: `/cancel_trace.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AdminApiController"],
		beego.ControllerComments{
			Method: "MakeIssueCache",
			Router: `/make_issue_cache.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:AesController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:AesController"],
		beego.ControllerComments{
			Method: "AesEncryptString",
			Router: `/encrypt.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "GetProjectListServer",
			Router: `/bet_record_list.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "Bet",
			Router: `/bet.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "TrendViewData",
			Router: `/tides.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "GetProjectDetailt",
			Router: `/project_detailt.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "GetAccountChangeList",
			Router: `/account_change_list.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "DropProject",
			Router: `/cancel_project.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:GameController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:GameController"],
		beego.ControllerComments{
			Method: "CurrentUserInfo",
			Router: `/current_user_info.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "Token",
			Router: `/token.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "LotteryData",
			Router: `/lottery_data.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "Load",
			Router: `/load-data.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "LoadIssues",
			Router: `/load-numbers.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "NoticeList",
			Router: `/notice_list.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PlatPrizeData",
			Router: `/plat_prize_data.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetGameMenu",
			Router: `/get_game_menu.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetBetRecord",
			Router: `/get_bet_record.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PrintProjects",
			Router: `/print_projects.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "Test",
			Router: `/test.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetDayProjects",
			Router: `/day.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "GetTraceListServer",
			Router: `/trace_list.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "GetTraceDetailServer",
			Router: `/trace_detailt.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "CancelTraceReserveServer",
			Router: `/cancel_issue_trace.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "StopTraceServer",
			Router: `/cancel_trace.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "GetTraceProjectDetail",
			Router: `/trace_project_detail.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:TraceController"],
		beego.ControllerComments{
			Method: "SelectLotteries",
			Router: `/get_select_lottery.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "InitUser",
			Router: `/init.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Balance",
			Router: `/balance.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Transfer",
			Router: `/transfer.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Transferinfo",
			Router: `/transferinfo.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Edit",
			Router: `/edit.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "UserPrizeSet",
			Router: `/user_prize_set.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["NewLottApi/controllers:UserController"] = append(beego.GlobalControllerRouter["NewLottApi/controllers:UserController"],
		beego.ControllerComments{
			Method: "Test",
			Router: `/test.do`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
