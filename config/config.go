package config

var (
	Ip      string
	Port    string
	Label   string
	Version = 1
)

func init() {
	Ip = "0.0.0.0"
	Port = "6666"
}
