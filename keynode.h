#pragma once
#include "key.h"
#include <vector>

class KeyNode {
public:
	Key key;
	KeyNode* parent;
	std::vector<KeyNode*> children;
	
	KeyNode();
	
	template <typename T>
	KeyNode(T val) : key(val), parent(nullptr) {}

	template <typename T, typename U>
	KeyNode(T val, U parent) : key(val), parent(parent) {}


	bool isRoot() const;						// Returns true if the tree has no parent
	bool hasChildren() const;					// Returns false if the tree has no children
	bool isDescendantOf(Key* key);				// Returns true if the node has a parent with the given key
	bool isDescendantOf(KeyType key);			// Returns true if the node has a parent with the given key
	bool isAncestorOf(Key* key) const;			// Returns true if the node has a child with the given key
	bool isAncestorOf(KeyType key) const;		// Returns true if the node has a child with the given key
	bool hasChild(Key* key) const;				// Returns true if the node has a child with the given key
	bool hasChild(KeyType key) const;			// Returns true if the node has a child with the given key
	bool hasDescendant(int key) const;			// Returns true if the node has a descendant with the given key
	bool hasDescendant(char key) const;			// Returns true if the node has a descendant with the given key
	bool hasDescendant(std::string key) const;	// Returns true if the node has a descendant with the given key

	std::vector<KeyNode*> getChildren() const;	// Returns a vector of all children

	KeyNode* getParent() const;				// Returns the parent of the node
	KeyNode* getRoot();						// Returns the root of the tree
};
