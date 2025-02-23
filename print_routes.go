package rapidgo

import "fmt"

// PrintRoutes prints all registered routes in a tree structure
func (e *Engine) PrintRoutes() {
	fmt.Println("\nRegistered Routes:")
	fmt.Println("================")
	for method, root := range e.Router.trees {
		fmt.Printf("\n[%s]\n", method)
		printNode(root, "", true)
	}
}

// Helper function to print the node tree
func printNode(n *node, prefix string, isLast bool) {
	if n == nil {
		return
	}

	// Print current node
	marker := "├──"
	if isLast {
		marker = "└──"
	}

	// Show handler presence with [handler]
	handlerMark := " "
	if n.handler != nil {
		handlerMark = " [handler]"
	}

	// Print node information
	fmt.Printf("%s%s %s%s\n", prefix, marker, n.path, handlerMark)

	// Prepare prefix for children
	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	// Print children
	for i, child := range n.children {
		isLastChild := i == len(n.children)-1
		printNode(child, childPrefix, isLastChild)
	}
}
