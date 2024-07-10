#pragma once
#include "function.h"

class NetworkNode;

class Connector {
public:
	NetworkNode* to;

};

class Connection : public Connector {
public:
	NetworkNode* from;

	Connection(NetworkNode* from, NetworkNode* to);
};

class MultiConnection : public Connector {
public:
	NetworkNode* left;
	NetworkNode* right;
	//Function* function;

	MultiConnection(NetworkNode* left, NetworkNode* right, NetworkNode* out);
};

