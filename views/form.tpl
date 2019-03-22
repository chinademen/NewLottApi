<!doctype html>
<html>
  	<head>
    	<title>游戏平台中间件</title>
    	<meta http-equiv="content-type" content="text/html; charset=utf-8">
		
	</head>
<style>
table,tr,td{
	border:1px #000000 solid;
}
</style>
<body>
<p><b>游戏平台中间件，api调用说明</b></p>
<p>
ajax返回值说明:
</p>
<p>
{
  "Status": 200,
  "Msg": "请求完成",
  "Data": {
    "account": "rrrrr"
  }
}
<p>



<p>验证接入商是否可用<br />
/access/info.do</p>
<table>
<form method="post" action="/access/info.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text" value="test"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><textarea name="params">accesscode=test</textarea></td>
<td>accesscode=test</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>加密校验，验证您的加密结果是否更接口的一样<br />
/unhook/encode.do</p>
<table>
<form method="post" action="/unhook/encode.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>aeskey</td>
<td><input name="aeskey" type="text"></td>
<td>AES加密的密钥</td>
</tr>
<tr>
<td>params</td>
<td><textarea name="params"></textarea></td>
<td>需要加密的参数，格式和post表头提交的一样，多个参数用&分割;例如：account=test01&password=a123456789&plat=PT</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>
	
<p>初始化会员/创建一个会员<br />
/game/init.do</p>
<table>
<form method="post" action="/game/init.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>传递的参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度5-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,6-20小写位英文和数字</p>
<p>other：扩展字段,比如代理</p>
<p>level：用户等级,网站的会员等级</p>
<p>test：是否测试账户,1=测试，区别在于测试可以产生的数据可以过滤，不计算报表</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>用户的转入转出<br />
/cash/transfer.do</p>
<table>
<form method="post" action="/cash/transfer.do">
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>加密参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,8-15小写位英文和数字</p>
<p>amount：转账的额度,整数</p>
<p>orderid：网站的订单号</p>
<p>ispush：on=推送，默认不推送</p>
<p>url：推送url</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>登录游戏平台<br />
/game/gamelogin.do</p>
<table>
<form method="post" action="/game/gamelogin.do">
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>加密参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,8-15小写位英文和数字</p>
<p>gamekind：扩展1</p>
<p>gamename：扩展2</p>
<p>gameid：扩展3</p>
<p>ip：ip</p>
<p>drivetype：h5或者pc， 默认pc</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>修改密码<br />
/game/changepwd.do</p>
<table>
<form method="post" action="/game/changepwd.do">
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>加密参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,8-15小写位英文和数字</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>查询用户的余额<br />
/cash/balance.do</p>
<table>
<form method="post" action="/cash/balance.do">
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>加密参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,8-15小写位英文和数字</p>
<p>plat：游戏平台编号，必须大写</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>查询投注记录<br />
/game/getbet.do</p>
<table>
<form method="post" action="/game/getbet.do">
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>加密参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>plat：游戏平台编号，必须大写</p>
<p>startdate：开始日期,格式：2015-09-28 10:10:10</p>
<p>enddate：结束日期,不写则默认到当前时间,格式同上</p>
<p>currpage：当前页,数字</p>
<p>pagenum：每页显示数目,默认是20条(最大)</p>
<p>注：如果数据不合法，则返回错误提示。最大查询日期区间是7天，如果超过7天，则查询的区间是从开始日期往后的七天数据</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>


<p>查询转账单据信息<br />
/cash/transferinfo.do</p>
<table>
<form method="post" action="/cash/transferinfo.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>accesscode</td>
<td><input name="accesscode" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><input name="params" type="text"></td>
<td>传递的参数</td>
</tr>
<tr>
<td>参数</td>
<td rowspan="2">
<p>account：游戏用户名(长度6-12位之间，只能字母和数字，首位字母)</p>
<p>password：游戏用户密码,8-15小写位英文和数字</p>
<p>orderid：转账的单据号</p>
</td>
</tr>
<tr>
<td>加密格式</td>
<td>aes(字段=值&字段=值,密钥)</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>



	</body>
</html>