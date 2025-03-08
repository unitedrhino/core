package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/trace"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	dp              *dispatcher //ws调度器
	once            sync.Once
	store           kv.Store
	nodeID          int64
	checkSubscribe  func(ctx context.Context, in *SubscribeInfo) error
	checkSubscribe2 func(ctx context.Context, in *SubscribeInfo) ([]map[string]any, error)
	connectID       atomic.Int64
)

const (
	errorCount    = 5                     //错误次数
	interval      = 10 * time.Second      //心跳间隔
	keepAliveType = websocket.PingMessage //心跳类型
)

type connection struct {
	r                  *http.Request
	uc                 *ctxs.UserCtx
	server             *Server
	ws                 *websocket.Conn //ws连接实例
	userID             int64           //ws连接实例唯一标识
	connectID          int64
	userSubscribeMutex sync.Mutex
	userSubscribe      map[string]map[string]struct{}
	closed             bool        //ws连接已关闭
	send               chan []byte //发送信息管道
	pingErr            atomic.Int64
}

// ws调度器
type dispatcher struct {
	s2cGzip            bool                            //发送的信息是否gzip压缩
	connPool           map[int64]map[int64]*connection //ws连接池 一个用户会有多个端接入,也就会有多个ws连接 第一个key是userID 第二个key是ConnectID
	mu                 sync.RWMutex                    // 互斥锁
	userSubscribeMutex sync.RWMutex
	userSubscribe      map[string]map[int64]*connection //第一个key是订阅参数的md5,第二个key是连接的id
}
type WsPublish struct {
	Code   string
	Data   any
	Params []string
}
type WsPublishes []WsPublish

// 创建ws调度器
func StartWsDp(s2cGzip bool, NodeID int64, event *eventBus.FastEvent, c cache.ClusterConf) {
	nodeID = NodeID
	once.Do(func() {
		store = kv.NewStore(c)
		dp = newDp(s2cGzip)
		event.Subscribe(fmt.Sprintf(eventBus.CoreApiUserPublish, ">"), func(ctx context.Context, t time.Time, body []byte) error {
			var pbs = WsPublishes{}
			logx.WithContext(ctx).Debugf("websocket StartWsDp.sendMessage nodeID:%v body:%v", nodeID, string(body))
			err := json.Unmarshal(body, &pbs)
			if err != nil {
				return err
			}
			for _, pb := range pbs {
				err := func() error {
					dp.userSubscribeMutex.RLock()
					defer dp.userSubscribeMutex.RUnlock()
					var sub map[int64]*connection
					var connectIDSet = map[int64]struct{}{}
					for _, param := range pb.Params {
						key := pb.Code + ":" + param
						sub = dp.userSubscribe[key]
						if sub == nil {
							continue
						}
						for _, c := range sub { //所有订阅者都需要发
							if _, ok := connectIDSet[c.connectID]; ok { //去重,一个连接只发一次
								continue
							}
							connectIDSet[c.connectID] = struct{}{}
							c.sendMessage(WsResp{
								WsBody: WsBody{
									Type: Pub,
									Path: pb.Code,
									Body: pb.Data,
								},
							})
						}
					}
					return nil
				}()
				if err != nil {
					logx.WithContext(ctx).Error(err)
					continue
				}
			}
			return nil
		})
		utils.Go(context.Background(), func() {
			t := time.NewTicker(time.Second * 10)
			for range t.C {
				func() {
					dp.userSubscribeMutex.RLock()
					defer dp.userSubscribeMutex.RUnlock()
					logx.Infof("websocket userSubscribeReadAll:%v", utils.Fmt(dp.userSubscribe))
				}()
			}
		})
	})
}

func RegisterSubscribeCheck(f func(ctx context.Context, in *SubscribeInfo) error) {
	checkSubscribe = f
}

func RegisterSubscribeCheck2(f func(ctx context.Context, in *SubscribeInfo) ([]map[string]any, error)) {
	checkSubscribe2 = f
}

// 创建ws调度器
func newDp(s2cGzip bool) *dispatcher {
	d := &dispatcher{
		s2cGzip:       s2cGzip,
		connPool:      make(map[int64]map[int64]*connection),
		userSubscribe: map[string]map[int64]*connection{},
	}
	return d
}

// 读ping心跳
func (c *connection) pingRead(message []byte) error {
	logx.Infof("websocket pingRead message:%s userID:%v", string(message), c.userID)
	c.pingErr.Store(0)

	return nil
}

// 读pong心跳
func (c *connection) pongRead(message []byte) error {
	logx.Infof("websocket pongRead message:%s userID:%v connectID:%v", string(message), c.userID, c.connectID)
	c.pingErr.Store(0)
	return nil
}

