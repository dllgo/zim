package utils

import "net"

const (
	//# 对称加密key
	CryptoKey = "Adba723b7fe06819"
)
 
//对称加密IP和端口，当做clientId
func GenCid() string {
	raw := []byte(getIntranetIp() + ":" + "1024")
	str, err := Encrypt(raw, []byte(CryptoKey))
	if err != nil {
		panic(err)
	}
	return str
}

//获取client key地址信息
func GetAddrByCid(cId string) (addr string) {
	//解密ClientId
	addr, _= Decrypt(cId, []byte(CryptoKey))
	
	return
}

//获取本机内网IP
func getIntranetIp() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}

	return ""
}