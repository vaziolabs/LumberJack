package core

import "time"

// AddActivity adds an activity entry to the node
func (n *Node) AddActivity(content interface{}, metadata map[string]interface{}, userID string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	entry := Entry{
		Content:   content,
		Metadata:  metadata,
		UserID:    userID,
		Timestamp: time.Now(),
	}
	n.Entries = append(n.Entries, entry)
}

// AddChild adds a child node with proper parent linking
func (n *Node) AddChild(child *Node) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	child.AddParent(n)
	n.Children[child.ID] = child
	return nil
}

// AddParent adds a parent node to the node
func (n *Node) AddParent(parent *Node) {
	if n.Parents == nil {
		n.Parents = make(map[string]string)
	}
	n.Parents[parent.ID] = parent.Name
}
