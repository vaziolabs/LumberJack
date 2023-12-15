#pragma once
#include "key.h"
#include "connection.h"
#include "value.h"
#include <vector>
#include <set>

// Optional type that can be either Connection or MultiConnection
class NetworkNode {
public:
	Key key;									// Key for a node i.e. "Name"
	std::vector<Value> values;					// Values for a node i.e. "John"
	std::set<Connector*> connections;			// Children of a node i.e. "John" -> "Smith"

	NetworkNode();

	template <typename T>
	NetworkNode(T key) : key(key) {}

	std::string getKeyType() const;


	template <typename T, typename U>
	NetworkNode(T key, U value) : key(key) {
		values.push_back(Value(value));
	}

	void addValue(ValueType value);
	void deleteValue(int index);
	int getIndex(Value value);
	int getIndex(ValueType value);
	void addConnection(Connector* node);
	void removeConnection(Connector* node);

	void print();
};
/**/