// 发送ping心跳
func (c *connection) pingSend() error {
	e := c.pingErr.Load()
	if e >= errorCount {
		logx.Infof("websocket connection timeout userID:%v connectID:%v pingErr:%v ", c.userID, c.connectID, e)
		//连续5次没有收到ping心跳 关闭连接
		return errors.System.AddMsg("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PingMessage, nowTime); err != nil {
		logx.Infof("websocket PingMessage userID:%v connectID:%v err:%v ", c.userID, c.connectID, err)
	}
	c.pingErr.Add(1)
	return nil
}

// 发送pong心跳
func (c *connection) pongSend() error {
	e := c.pingErr.Load()
	if e >= errorCount {
		logx.Infof("websocket connection timeout userID:%v connectID:%v pingErr:%v ", c.userID, c.connectID, e)
		//连续5次没有收到pong心跳 关闭连接
		return errors.System.AddMsg("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PongMessage, nowTime); err != nil {
		logx.Infof("websocket PongMessage userID:%v connectID:%v err:%v ", c.userID, c.connectID, err)
	}
	c.pingErr.Add(1)
	return nil
}

// 发送订阅信息
func SendSub(ctx context.Context, body WsResp) {
	clientToken := trace.TraceIDFromContext(ctx)
	body.Handler = map[string]string{"Traceparent": clientToken}
}

func AddConnPool(userID int64, conn *connection) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	if dp.connPool[userID] == nil {
		dp.connPool[userID] = map[int64]*connection{}
	}
	dp.connPool[userID] = map[int64]*connection{conn.connectID: conn}
}

// 创建ws连接
func NewConn(ctx context.Context, userID int64, server *Server, r *http.Request, wsConn *websocket.Conn) *connection {
	conn := &connection{
		server:        server,
		ws:            wsConn,
		uc:            ctxs.GetUserCtx(ctx),
		r:             r,
		userID:        userID,
		userSubscribe: map[string]map[string]struct{}{},
		connectID:     connectID.Add(1),
		send:          make(chan []byte, 10000),
	}
	AddConnPool(userID, conn)
	logx.Infof("websocket 创建连接成功 RemoteAddr::%s userID:%v connectID:%v uc:%v",
		wsConn.RemoteAddr().String(), userID, conn.connectID, utils.Fmt(conn.uc))
	resp := WsResp{WsBody: WsBody{}}
	clientToken := trace.TraceIDFromContext(ctx)
	resp.Handler = map[string]string{
		"Traceparent": clientToken,
		"connectID":   cast.ToString(conn.connectID),
	}
	conn.sendMessage(resp)
	return conn
}

// 开启读取进程
func (c *connection) StartRead() {
	defer func() {
		c.Close("read message error")
	}()
	c.ws.SetPongHandler(func(message string) error {
		c.pongRead([]byte(message))
		return nil
	})
	c.ws.SetPingHandler(func(message string) error {
		c.pingRead([]byte(message))
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			logx.Errorf("%s.websocket ReadMessage message:%s userID:%v connectID:%v err:%v",
				utils.FuncName(), string(message), c.userID, c.connectID, err)
			c.Close("read err")
			break
		}

		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			c.errorSend(WsReq{}, errors.Type.AddDetail("error reading message"))
			continue
		}
		c.pingErr.Store(0)
		c.handleRequest(message)
	}
}
func (c *connection) errorSend(req WsReq, data error) {
	e := errors.Fmt(data)
	resp := WsResp{
		WsBody: WsBody{Type: req.Type, Path: req.Path},
		Code:   e.GetCode(),
		Msg:    e.GetI18nMsg(""),
	}
	c.sendMessage(resp)
}

func (c *connection) handleRequest(message []byte) {
	var body WsReq
	err := json.Unmarshal(message, &body)
	if err != nil {
		c.errorSend(body, errors.Parameter)
		return
	}
	ctx, cancel := context.WithTimeout(ctxs.SetUserCtx(context.Background(), c.uc), 50*time.Second)
	defer cancel()
	utils.Recover(ctx)
	//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, string(body.Type), "")
	defer span.End()
	logx.WithContext(ctx).Infof("%s.websocket message:%s userID:%v connectID:%v", utils.FuncName(), string(message), c.userID, c.connectID)
	if err := isDataComplete(body.Type, body); err != nil {
		c.errorSend(body, err)
		return
	}
	if len(body.Handler) > 0 {
		for k, v := range body.Handler {
			c.r.Header.Set(k, v)
		}
	}
	switch body.Type {
	//case Control:
	//	downControl(c, body)
	case Sub:
		subscribeHandle(ctx, c, body)
	case UnSub:
		unSubscribeHandle(ctx, c, body)
	case UpPing:
		var resp = WsResp{WsBody: WsBody{}}
		resp.WsBody.Type = DownPong
		c.sendMessage(resp)
	default:
	}
}

func isDataComplete(wsType WsType, body WsReq) error {
	if wsType == "" {
		return errors.Parameter.AddDetail("type is  null")
	}
	switch wsType {
	case Control:
		if body.Path == "" || body.Method == "" || body.Body == "" {
			return errors.Parameter.AddDetail("path|method|body is  null")
		}
	case Sub, UnSub:
		if _, ok := body.Body.(map[string]interface{}); !ok {
			return errors.Parameter.AddDetail("body is  null")
		}
	case Pub:
		return errors.NotRealize
	default:
	}
	return nil
}

