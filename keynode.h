#pragma once
#include "key.h"
#include <vector>
#include <list>

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

	bool isRoot() const;						// Returns true if the node has no parent
	bool leaf() const;							// Returns true if the node has no children
	bool isDescendantOf(Key* key);				// Returns true if the node has a parent with the given key
	bool isDescendantOf(KeyType key);			// Returns true if the node has a parent with the given key
	bool hasChild(Key* key) const;				// Returns true if the node has a child with the given key
	bool hasChild(KeyType key) const;			// Returns true if the node has a child with the given key
	bool hasDescendant(int key) const;			// Returns true if the node has a descendant with the given key
	bool hasDescendant(char key) const;			// Returns true if the node has a descendant with the given key
	bool hasDescendant(std::string key) const;	// Returns true if the node has a descendant with the given key

	std::vector<KeyNode*> getChildren() const;	// Returns a vector of all children
	std::list<KeyNode*> getAncestors();			// Returns a list of all ancestors

	KeyNode* getParent() const;					// Returns the parent of the node
	KeyNode* getRoot();							// Returns the root of the tree
	KeyNode* findDescendant(Key* key);			// Returns the descendant with the given key
	KeyNode* findDescendant(KeyType key);		// Returns the descendant with the given key
};
