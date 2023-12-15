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
	KeyTree key_tree;	// Used to track Keys for Values
	std::set<Value> values;							// Values in the network	

	std::set<Key> node_keys;						// Keys for nodes
	std::set<NetworkNode> nodes;					// Nodes in the network
	std::map<Key*, NetworkNode*> nodes_map;			// Key is the node's key used to access values
	
	std::set<Connector> connections;				// Connections between nodes
	std::map<Key*, Connector*> connections_map;		// Key is the from node's key

	Network() : key_tree() {}

	Network(KeyTree key_tree) : key_tree(key_tree) {}
	Network(KeyType root_key) : key_tree(root_key) {}
	
	NetworkNode* addNode(KeyType key);
	void addNode(NetworkNode* node);
	void removeNode(NetworkNode* node);

	void addConnection(Connection* connection);
	void addConnection(MultiConnection* connection);
	void removeConnection(Connection* connection);
	void removeConnection(MultiConnection* connection);

	void connect(KeyType from, KeyType to);
	void connect(NetworkNode* from, NetworkNode* to);
	void disconnect(NetworkNode* from, NetworkNode* to);

	void print();
};
