#pragma once
#include "key.h"
#include <vector>

class KeyNode {
public:
	Key key;
	KeyNode* parent;
	std::vector<KeyNode*> children;

	KeyNode() : key(Key()), parent(nullptr) {} 
	KeyNode(Key key) : key(key), parent(nullptr) {}
};
