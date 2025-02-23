package rapidgo

import "strings"

// node represents a radix tree node used for routing.
type node struct {
	path     string         // Static or wildcard path segment
	children []*node        // Child nodes
	handler  func(*Context) // Handler for the route, if applicable
	isWild   bool           // True if this node represents a wildcard (:param)
}

func (n *node) insert(path string, handler func(*Context)) {
	// Remove any leading/trailing slashes and split into segments.
	segments := strings.Split(strings.Trim(path, "/"), "/")
	current := n

	for _, segment := range segments {
		if segment == "" {
			continue // Skip empty segments.
		}

		var child *node
		// Look for an existing child node that exactly matches this segment.
		for _, c := range current.children {
			if c.path == segment {
				child = c
				break
			}
		}

		// If no child was found, create a new node for this segment.
		if child == nil {
			child = &node{
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

func (n *node) search(path string, params map[string]string) func(*Context) {
	// Split the path by "/" (ignoring empty segments).
	segments := strings.Split(strings.Trim(path, "/"), "/")
	current := n

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		var found *node
		// Check all children for a match.
		for _, child := range current.children {
			// If the segment exactly matches the child's path, choose it.
			if child.path == segment {
				found = child
				break
			}
			// Otherwise, if the child is a wildcard (e.g. ":id"), match it.
			if child.isWild {
				found = child
				// Extract the parameter value (exclude the ':' from the key).
				params[child.path[1:]] = segment
				break
			}
		}

		// If no matching child is found, return nil.
		if found == nil {
			return nil
		}
		current = found
	}

	return current.handler
}