func downControl(c *connection, body WsReq) {
	reqBody, err := getRequestBody(body.Body)
	if err != nil {
		// 处理编码错误
	}
	bodyBytes, err := json.Marshal(body.Body)
	length := len(bodyBytes)
	header := c.r.Header
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", strconv.Itoa(length))
	r := &http.Request{
		Method: body.Method,
		Host:   c.r.Host,
		URL: &url.URL{
			Path: body.Path,
		},
		Header:        header,
		Body:          reqBody,
		ContentLength: int64(length),
	}
	w := response{req: &body, resp: WsResp{WsBody: WsBody{Handler: map[string]string{}, Type: ControlRet}}}
	c.server.ServeHTTP(&w, r)
	if token := w.Header().Get(ctxs.UserSetTokenKey); token != "" { //登录态保持更新
		c.r.Header.Set(ctxs.UserSetTokenKey, token)
	}
	c.sendMessage(w.resp)
}

// 将请求体转换为io.ReadCloser类型
func getRequestBody(body interface{}) (io.ReadCloser, error) {
	var reqBody io.ReadCloser
	if body != nil {
		switch body.(type) {
		case string:
			reqBody = io.NopCloser(bytes.NewBufferString(body.(string)))
		case []byte:
			reqBody = io.NopCloser(bytes.NewBuffer(body.([]byte)))
		case map[string]interface{}:
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = io.NopCloser(bytes.NewReader(bodyBytes))
		default:
			// 处理其他类型
		}
	}
	return reqBody, nil
}

// 开启发送进程
func (c *connection) StartWrite() {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		//发送心跳
		case <-ticker.C:
			if c.closed {
				return
			}
			var err error
			switch keepAliveType {
			case websocket.PingMessage:
				err = c.pingSend()
			case websocket.PongMessage:
				err = c.pongSend()
			}
			if err != nil {
				logx.Errorf("websocket pingPong keepAliveType:%v  userID:%v,connectID:%v writeMessage:%v err:%v",
					keepAliveType, c.userID, c.connectID, err)
				c.Close("connection timeout")
				return
			}
		//发送信息
		case message := <-c.send:
			if c.closed {
				return
			}
			logx.Infof("websocket userID:%v,connectID:%v writeMessage:%v", c.userID, c.connectID, string(message))
			if err := c.writeMessage(websocket.TextMessage, message); err != nil {
				logx.Errorf("websocket StartWrite  userID:%v,connectID:%v writeMessage:%v err:%v",
					c.userID, c.connectID, string(message), err)
				c.Close("send message error")
				return
			}
		}
	}
}

// 关闭ws连接
func (c *connection) Close(msg string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	_, ok := dp.connPool[c.userID]
	if ok || !c.closed {
		c.closed = true
		func() {
			defer func() {
				recover()
			}()
			close(c.send)

		}()
		delete(dp.connPool[c.userID], c.connectID)
		if len(dp.connPool[c.userID]) == 0 {
			delete(dp.connPool, c.userID)
		}
		c.userSubscribeMutex.Lock()
		defer c.userSubscribeMutex.Unlock()
		for key := range c.userSubscribe {
			func() {
				dp.userSubscribeMutex.Lock()
				defer dp.userSubscribeMutex.Unlock()
				sub, ok := dp.userSubscribe[key]
				if !ok {
					return
				}
				delete(sub, c.connectID)
			}()
		}
		err := c.ws.Close()
		logx.Infof("websocket 关闭连接 msg:%v  userID:%v connectID:%v err:%v", msg, c.userID, c.connectID, err)
	}
}

// 发送信息
func (c *connection) sendMessage(body WsResp) {
	if body.Code == 0 {
		body.Code = errors.OK.Code
		body.Msg = errors.OK.GetMsg()
	}
	message, _ := json.Marshal(body)
	if !c.closed {
		c.send <- message
	}
}

// 写消息
func (c *connection) writeMessage(messageType int, message []byte) error {
	if message == nil {
		logx.Infof("websocket error message: is  null ")
	}
	switch messageType {
	case websocket.PingMessage, websocket.PongMessage:
		err := c.ws.WriteControl(messageType, message, time.Time{})
		if err != nil {
			logx.Infof("%s.websocket  error message::%s userID:%v connectID:%v", utils.FuncName(), string(message), c.userID, c.connectID)
		}
	case websocket.TextMessage:
		err := c.ws.WriteMessage(messageType, message)
		if err != nil {
			logx.Infof("%s.websocket error message::%s userID:%v connectID:%v", utils.FuncName(), string(message), c.userID, c.connectID)
		}
	}
	logx.Debugf("%s.websocket message:%s userID:%v connectID:%v", utils.FuncName(), string(message), c.userID, c.connectID)
	return nil
}
