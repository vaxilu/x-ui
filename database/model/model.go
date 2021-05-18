package model

import "time"

type Protocol string

const (
	VMess       Protocol = "vmess"
	VLESS       Protocol = "vless"
	Dokodemo    Protocol = "Dokodemo-door"
	Http        Protocol = "http"
	Trojan      Protocol = "trojan"
	Shadowsocks Protocol = "shadowsocks"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Inbound struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId     int       `json:"user_id"`
	Up         int64     `json:"up"`
	Down       int64     `json:"down"`
	Remark     string    `json:"remark"`
	Enable     bool      `json:"enable"`
	ExpiryTime time.Time `json:"expiry_time"`

	// config part
	Listen         string   `json:"listen"`
	Port           int      `json:"port"`
	Protocol       Protocol `json:"protocol"`
	Settings       string   `json:"settings"`
	StreamSettings string   `json:"stream_settings"`
	Tag            string   `json:"tag"`
	Sniffing       string   `json:"sniffing"`
}
