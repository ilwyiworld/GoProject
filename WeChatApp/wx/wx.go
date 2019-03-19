package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"github.com/clbanning/mxj"
)

type weixinQuery struct {	//微信get请求参数
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`	//post请求不存在
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {

	weixinClient := new(WeixinClient)

	weixinClient.Token = token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()

	//验证签名
	if weixinClient.Query.Signature != weixinClient.signature() {
		return nil, errors.New("Invalid Signature.")
	}

	return weixinClient, nil
}

func (this *WeixinClient) initWeixinQuery() {	//获取微信发来的get请求的参数
	var q weixinQuery
	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")
	this.Query = q
}

func (this *WeixinClient) signature() string {
	// 1）将token、timestamp、nonce三个参数进行字典序排序
	// 2）将三个参数字符串拼接成一个字符串进行sha1加密
	// 3）开发者获得加密后的字符串可与signature对比，标识该请求来源于微信
	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *WeixinClient) initMessage() error {

	//获取post请求数据
	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	//解析xml成为map[string]interface{}
	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

func (this *WeixinClient) text() {

	//获取文本数据
	inMsg, ok := this.Message["Content"].(string)

	if !ok {
		return
	}

	var reply TextMessage

	reply.InitBaseData(this, "text")
	reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))

	//数据转成xml格式
	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

func (this *WeixinClient) Run() {

	//解析post请求数据
	err := this.initMessage()

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	//微信消息是否是文本类型（string）
	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	//不同的消息类型
	switch MsgType {
	case "text":
		this.text()
		break
	case "voice":
		break;
	default:
		break
	}

	return
}