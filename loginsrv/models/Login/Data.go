package Login

//----------------------------------------------- REQUEST -------------------------------------------------
//客户端请求：验证码
type ReqCaptcha struct {
	Flag string `json:"flag"`
	//Platform int    `json:"platform"`
}

//客户端请求：注册
type ReqRegist struct {
	AccountName          string `json:"userAccount"`
	Password             string `json:"password"`
	Captcha              string `json:"captcha"`
	RegistrationPlatform int    `json:"registrationPlatform"` //0苹果，1安卓，2，Wap，3 PC网页端
	Flag                 string `json:"flag"`                 //标识码用于判定请求验证码的机器和注册的机器是否一致
	AccountType          int    `json:"accountType"`          //账户类型 0 为试玩用户,1为普通注册用户
	InviteCode           string `json:"inviteCode"`           //邀请码(邀请码可以为空)
	Ip                   string `json:"ip"`                   //Ip 注册IP,web 和 wap 是由客户端发来,app由服务器从消息来源获取
}

//客户端请求：登陆
type ReqLogin struct {
	AccountName string `json:"accountName"`
	Password    string `json:"password"`
	Flag        string `json:"flag"`
	Ip          string `json:"ip"`
}

//预登录请求
type ReqPrepareLogin struct {
	AccountName        string `json:"accountName"`
	Token              string `json:"token"`
	Flag               string `json:"flag"`
	AccountId          int    `json:"accountId"`
	AccountType        int    `json:"accountType"`
	InviteCode         string `json:"inviteCode"`
	RegistTimeStamp    int64  `json:"registTime"`
	RegistIp           string `json:"registIp"`
	LastLoginTimeStamp int64  `json:"lastLoginTimeStamp"`
	LastLoginIp        string `json:"lastLoginIp"`
}

//----------------------------------------------- RESPONSE -------------------------------------------------
//RespCommon 服务器应答：公用回复
type RespCommon struct {
	Token string `json:"token"`
	State int    `json:"state"`
}

type RespLogin struct {
	State       int    `json:"state"`
	Token       string `json:"token"`
	GameIp      string `json:"gameIp"`
	AccountType int    `json:"account"`
}

//----------------------------------------------- MEMBER -------------------------------------------------

//以map类型保存标识符与验证码
var VerifyInfo = make(map[string]string)

//GameServerIp Game服务器启动时会通知Login有那些Game服务器启动 //key 为服务器ID value 为服务器ip和端口字符串
var GameServerIp = make(map[int]string)

//验证码图片的长宽定义
const (
	Dx = 60
	Dy = 30
)

//----------------------------------------------- Game 服务器发来的消息结构 -------------------------------------------------
//服务器注册信息
type ServerRegistInfo struct {
	Id     int    `json:"id"`     //服务器id
	Port   int    `json:"port"`   //服务器端口(注意:这个服务器端口是指,注册服务器用于接收消息的端口)
	Cipher []byte `json:"cipher"` //rsa密文
}
