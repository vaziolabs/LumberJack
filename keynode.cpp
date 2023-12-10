#include "keynode.h"

KeyNode::KeyNode() { this->key = NULL; this->parent = nullptr; }

bool KeyNode::isRoot() const { return this->parent == nullptr; }
bool KeyNode::leaf() const { return this->children.size() > 0; }
bool KeyNode::isDescendantOf(Key* key) { return this->isDescendantOf(key->getValue()); }
bool KeyNode::hasChild(Key* key) const { return this->hasChild(key->getValue()); }

std::vector<KeyNode*> KeyNode::getChildren() const { return this->children; }

KeyNode* KeyNode::getParent() const { return this->parent; }
KeyNode* KeyNode::findDescendant(Key* key) { return this->findDescendant(key->getValue()); }

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

KeyNode* KeyNode::findDescendant(KeyType key) {
	if (this->key == key) {
		return this;
	}

	std::vector<KeyNode*> childs = this->children;

	while (childs.size() > 0) {
		std::vector<KeyNode*> newChilds;

		for (int i = 0; i < childs.size(); i++) {
			if (childs[i]->key == key) {
				return childs[i];
			}

			for (int j = 0; j < childs[i]->children.size(); j++) {
				newChilds.push_back(childs[i]->children[j]);
			}
		}

		childs = newChilds;
	}

	return nullptr;
}
std::list<KeyNode*> KeyNode::getAncestors()
{
	std::list<KeyNode*> ancestors;

	KeyNode* current = this;

	while (current->parent != nullptr) {
		ancestors.push_back(current->parent);
		current = current->parent;
	}

	return ancestors;
}
