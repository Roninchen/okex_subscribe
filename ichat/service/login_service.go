package service

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	e "okex/ichat/enum" /* 取个别名 */
	m "okex/ichat/model"
	t "okex/ichat/tools"
	"okex/ichat/util"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var  client  http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	client = http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not allow redirect
		},
	}
}
/* 从微信服务器获取登陆uuid */
func GetUUIDFromWX() (string, error) {
	paraMap := e.GetUUIDParaEnum()
	paraMap[e.TimeStamp] = fmt.Sprintf("%d", time.Now().Unix())

	resp, err := http.Get(e.UUID_URL + t.GetURLParams(paraMap))
	if err != nil {
		return "", errors.New("访问微信服务器API失败:" + err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("获取微信API反馈UUID数据失败:" + err.Error())
	}

	/* 正则解析uuid ,FindStringSubmatch首先返回整体匹配结果，然后返回每个子表达式的匹配结果*/
	reg := regexp.MustCompile(`^window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+)";$`)
	matches := reg.FindStringSubmatch(string(bodyBytes))
	if len(matches) != 3 {
		return "", errors.New("解析微信UUID API数据失败:" + string(bodyBytes))
	}
	status, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return "", errors.New("解析微信UUID API数据失败:" + err.Error())
	}

	if status != 200 {
		return "", errors.New(fmt.Sprintf("微信返回的状态错误，请排查网络故障，如仍有问题，可能是微信更改了API，请联系开发者。status:%d", status))
	}

	return matches[2], nil
}

/* 下载URL指向的JPG，保存到指定路径下的qrcode.jpg文件 */
func DownloadImagIntoDir(url string, dirPath string) error {
	//检查并创建临时目录
	if !isDirExist(dirPath) {
		os.Mkdir(dirPath, 0755)
		fmt.Println("dir %s created", dirPath)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dst, err := os.Create(dirPath + "/qrcode.jpg")
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, resp.Body)
	return err
}

/* 判断微信是否登陆 */
func CheckLogin(uuid string) (int64, string) {
	var timestamp int64 = time.Now().UnixNano() / 1000000
	paraMap := e.GetLoginParaEnum()
	paraMap[e.UUID] = uuid
	paraMap[e.TimeStamp] = fmt.Sprintf("%d", timestamp)
	paraMap[e.R] = fmt.Sprintf("%d", ^(int32)(timestamp))
	var paraStr = t.GetURLParams(paraMap)

	resp, err := http.Get(e.LOGIN_URL + paraStr)
	if err != nil {
		return 0, "访问微信服务器API失败:" + err.Error()
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "获取微信API反馈登陆数据失败:" + err.Error()
	}

	var resultStr = string(bodyBytes)

	reg := regexp.MustCompile(`^window.code=(\d+);`)
	matches := reg.FindStringSubmatch(resultStr)
	if len(matches) < 2 {
		return 0, "预期的返回结果格式不匹配"
	}

	status, err := strconv.ParseInt(matches[1], 10, 64)

	return status, resultStr
}

/**
 * 处理微信登陆成功时返回的其他登陆数据
 * 首先根据回调URI再次Get一次微信服务器，得到XML格式的登陆数据
 * 解析XML，向map中压入相关数据
 */
func ProcessLoginInfo(loginContent string) (m.LoginMap, error) {
	resultMap := m.LoginMap{}

	re := regexp.MustCompile(`window.redirect_uri="(\S+)";`)
	matches := re.FindStringSubmatch(loginContent)
	resultMap.Info = make(map[string]string)
	resultMap.Info["url"] = matches[1]

	req, _ := http.NewRequest("GET", resultMap.Info["url"], nil)
	req.Header.Add("User-Agent", e.USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return resultMap,err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resultMap,err
	}
	result := m.LoginCallbackXMLResult{}
	err = xml.Unmarshal(b, &result)
	if err != nil {
		return resultMap,err
	}

	resultMap.BaseRequest.DeviceID = util.GetRandomID(15)
	resultMap.BaseRequest.SKey = result.SKey
	resultMap.BaseRequest.Sid = result.WXSid
	resultMap.BaseRequest.Uin = result.WXUin
	resultMap.PassTicket = result.PassTicket
	resultMap.Info["pass_ticket"] = result.PassTicket
	resultMap.Info["logintime"] = util.GetTimestamp()
	resultMap.Info["url"] = resultMap.Info["url"][:strings.LastIndex(resultMap.Info["url"], "/")]

	indexUrl := req.URL.Hostname()
	if detailedUrl, ok := e.DetailedUrls[indexUrl]; ok {
		resultMap.Info["fileUrl"] = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin", detailedUrl[0])
		resultMap.Info["syncUrl"] = fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin", detailedUrl[1])
	} else {
		resultMap.Info["fileUrl"] = resultMap.Info["url"]
		resultMap.Info["syncUrl"] = resultMap.Info["url"]
	}

	/* 收集Cookie */
	resultMap.Cookies = resp.Cookies()
	return resultMap, nil
}

