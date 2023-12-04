package crypto

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type MD5 struct {
}

func (*MD5) Crypt(v1 []byte, v2 []byte, uppercase bool) string {
	var rs = make([]string, 0)
	m := md5.New()
	m.Write(v1)
	bm := m.Sum(v2)
	for _, v := range bm {
		if uppercase {
			rs = append(rs, fmt.Sprintf("%02X", v))
		} else {
			rs = append(rs, fmt.Sprintf("%02x", v))
		}
	}
	return strings.Join(rs, "")
	/*return hex.EncodeToString(h.Sum(format.ToByte(v2)))*/ //第二种返回字符串的方法，返回的参数是小写
}
