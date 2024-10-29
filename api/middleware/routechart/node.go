package routechart

import "github.com/go-chujang/demo-aa/api/ctxutil"

type node struct {
	children map[string]*node
	policy   *ctxutil.Policy
}

func newNode() *node {
	return &node{
		children: make(map[string]*node),
	}
}

func (n *node) insert(pathSegments []string, policy *ctxutil.Policy) {
	current := n
	for _, segment := range pathSegments {
		if len(segment) > 0 && segment[0] == ':' {
			segment = "*"
		}
		if _, exists := current.children[segment]; !exists {
			current.children[segment] = newNode()
		}
		current = current.children[segment]
	}
	current.policy = policy
}

func (n *node) search(pathSegments []string) (*ctxutil.Policy, bool) {
	current := n
	for _, segment := range pathSegments {
		if child, exists := current.children[segment]; exists {
			current = child
		} else if child, exists := current.children["*"]; exists {
			current = child
		} else {
			return nil, false
		}
	}
	return current.policy, true
}