/* 初始化微信登陆数据 */
func WebInit(loginMap *m.LoginMap) error {
	data := map[string]interface{}{}
	data["BaseRequest"] = loginMap.BaseRequest

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", loginMap.Info["url"]+"/webwxinit", strings.NewReader(string(b)))
	req.URL.RawQuery = url.Values{
		"r":           {util.GetR()},
		"pass_ticket": {loginMap.Info["pass_ticket"]},
		"lang":        {"zh_CN"},
	}.Encode()
	req.Header.Add("ContentType", e.JSON_HEADER)
	req.Header.Add("User-Agent", e.USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	initInfo := m.InitInfo{}
	err = json.Unmarshal(b, &initInfo)
	if err != nil {
		return err
	}
	loginMap.User = initInfo.User
	loginMap.SelfNickName = initInfo.User.NickName
	loginMap.SelfUserName = initInfo.User.UserName
	loginMap.SyncKeys = initInfo.SyncKeys
	loginMap.SyncKeyStr = initInfo.SyncKeys.ToString()

	fmt.Println(loginMap)

	return nil
}

/**
 * 通知微信服务器状态变化，只要通知即可，无需处理返回数据
 * {"BaseRequest":{"Uin":154158775,"Sid":"/nxZxJ0LclxmOw8v","Skey":"@crypt_3aaab8d5_cdfa952ec95e594b100f44aba942a73c","DeviceID":"e390742104557152"},"Code":3,"FromUserName":"@fc96d593487db4fb92b9a633aec8293b","ToUserName":"@fc96d593487db4fb92b9a633aec8293b","ClientMsgId":1504571331980}
 */
func NotifyStatus(loginMap *m.LoginMap) error {
	urlMap := map[string]string{
		e.PassTicket: loginMap.PassTicket}

	statusNotifyJsonData := map[string]interface{}{}
	statusNotifyJsonData["BaseRequest"] = loginMap.BaseRequest
	statusNotifyJsonData["Code"] = 3
	statusNotifyJsonData["FromUserName"] = loginMap.SelfUserName
	statusNotifyJsonData["ToUserName"] = loginMap.SelfUserName
	statusNotifyJsonData["ClientMsgId"] = time.Now().UnixNano() / 1000000

	jsonBytes, err := json.Marshal(statusNotifyJsonData)
	if err != nil {
		return err
	}

	//fmt.Println("notify json:\n" + string(jsonBytes))

	_, err = http.Post(e.STATUS_NOTIFY_URL+t.GetURLParams(urlMap), e.JSON_HEADER, strings.NewReader(string(jsonBytes)))

	return err
}

func ShowMobileLogin(loginMap *m.LoginMap) {
	data := map[string]interface{}{}
	data["BaseRequest"] = loginMap.BaseRequest
	data["Code"] = 3
	data["FromUserName"] = loginMap.User.UserName
	data["ToUserName"] = loginMap.User.UserName
	data["ClientMsgId"] = util.GetTimestamp()

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, _ := http.NewRequest("POST", loginMap.Info["url"]+"/webwxstatusnotify", strings.NewReader(string(b)))
	req.Header.Add("ContentType", e.JSON_HEADER)
	req.Header.Add("User-Agent", e.USER_AGENT)

	resp, err := client.Do(req)
	bytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func isDirExist(path string) bool {
	p, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return p.IsDir()
	}
}

/**
 * 附录：
 * 微信登陆后回调URL的XML返回结果：
 * <error>
 *   <ret>0</ret>
 *   <message></message>
 *   <skey>@crypt_3aaab8d5_aa9febb1c57122a4569c1b1dc4772eac</skey>
 *   <wxsid>vjqCszEkQQw9jep1</wxsid>
 *   <wxuin>154158775</wxuin>
 *   <pass_ticket>wbFO7Vqg%2BpADuIcrQPDM1e0KjmNvgsH8jYAEoT0FtSY%3D</pass_ticket>
 *   <isgrayscale>1</isgrayscale>
 * </error>
 *
 */
