#pragma once
#include "key.h"
#include "value.h"
#include "connection.h"
#include <map>
#include <set>

struct NetworkNode {
	Key key;
	std::map<Key*, std::set<Value>*> datamap;
	// TODO: figure out how to determine a pointer to some values

	NetworkNode() : key() {}
};
