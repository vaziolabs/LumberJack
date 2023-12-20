#include "networknode.h"

NetworkNode::NetworkNode() {
	this->key = NULL;
	this->values = std::vector<Value>();
	this->connections = std::set<Connector*>();
}
std::string NetworkNode::getKeyType() const { return key.type(); }

 void NetworkNode::addValue(ValueType value) {
	values.push_back(Value(value));
}

void NetworkNode::deleteValue(int index) {
	values.erase(values.begin() + index);
}

int NetworkNode::getIndex(Value value) {
	auto it = std::find(values.begin(), values.end(), value);

	if (it != values.end()) {
		return it - values.begin();
	}

	return -1;
}

int NetworkNode::getIndex(ValueType value) {
	auto it = std::find_if(values.begin(), values.end(), [value](Value v) { return v.value() == value; });

	if (it != values.end()) {
		return it - values.begin();
	}

	return -1;
} 
/**/
void NetworkNode::addConnection(Connector* node) {
	connections.insert(node);
}

void NetworkNode::removeConnection(Connector* node) {
	connections.erase(node);
}
void NetworkNode::print() {
	std::cout << "\tNode Key: " << this->key << std::endl;
	
	std::cout << "\tValues:" << std::endl;
	for (int i = 0; i < this->values.size(); i++) {
		std::cout << "\t\t[" << i << "] - " << this->values[i] << std::endl;
	}

	std::cout << "\tConnections:" << std::endl;
	for (auto &connection : this->connections) {
		std::cout << "\t\t-> " << connection->to->key << std::endl;
	}
}

inline std::ostream& operator<<(std::ostream& os, const NetworkNode& obj) {
	os << obj.key;
	return os;
}

inline bool operator==(const NetworkNode& lhs, const NetworkNode& rhs) {
	return lhs.key == rhs.key;
}

inline bool operator!=(const NetworkNode& lhs, const NetworkNode& rhs) {
	return !(lhs.key == rhs.key);
}

inline bool operator<(const NetworkNode& lhs, const NetworkNode& rhs) {
	return lhs.key < rhs.key;
}

inline bool operator>(const NetworkNode& lhs, const NetworkNode& rhs) {
	return lhs.key > rhs.key;
}

inline bool operator<=(const NetworkNode& lhs, const NetworkNode& rhs) {
	return lhs.key <= rhs.key;
}

inline bool operator>=(const NetworkNode& lhs, const NetworkNode& rhs) {
	return lhs.key >= rhs.key;
}
/**/