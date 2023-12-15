#include "connection.h"
Connection::Connection(NetworkNode* from, NetworkNode* to) {
	this->from = from;
	this->to = to;
}

MultiConnection::MultiConnection(NetworkNode* left, NetworkNode* right, NetworkNode* out) {
	this->left = left;
	this->right = right;
	this->to = out;
}
/*
*/