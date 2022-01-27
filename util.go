
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

func trimRightStr(str string)(string){
	i := len(str)
	for i > 0 && str[i - 1] == ' ' {
		i--
	}
	return str[:i]
}

func SplitNode(cmd string)(node string, _ string){
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

func delDuplication(list []string)([]string){
	if len(list) <= 1 {
		return list
	}
	sort.Strings(list)
	j := 0
	for i := 1; i < len(list); i++ {
		if list[i] != list[j] {
			j++
			if j < i {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	list = list[:j + 1]
	sort.Strings(list)
	return list
}
