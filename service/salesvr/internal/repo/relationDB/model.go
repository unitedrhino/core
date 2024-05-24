package relationDB

import (
	"database/sql"
	"time"
)

// 示例
type SaleExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

type SaleDistributionBrokerageDetail struct {
	ID                     int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	SerialNumber           string    `gorm:"column:serial_number;type:varchar(64);comment:设备ID" json:"serial_number"`
	CardNumber             string    `gorm:"column:card_number;type:varchar(64);comment:物联网卡号" json:"card_number"`
	SaleDistributionUserID int64     `gorm:"column:distribution_user_id;type:bigint(20);comment:分销商用户ID" json:"sale_distribution_user_id"`
	SaleDistributionName   string    `gorm:"column:distribution_name;type:varchar(100);comment:分销公司名称" json:"sale_distribution_name"`
	GoodID                 int64     `gorm:"column:good_id;type:bigint(20);comment:套餐ID" json:"good_id"`
	GoodName               string    `gorm:"column:good_name;type:varchar(100);comment:套餐名称" json:"good_name"`
	UserID                 int64     `gorm:"column:user_id;type:bigint(20);comment:用户ID" json:"user_id"`
	PayPrice               int64     `gorm:"column:pay_price;type:decimal(10,2);comment:实际付款金额" json:"pay_price"`
	Brokerage              int64     `gorm:"column:brokerage;type:decimal(10,2);comment:获得佣金;NOT NULL" json:"brokerage"`
	CreateTime             time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_time"`
}

func (m *SaleDistributionBrokerageDetail) TableName() string {
	return "sale_distribution_brokerage_detail"
}

type SaleDistributionEarning struct {
	ID                       int64        `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Money                    int64        `gorm:"column:money;type:decimal(10,2);comment:金额" json:"money"`
	SaleDistributionUserID   int64        `gorm:"column:distribution_user_id;type:bigint(20);comment:分销用户ID" json:"sale_distribution_user_id"`
	SaleDistributionUserName int64        `gorm:"column:distribution_user_name;type:bigint(50);comment:分销用户名称" json:"sale_distribution_user_name"`
	CompanyName              string       `gorm:"column:company_name;type:varchar(255);comment:公司名称;NOT NULL" json:"company_name"`
	BankName                 string       `gorm:"column:bank_name;type:varchar(255);comment:开户银行" json:"bank_name"`
	CorporateAccount         string       `gorm:"column:corporate_account;type:varchar(255);comment:银行账户" json:"corporate_account"`
	ApplicationType          int64        `gorm:"column:application_type;type:tinyint;default:1;comment:审核结果(1=发起审核,2=一级审核通过,3=管理员审核通过,4=管理员拒绝,5=上级经销商拒绝)" json:"application_type"`
	ReviewContent            string       `gorm:"column:review_content;type:varchar(255);comment:审核内容批语记录" json:"review_content"`
	TopUserID                int64        `gorm:"column:top_user_id;type:bigint(20);default:1;comment:上级经销商ID" json:"top_user_id"`
	TopUserName              string       `gorm:"column:top_user_name;type:varchar(50);comment:上级经销商名称" json:"top_user_name"`
	CreateTime               sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
}

func (m *SaleDistributionEarning) TableName() string {
	return "sale_distribution_earning"
}

// 分销商品表
type SaleDistributionGoods struct {
	GoodID            int64        `gorm:"column:good_id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:商品ID" json:"good_id"`
	GoodName          string       `gorm:"column:good_name;type:varchar(50);comment:商品名称;NOT NULL" json:"good_name"`
	SitType           int64        `gorm:"column:sit_type;type:enum('1','2','3','4');default:1;comment:套餐类型:1计时，2计量，3语音，4短信" json:"sit_type"`
	SitAmount         int64        `gorm:"column:sit_amount;type:int(11);default:0;comment:套餐内容" json:"sit_amount"`
	UserID            int64        `gorm:"column:user_id;type:bigint(20);comment:创建人ID" json:"user_id"`
	GoodSort          int64        `gorm:"column:good_sort;type:int(11);comment:排序" json:"good_sort"`
	GoodApply         int64        `gorm:"column:good_apply;type:enum('product','user');default:product;comment:商品适用类型(product=设备产品,user=用户套餐)" json:"good_apply"`
	GoodApplyID       string       `gorm:"column:good_apply_id;type:varchar(200);comment:商品适应类型ID" json:"good_apply_id"`
	GoodType          int64        `gorm:"column:good_type;type:enum('ordinary','activity','system');default:ordinary;comment:商品类型(ordinary=普通商品,activity=活动商品system=系统商品)" json:"good_type"`
	ActivityStartTime sql.NullTime `gorm:"column:activity_start_time;type:datetime;comment:活动商品开始时间" json:"activity_start_time"`
	ActivityEndTime   sql.NullTime `gorm:"column:activity_end_time;type:datetime;comment:活动商品结束时间" json:"activity_end_time"`
	Limited           int64        `gorm:"column:limited;type:int(11);comment:限购数量" json:"limited"`
	CreateTime        sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime        sql.NullTime `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
}

