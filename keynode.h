#pragma once
#include "key.h"
#include <vector>

class KeyNode {
public:
	Key key;
	KeyNode* parent;
	std::vector<KeyNode*> children;
	
	KeyNode();
	KeyNode(Key key);
	KeyNode(Key key, KeyNode* parent);
};
