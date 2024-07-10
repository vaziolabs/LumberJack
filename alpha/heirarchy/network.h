#pragma once

#include "key.h"
#include "value.h"
#include "networknode.h"
#include "value.h"
#include "keytree.h"
#include "connection.h"
#include <set>
#include <map>

class Network {
public:
	//KeyTree key_tree;								// Used to track Keys for Values
	std::set<NetworkNode*> nodes;					// Nodes in the network
	std::set<Connector*> connections;				// Connections between nodes

	Network();
	~Network();
	
	//void addNode(NetworkNode* node);
	/*void removeNode(NetworkNode* node);

	void addConnection(Connection* connection);
	void addConnection(MultiConnection* connection);
	void removeConnection(Connection* connection);
	void removeConnection(MultiConnection* connection);

	void connect(KeyType from, KeyType to);
	void connect(NetworkNode* from, NetworkNode* to);
	void disconnect(NetworkNode* from, NetworkNode* to);

	void print();*/
};