type SaleDistributionGoodsPrice struct {
	ID                     int64        `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	GoodID                 int64        `gorm:"column:good_id;type:bigint(20);comment:商品ID;NOT NULL" json:"good_id"`
	CostPrice              int64        `gorm:"column:cost_price;type:decimal(10,2);comment:结算价(上级结算价)" json:"cost_price"`
	SalePrice              int64        `gorm:"column:sale_price;type:decimal(10,2);comment:销售价(下级结算价)" json:"sale_price"`
	SaleFakePrice          int64        `gorm:"column:sale_fake_price;type:decimal(10,2);comment:销售划线价" json:"sale_fake_price"`
	SaleDistributionUserID int64        `gorm:"column:distribution_user_id;type:bigint(20);default:0;comment:经销商ID(0为设置销售价)" json:"sale_distribution_user_id"`
	UserID                 int64        `gorm:"column:user_id;type:bigint(20);comment:创建人ID" json:"user_id"`
	CreateTime             sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	Status                 int64        `gorm:"column:status;type:tinyint(1);default:1;comment:(1=上架,2=下架)" json:"status"`
	SaleType               int64        `gorm:"column:sale_type;type:enum('fx','my');default:my;comment:销售类型(1=自销,2=分销)" json:"sale_type"`
	SaleUserID             int64        `gorm:"column:sale_user_id;type:bigint(20);comment:设置用户ID" json:"sale_user_id"`
}

func (m *SaleDistributionGoodsPrice) TableName() string {
	return "sale_distribution_goods_price"
}

type SaleDistributionOrder struct {
	ID                       int64        `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	OrderID                  string       `gorm:"column:order_id;type:varchar(64);comment:订单ID;NOT NULL" json:"order_id"`
	ThirdOrderID             string       `gorm:"column:third_order_id;type:varchar(64);comment:第三方支付订单" json:"third_order_id"`
	ThirdInfo                int64        `gorm:"column:third_info;type:enum('1','2');default:1;comment:第三方(1=微信,2=支付宝)" json:"third_info"`
	SerialNumber             string       `gorm:"column:serial_number;type:varchar(64);comment:设备编号" json:"serial_number"`
	CardNumber               string       `gorm:"column:card_number;type:varchar(64);comment:物联网卡号" json:"card_number"`
	PayPrice                 int64        `gorm:"column:pay_price;type:decimal(10,2);comment:支付金额;NOT NULL" json:"pay_price"`
	GoodID                   int64        `gorm:"column:good_id;type:bigint(20);comment:套餐ID;NOT NULL" json:"good_id"`
	GoodName                 string       `gorm:"column:good_name;type:varchar(100);comment:套餐名称" json:"good_name"`
	GoodPrice                int64        `gorm:"column:good_price;type:decimal(10,2);comment:套餐价格;NOT NULL" json:"good_price"`
	SaleDistributionUserID   int64        `gorm:"column:distribution_user_id;type:bigint(20);comment:分销用户ID" json:"sale_distribution_user_id"`
	SaleDistributionUserName string       `gorm:"column:distribution_user_name;type:varchar(100);comment:分销用户名称" json:"sale_distribution_user_name"`
	UserID                   int64        `gorm:"column:user_id;type:bigint(20);comment:支付用户ID;NOT NULL" json:"user_id"`
	UserName                 string       `gorm:"column:user_name;type:varchar(100);comment:支付用户名称" json:"user_name"`
	CreateTime               sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	ActiveTime               sql.NullTime `gorm:"column:active_time;type:datetime;comment:生效时间" json:"active_time"`
	EndTime                  sql.NullTime `gorm:"column:end_time;type:datetime;comment:到期时间" json:"end_time"`
	Status                   int64        `gorm:"column:status;type:enum('1','2','3');default:1;comment:状态(1=待支付,2=已支付,3=支付失败)" json:"status"`
	TicketID                 int64        `gorm:"column:ticket_id;type:bigint(20);comment:优惠券ID" json:"ticket_id"`
	TicketName               string       `gorm:"column:ticket_name;type:varchar(100);comment:优惠券名称" json:"ticket_name"`
	TicketPrice              int64        `gorm:"column:ticket_price;type:decimal(10,2);comment:优惠券价格" json:"ticket_price"`
}

