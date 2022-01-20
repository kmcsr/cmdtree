
package cmdtree

import (
	strings "strings"
	sort "sort"
)

func trimLeftStr(str string)(string){
	i := 0
	for i < len(str) && str[i] == ' ' {
		i++
	}
	return str[i:]
}

func splitNode(cmd string)(node string, _ string){
	i := strings.IndexByte(cmd, ' ')
	if i < 0 { i = len(cmd) }
	return cmd[:i], cmd[i:]
}

func strInList(str string, list []string)(bool){
	switch len(list) {
	case 0:
		return false
	case 1:
		return list[0] == str
	case 2:
		return list[0] == str || list[1] == str
	}
	i := sort.SearchStrings(list, str)
	return 0 <= i && i < len(list) && list[i] == str
}
