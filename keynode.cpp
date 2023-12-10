#include "keynode.h"

KeyNode::KeyNode() { this->key = NULL; this->parent = nullptr; }

bool KeyNode::isRoot() const { return this->parent == nullptr; }
bool KeyNode::hasChildren() const { return this->children.size() > 0; }
bool KeyNode::isDescendantOf(Key* key) { return this->isDescendantOf(key->getValue()); }
bool KeyNode::hasChild(Key* key) const { return this->hasChild(key->getValue()); }
bool KeyNode::isAncestorOf(Key* key)  const{ return this->isAncestorOf(key->getValue()); }

std::vector<KeyNode*> KeyNode::getChildren() const { return this->children; }

KeyNode* KeyNode::getParent() const { return this->parent; }

bool KeyNode::isDescendantOf(KeyType key) {
	if (this->isRoot()) {
		return false;
	}

	KeyNode* current = this;
	 
	while (current != nullptr) {
		if (current->key == key) {
			return true;
		}

		current = current->parent;
	}

	return false;
}

bool KeyNode::hasChild(KeyType key) const {
	for (int i = 0; i < this->children.size(); i++) {
		if (this->children[i]->key == key) {
			return true;
		}
	}

	return false;
}

bool KeyNode::isAncestorOf(KeyType key) const {
	// needs to iterate through all children of all children
	std::vector<KeyNode*> childs = this->children;
	
	// if any of them have the key, return true
	if (this->hasChild(key)) {
		return true;
	}

	while (childs.size() > 0) {
		std::vector<KeyNode*> newChilds;

		// for each child, check if it has the key
		for (int i = 0; i < childs.size(); i++) {
			// if it does, return true
			if (childs[i]->hasChild(key)) {
				return true;
			}

			// else add its children to the newChilds vector
			for (int j = 0; j < childs[i]->children.size(); j++) {
				newChilds.push_back(childs[i]->children[j]);
			}
		} 

		// set childs to newChilds and repeat until childs is empty
		childs = newChilds;
	}
	
	// else return false
	return false;
}

KeyNode* KeyNode::getRoot() {
	KeyNode* current = this;

	while (current->parent != nullptr) {
		current = current->parent;
	}

	return current;
}

bool KeyNode::hasDescendant(int key) const {
	if (this->key == key) {
		return true;
	}

	std::vector<KeyNode*> childs = this->children;

	while (childs.size() > 0) {
		std::vector<KeyNode*> newChilds;

		for (int i = 0; i < childs.size(); i++) {
			if (childs[i]->key == key) {
				return true;
			}

			for (int j = 0; j < childs[i]->children.size(); j++) {
				newChilds.push_back(childs[i]->children[j]);
			}
		}

		childs = newChilds;
	}

	return false;
}

bool KeyNode::hasDescendant(char key) const {
	if (this->key == key) {
		return true;
	}

	std::vector<KeyNode*> childs = this->children;

	while (childs.size() > 0) {
		std::vector<KeyNode*> newChilds;

		for (int i = 0; i < childs.size(); i++) {
			if (childs[i]->key == key) {
				return true;
			}

			for (int j = 0; j < childs[i]->children.size(); j++) {
				newChilds.push_back(childs[i]->children[j]);
			}
		}

		childs = newChilds;
	}

	return false;
}

bool KeyNode::hasDescendant(std::string key) const {
	if (this->key == key) {
		return true;
	}

	std::vector<KeyNode*> childs = this->children;

	while (childs.size() > 0) {
		std::vector<KeyNode*> newChilds;

		for (int i = 0; i < childs.size(); i++) {
			if (childs[i]->key == key) {
				return true;
			}

			for (int j = 0; j < childs[i]->children.size(); j++) {
				newChilds.push_back(childs[i]->children[j]);
			}
		}

		childs = newChilds;
	}

	return false;
}