package rapidgo

import "strings"

// node represents a radix tree node used for routing.
type Node struct {
	path     string         // Static or wildcard path segment
	children []*Node        // Child nodes
	handler  func(*Context) // Handler for the route, if applicable
	isWild   bool           // True if this node represents a wildcard (:param)
}

func (n *Node) insert(path string, handler func(*Context)) {
	// Remove any leading/trailing slashes and split into segments.
	segments := strings.Split(strings.Trim(path, "/"), "/")
	current := n

	for _, segment := range segments {
		if segment == "" {
			continue // Skip empty segments.
		}

		var child *Node
		// Look for an existing child node that exactly matches this segment.
		for _, c := range current.children {
			if c.path == segment {
				child = c
				break
			}
		}

		// If no child was found, create a new node for this segment.
		if child == nil {
			child = &Node{
				path:   segment,
				isWild: len(segment) > 0 && segment[0] == ':',
			}
			current.children = append(current.children, child)
		}
		// Move to the child node for the next segment.
		current = child
	}

	// Once all segments have been processed, assign the handler to the final node.
	current.handler = handler
}

func (n *Node) search(path string, params map[string]string) func(*Context) {
	// Split the path by "/" (ignoring empty segments).
	segments := strings.Split(strings.Trim(path, "/"), "/")
	current := n

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		var found *Node
		var wildCard *Node
		// Check all children for a match.
		for _, child := range current.children {
			// Exact match takes precedence.
			if child.path == segment {
				found = child
				break
			}
			// Keep track of the first wildcard match as a fallback.
			if child.isWild && wildCard == nil {
				wildCard = child
			}
		}

		// If no matching child is found, return nil.
		// Use the exact match if found; otherwise, fall back to the wildcard.
		if found != nil {
			current = found
		} else if wildCard != nil {
			// Extract the parameter value (exclude the ':' from the key).
			params[wildCard.path[1:]] = segment
			current = wildCard
		} else {
			// No matching child found.
			return nil
		}
	}

	return current.handler
}
