package gee

import (
	"fmt"
	"log"
)

type TrieNode struct {
	part string // 记录路由中的一部分 比如 hello
	next map[string]*TrieNode // 记录它所有的子节点
	isPath bool // 记录是否是一个路由
	isWild bool // 记录是否通配
	path string // 记录路径
}

func NewTrie() *TrieNode {
	return &TrieNode{
		next: make(map[string]*TrieNode),
		isPath: false,
	}
}

// 查找当前节点满足part的所有子节点
func (this *TrieNode) findChildren(part string) []*TrieNode {
	nodes := make([]*TrieNode, 0)
	for _, v := range this.next {
		if v.part == part || v.isWild {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

// Search  用于查找是否具有该路径
// 可能存在多个通配符情况 如 /hello/:book/go /hello/:name/Tom
// 所以需要精确找到符合路径的节点  考虑采用深度优先  有可能找第一次就找到了正确地路径
func (this *TrieNode) Search(parts []string) *TrieNode{
	res := new(TrieNode)
	found := false
	fmt.Printf("parts len is %v, parts is %v\n", len(parts), parts)
	var dfs func(i int, parts []string, node *TrieNode)
	dfs = func (i int, parts []string, node *TrieNode) {
		if i == len(parts) {
			fmt.Printf("进入第%v个part,当前node是%v\n",i,node)
			found = true
			return
		}
		part := parts[i]
		nodes := node.findChildren(part)
		fmt.Printf("i is %v,part is %v nodes is %v\n", i, part,nodes)
		for _, v := range nodes {
			fmt.Printf("log v %v::\n", v)
			// 找到了
			if found {
				return
			}
			if v.part == part || v.isWild {
				res = v
				if v.part == "*" {
					found = true
					return
				}
				dfs(i + 1, parts, v)
			} else {
				continue
			}
		}
	}
	for i := 0; i < len(parts); i++ {
		if !found {
			dfs(i, parts, this)
		}
	}
	if found && res.isPath {
		return res
	}
	return nil
} 


// 向trie树添加路由信息
// 规则：如果此路径下已经存在*通配符，则此路径下不允许添加任何URL
// 如果已经有其他路径，则不允许添加*通配符
func (this *TrieNode) Insert(parts []string, path string) {
	// fmt.Printf("path is %s len is %v",path, parts)
	for _, v := range parts {
		if v != "" {
			if v[0] == '*' {
				v = "*"
				if len(this.next) > 0 {
					log.Println("当前目录下已经存在路径，无法申请*通配符路径")
					return
				}
			}	
		}
		if this.next["*"] != nil {
			log.Println("当前目录下已经存在*通配符路径，无法申请路径")
			return
		}
		if this.next[v] == nil {
			node := &TrieNode{
				part: v,
				next: make(map[string]*TrieNode),
				isWild: false,
			}
			if v != "" {
				if v[0] == ':' || v[0] == '*' {
					node.isWild = true
				}
			}
			this.next[v] = node
		}
		this = this.next[v]
		// fmt.Printf("gee insert %s, this: %v\n", parts[len(parts) - 1], this)
	} 
	this.isPath = true
	this.path = path
}