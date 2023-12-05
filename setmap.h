#pragma once
#include "key.h"
#include "value.h"
#include <set>
#include <map>
#include <vector>

using KeySet = std::set<Key, KeyComparator>;
using ValueVector = std::vector<Value>;
using KeyVector = std::vector<Key>;

class SetMap {
private:
	std::map<KeySet, ValueVector> map;

public:
	SetMap() {}

	void addKV(Key key, Value value);
	void insertKey(int index, Key key);
	void addValue(int index, Value value);
	void addValueToSet(KeySet keys, Value value);
	void distributedValue(Key key, Value value);
	void addSet(KeySet keys, ValueVector values);

	KeyVector getKeys();
	ValueVector getValues();
	ValueVector getValuesAtIndex(int index);
	KeyVector getKeyVectorFromIndex(int index);
	KeySet getKeySetFromIndex(int index);
	ValueVector getValuesFromSet(KeySet keys);
	int countSetsWithKey(Key key);
	int countMatchingSets(KeySet keys);
	std::list<int> indicesOfKey(Key key);
	std::list<ValueVector> listOfVectorsFromKey(Key key);
};