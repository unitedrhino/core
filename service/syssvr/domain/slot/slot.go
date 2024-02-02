package slot

type Info struct {
	Code     string            `json:"code"`     // 鉴权的编码
	SlotCode string            `json:"slotCode"` //slot的编码
	Method   string            `json:"method"`   // 请求方式 GET  POST
	Uri      string            `json:"uri"`      // 参考: /api/v1/system/user/self/captcha?fwefwf=gwgweg&wefaef=gwegwe
	Proto    string            `json:"proto"`    //http https
	Hosts    []string          `json:"hosts"`    //访问的地址 host or host:port
	Body     string            `json:"body"`     // body 参数模板
	Handler  map[string]string `json:"handler"`  //http头
	AuthType string            `json:"authType"` //鉴权类型 core
}
