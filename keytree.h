#pragma once
#include "key.h"
#include "keynode.h"
#include <string>


class KeyTree {
public:
	KeyNode* root;

	KeyTree() : root(new KeyNode()) {}
	KeyTree(KeyNode* root) : root(root) {}
	KeyTree(KeyType root_key) : root(new KeyNode(root_key)) {}

	void remove(int id) const;
	void remove(char id) const;
	void remove(std::string id) const;

	KeyNode* insert(KeyNode *key, KeyType parent_key);
	KeyNode* search(KeyType key) const;
	
	void print() const;
};