func (m *SaleDistributionOrder) TableName() string {
	return "sale_distribution_order"
}

// 分销优惠券表
type SaleDistributionTicket struct {
	TicketID        int64        `gorm:"column:ticket_id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:优惠券ID" json:"ticket_id"`
	TicketName      string       `gorm:"column:ticket_name;type:varchar(50);comment:优惠券名称" json:"ticket_name"`
	Threshold       int64        `gorm:"column:threshold;type:decimal(11,2);default:0.00;comment:使用门槛" json:"threshold"`
	Discount        int64        `gorm:"column:discount;type:decimal(11,2);comment:优惠价格" json:"discount"`
	TicketApplyType int64        `gorm:"column:ticket_apply_type;type:enum('0','1','2');default:0;comment:适用对象(0设备,1=用户,2=通用)" json:"ticket_apply_type"`
	TicketApplyID   string       `gorm:"column:ticket_apply_id;type:varchar(200);comment:适用对象ID" json:"ticket_apply_id"`
	TicketGoodID    string       `gorm:"column:ticket_good_id;type:varchar(200);comment:适用套餐ID" json:"ticket_good_id"`
	Timing          int64        `gorm:"column:timing;type:enum('1','2');default:1;comment:下发时机(1=设备绑定触发,2=自主下发)" json:"timing"`
	CreateTime      sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	EndTime         sql.NullTime `gorm:"column:end_time;type:datetime;comment:到期时间" json:"end_time"`
}

func (m *SaleDistributionTicket) TableName() string {
	return "sale_distribution_ticket"
}

// 经销商入驻申请表
type SaleDistributionUserApply struct {
	ID               int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	UserID           int64     `gorm:"column:user_id;type:bigint(20);comment:用户ID主键;NOT NULL" json:"user_id"`
	UserName         string    `gorm:"column:user_name;type:varchar(255);comment:用户名称;NOT NULL" json:"user_name"`
	CompanyName      string    `gorm:"column:company_name;type:varchar(255);comment:公司名称;NOT NULL" json:"company_name"`
	LegalName        string    `gorm:"column:legal_name;type:varchar(255);comment:公司法人名字;NOT NULL" json:"legal_name"`
	Mobile           string    `gorm:"column:mobile;type:varchar(16);comment:联系电话" json:"mobile"`
	CompanyCode      string    `gorm:"column:company_code;type:varchar(255);comment:统一社会信用代码;NOT NULL" json:"company_code"`
	CompanyPath      string    `gorm:"column:company_path;type:varchar(255);comment:公司地址" json:"company_path"`
	BankName         string    `gorm:"column:bank_name;type:varchar(255);comment:开户银行" json:"bank_name"`
	CorporateAccount string    `gorm:"column:corporate_account;type:varchar(255);comment:银行账户" json:"corporate_account"`
	ParentID         int64     `gorm:"column:parent_id;type:bigint(20);default:0;comment:上级id" json:"parent_id"`
	AccountType      int64     `gorm:"column:account_type;type:int(11);comment:状态：有效，无效，注销，封号" json:"account_type"`
	ApplicationType  int64     `gorm:"column:application_type;type:int(11);comment:申请是否通过" json:"application_type"`
	ReviewContent    string    `gorm:"column:review_content;type:varchar(255);comment:审核内容批语记录" json:"review_content"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;NOT NULL" json:"create_time"`
	ApplicationTime  time.Time `gorm:"column:application_time;type:datetime;default:CURRENT_TIMESTAMP;comment:通过申请时间;NOT NULL" json:"application_time"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"update_time"`
	Remark           int64     `gorm:"column:remark;type:bigint(20) unsigned;comment:分销系统用户表ID" json:"remark"`
}

