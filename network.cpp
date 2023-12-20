#include "network.h"

Network::Network() {
	NetworkNode* root = new NetworkNode("/");
	
	this->nodes = { root };
	this->connections = {};
}

Network::~Network() {
	for (auto& node : nodes) {
		delete node;

	}

	for (auto& connection : connections) {
		delete connection;
	}
}

/*
void Network::addNode(NetworkNode* node) {
	// make sure node does not already exist
	if (nodes.contains(node)) {
		throw std::exception("Network addNode Error: Node already exists with that key type.");
		return;
	}

	nodes.insert(node);
}

void Network::removeNode(NetworkNode* node) {
	nodes.erase(*node);
	nodes_map.erase(&node->key);
}

void Network::addConnection(Connection* connection) {
	// make sure connection does not already exist
	if (connections_map.find(&connection->from->key) != connections_map.end()) {
		throw std::exception("Network addConnection Error: Connection already exists");
		return;
	}

	connections.insert(*connection);
	connections_map.insert({ &connection->from->key, connection });
}

void Network::addConnection(MultiConnection* connection) {
	// make sure connection does not already exist
	if (connections_map.find(&connection->left->key) != connections_map.end() ||
		connections_map.find(&connection->right->key) != connections_map.end()) {
		throw std::exception("Network addConnection Error: Connection already exists");
		return;
	}

	connections.insert(*connection);
	connections_map.insert({ &connection->left->key, connection });
	connections_map.insert({ &connection->right->key, connection });
}

void Network::removeConnection(Connection* connection) {
	connections.erase(*connection);
	connections_map.erase(&connection->from->key);
}

void Network::removeConnection(MultiConnection* connection) {
	connections.erase(*connection);
	connections_map.erase(&connection->left->key);
	connections_map.erase(&connection->right->key);
}

void Network::connect(KeyType from, KeyType to) {
	NetworkNode* from_node = nullptr;
	NetworkNode* to_node = nullptr;

	// find from node
	for (auto it = nodes_map.begin(); it != nodes_map.end(); it++) {
		if (it->first->value() == from) {
			from_node = it->second;
			break;
		}
	}

	// find to node
	for (auto it = nodes_map.begin(); it != nodes_map.end(); it++) {
		if (it->first->value() == to) {
			to_node = it->second;
			break;
		}
	}

	// make sure nodes exist
	if (from_node == nullptr) {
		throw std::exception("Network connect Error: 'from' Node does not exist.");
		return;
	}

	if (to_node == nullptr) {
		throw std::exception("Network connect Error: 'to' Node does not exist.");
		return;
	}

	// make sure connection does not already exist
	for (auto it = connections_map.begin(); it != connections_map.end(); it++) {
		if (it->first->value() == from) {
			throw std::exception("Network connect Error: Connection already exists");
			return;
		}
	}

	// create connection
	Connection* connection = new Connection(from_node, to_node);
	connections.insert(*connection);
	connections_map.insert({ &connection->from->key, connection });
}

void Network::connect(NetworkNode* from, NetworkNode* to) {
	// make sure nodes exist
	if (from == nullptr || to == nullptr) {
		throw std::exception("Network connect Error: Invalid input. Node cannot be a null type.");
		return;
	}

	// make sure connection does not already exist
	for (auto it = connections_map.begin(); it != connections_map.end(); it++) {
		if (it->first->value() == from->key.value()) {
			throw std::exception("Network connect Error: Connection already exists");
			return;
		}
	}

	// create connection
	Connection* connection = new Connection(from, to);
	connections.insert(*connection);
	connections_map.insert({ &connection->from->key, connection });
}

void Network::disconnect(NetworkNode* from, NetworkNode* to) {
	// make sure nodes exist
	if (from == nullptr || to == nullptr) {
		throw std::exception("Network disconnect Error: Node does not exist");
		return;
	}

	// make sure connection exists
	for (auto it = connections_map.begin(); it != connections_map.end(); it++) {
		if (it->first->value() == from->key.value()) {
			throw std::exception("Network disconnect Error: Connection does not exist");
			return;
		}
	}

	// remove connection
	connections.erase(*connections_map[&from->key]);
	connections_map.erase(&from->key);
}

void Network::print() {
	std::cout << "Network:" << std::endl;
	std::cout << "\tNodes:" << std::endl;
	for (NetworkNode node : nodes) {
		node.print();
		std::cout << std::endl;
	}

	std::cout << "\tConnections:" << std::endl;
	for (auto& connector : connections) {
		// check if connection is a MultiConnection or Connection
		if (typeid(connector) == typeid(Connection)) {
			Connection* c = (Connection*)&connector;
			std::cout << "\t\t" << c->from->key << " -> " << c->to->key << std::endl;
		}
		else if (typeid(connector) == typeid(MultiConnection)) {
			MultiConnection* c = (MultiConnection*)&connector;
			std::cout << "\t\t" << c->left->key << " -> " << c->to->key << std::endl;
			std::cout << "\t\t" << c->right->key << " -> " << c->to->key << std::endl;
		}
	}
}
*/