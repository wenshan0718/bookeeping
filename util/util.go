//工具
package util

import (
	"strings"
)

/**
 * 读取cookie
 * cookies http.header.Cookie 中的字符串
	name cookie key的名称
*/
func GetCookieByName(cookies string, name string) string {
	if len(cookies) > 0 {
		c_start := strings.Index(cookies, name+"=")
		if c_start != -1 {
			c_start = c_start + len(name) + 1
			c_end := strings.Index(string([]byte(cookies)[c_start:]), ";")
			if c_end == -1 {
				c_end = len(cookies)
			} else {
				c_end = c_end + c_start
			}
			return string([]byte(cookies)[c_start:c_end])
		}
	}
	return ""
}