func (m *SaleDistributionUserApply) TableName() string {
	return "sale_distribution_user_apply"
}

// 分销用户佣金
type SaleDistributionUserBrokerage struct {
	ID                       int64        `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:ID" json:"id"`
	SaleDistributionUserID   int64        `gorm:"column:distribution_user_id;type:bigint(20);comment:分销用户ID;NOT NULL" json:"sale_distribution_user_id"`
	SaleDistributionUserName string       `gorm:"column:distribution_user_name;type:varchar(100);comment:分销用户名" json:"sale_distribution_user_name"`
	AllBrokerage             int64        `gorm:"column:all_brokerage;type:decimal(10,2);default:0.00;comment:获得的总佣金" json:"all_brokerage"`
	AlreadyBrokerage         int64        `gorm:"column:already_brokerage;type:decimal(10,2);default:0.00;comment:已提现的佣金" json:"already_brokerage"`
	CreateTime               sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
}

func (m *SaleDistributionUserBrokerage) TableName() string {
	return "sale_distribution_user_brokerage"
}

type SaleDistributionUserDevice struct {
	UserID                int64  `gorm:"column:user_id;type:bigint(20);primary_key;comment:分销商ID" json:"user_id"`
	DeviceSerialNumber    string `gorm:"column:device_serial_number;type:varchar(64);comment:设备表设备编号;NOT NULL" json:"device_serial_number"`
	SaleDistributionLevel string `gorm:"column:distribution_level;type:varchar(1);default:1;comment:分销等级" json:"sale_distribution_level"`
}

func (m *SaleDistributionUserDevice) TableName() string {
	return "sale_distribution_user_device"
}

// 分销流水表
type SaleDistributionWater struct {
	WaterID                int64        `gorm:"column:water_id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:流水ID" json:"water_id"`
	DeviceID               int64        `gorm:"column:device_id;type:bigint(20);comment:设备ID" json:"device_id"`
	Typeof                 int64        `gorm:"column:typeof;type:enum('0','1');default:0;comment:类型(1=设备套餐流水,2=用户套餐流水)" json:"typeof"`
	GoodsID                int64        `gorm:"column:goods_id;type:bigint(20);comment:商品ID" json:"goods_id"`
	GoodsPrice             int64        `gorm:"column:goods_price;type:decimal(10,2);comment:商品原价" json:"goods_price"`
	GoodsName              string       `gorm:"column:goods_name;type:varchar(50);comment:商品名称" json:"goods_name"`
	ActualPrice            int64        `gorm:"column:actual_price;type:decimal(10,2);comment:实际付款价格" json:"actual_price"`
	TickerID               int64        `gorm:"column:ticker_id;type:bigint(20);comment:优惠券ID" json:"ticker_id"`
	TicketPrice            int64        `gorm:"column:ticket_price;type:decimal(10,2);comment:抵扣金额" json:"ticket_price"`
	TicketName             string       `gorm:"column:ticket_name;type:varchar(50);comment:优惠券名称" json:"ticket_name"`
	UserID                 int64        `gorm:"column:user_id;type:bigint(20);comment:用户ID" json:"user_id"`
	UserName               string       `gorm:"column:user_name;type:varchar(50);comment:用户名称" json:"user_name"`
	SaleDistributionUserID int64        `gorm:"column:distribution_user_id;type:bigint(20);comment:分销商ID" json:"sale_distribution_user_id"`
	CreateTime             sql.NullTime `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
}

func (m *SaleDistributionWater) TableName() string {
	return "sale_distribution_water"
}
