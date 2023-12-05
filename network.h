#pragma once

#include "key.h"
#include "value.h"
#include "networknode.h"
#include "value.h"
#include "keytree.h"
#include <set>
#include <map>

class Network {
public:
	KeyTree key_tree;
	std::set<Value> values;
	std::set<NetworkNode> nodes;
	std::set<Connection> connections;
	std::map<Key, std::set<Connection>*> connections_map;

	Network() : key_tree() {}
/*
	void addNode(NetworkNode* node);
	void removeNode(NetworkNode* node);

	void addConnection(Connection* connection);
	void removeConnection(Connection* connection);

	void connect(Key from, Key to);
	void disconnect(Key from, Key to);

	void print();
	*/
};
