#include "setmap.h"

void SetMap::addKV(Key key, Value value) {
	// Create a new set of key
	KeySet keys;
	keys.insert(key);

	// Create a new vector of value
	ValueVector values;
	values.push_back(value);

	// Add the new set and vector to the map
	this->addSet(keys, values);
}

void SetMap::insertKey(int index, Key key) {
	// Get old set
	KeySet keys = this->getKeySetFromIndex(index);
	ValueVector values = this->map[keys];

	// Remove old set
	map.erase(keys);

	// Insert new key
	keys.insert(key);
	this->addSet(keys, values);
}

void SetMap::addValue(int index, Value value) {
	auto it = map.begin();

	std::advance(it, index);

	it->second.push_back(value);
}

void SetMap::addValueToSet(KeySet keys, Value value) { map[keys].push_back(value); }

void SetMap::distributedValue(Key key, Value value) {
	std::list<int> indices = this->indicesOfKey(key);

	for (auto& index : indices) {
		this->addValue(index, value);
	}
}

void SetMap::addSet(KeySet keys, ValueVector values) {
	map.insert(std::pair<KeySet, ValueVector>(keys, values));
}


KeyVector SetMap::getKeys() {
		KeyVector keys;

		for (auto& [k, v] : map) {
			for (auto& key : k) {
				keys.push_back(key);
			}
		}

		return keys;
}

std::list<int> SetMap::indicesOfKey(Key key) {
	std::list<int> indices;

	int i = 0;
	for (auto& [keys, vals] : map) {
		if (keys.find(key) != keys.end()) {
			indices.push_back(i);
		}
		i++;
	}

	return indices;
}

ValueVector SetMap::getValues() {
	ValueVector values;

	for (auto& [k, v] : map) {
		for (auto& value : v) {
			values.push_back(value);
		}
	}

	return values;
}

ValueVector SetMap::getValuesAtIndex(int index) {
	auto it = map.begin();

	std::advance(it, index);

	return it->second;
}

KeyVector SetMap::getKeyVectorFromIndex(int index) {
	KeyVector keys;
	auto it = map.begin();

	std::advance(it, index);

	for (auto& key : it->first) {
		keys.push_back(key);
	}

	return keys;
}

KeySet SetMap::getKeySetFromIndex(int index) {
	KeySet keys;
	auto it = map.begin();

	std::advance(it, index);

	return it->first;
}

ValueVector SetMap::getValuesFromSet(KeySet keys) { return map[keys]; }

int SetMap::countSetsWithKey(Key key) {
	int count = 0;

	for (auto& [keys, vals] : map) {
		if (keys.find(key) != keys.end()) {
			count++;
		}
	}

	return count;
}

int SetMap::countMatchingSets(KeySet keys) {
	int count = 0;

	for (auto& [k, v] : map) {
		if (k == keys) {
			count++;
		}
	}

	return count;
}

std::list<std::vector<Value>> SetMap::listOfVectorsFromKey(Key key) {
	std::list<std::vector<Value>> values;

	for (auto& [keys, vals] : map) {
		if (keys.find(key) != keys.end()) {
			values.push_back(vals);
		}
	}

	return values;
}
