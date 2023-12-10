#include "keynode.h"

KeyNode::KeyNode() {
	this->key = NULL;
	this->parent = nullptr;
}

KeyNode::KeyNode(Key key) {
	this->key = key;
	this->parent = nullptr;
}

KeyNode::KeyNode(Key key, KeyNode* parent) {
	this->key = key;
	this->parent = parent;
}