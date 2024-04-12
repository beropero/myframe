package kilon

import (
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的当前部分
	children []*node // 孩子节点列表
	isWild   bool    // 是否为模糊搜索，当含有通配符":"或者"*"时为true
}

// 插入
func (n *node) insert(pattern string, parts []string, index int) {
	// 进来的时候 n.part = parts[index-1] 即最后一个 part 则 直接设置patten
	if len(parts) == index {
		n.pattern = pattern
		return
	}
	// 还需匹配part
	// 先在n的children列表中匹配part
	part := parts[index]
	child := n.matchChild(part)
	// 如果没有找到，则构建一个child并插入n.children列表中
	if child == nil {
		child = &node{
			part: part,
			// 含有通配符":"或者"*"时为true
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 插入n.children列表
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, index+1)
}

// 查找匹配child
func (n *node) matchChild(part string) *node {
	// 遍历 n.children 查找part相同的
	for _, child := range n.children {
		// 如果找到匹配返回child
		if child.part == part || child.isWild {
			return child
		}
	}
	// 没找到返回nil
	return nil
}

// 搜索
func (n *node) search(parts []string, index int) *node {
	// 如果匹配将节点返回
	if len(parts) == index || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[index]
	// 获取匹配的所有孩子节点
	nodes := n.matchChildren(part)

	// 递归搜索匹配的child节点
	for _, child := range nodes {
		result := child.search(parts, index+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 查找匹配的孩子节点，由于有通配符，所以可能会有多个匹配，因此返回一个节点列表
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child) // 将符合的孩子节点添入返回列表
		}
	}
	return nodes
}
