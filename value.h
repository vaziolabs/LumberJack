#pragma once
#include "function.h"
#include <string>
#include <typeinfo>
#include <variant>
#include <sstream>
#include <functional>
#include <any>

// TODO: Are we going to use other types 
//		 (e.g. float, long, etc.)?
//		 (enum, struct, class, etc.)?
//		 (vector, list, map, etc.)?
//		 (function, pointer, etc.)?
//		 (void, null, etc.)?
//       or do we want to use any/void or a generic <T>?

class Value {
	std::variant<bool, int, double, char, std::string, Function<Value>> value;

public:
	template <typename T>
	Value(T val) : value(val) {}

	const type_info& getType();
	std::string getTypeName();
	auto getValue() -> decltype(value);

	friend std::ostream& operator<<(std::ostream& os, const Value& obj);

private:
	bool isBool() const;
	bool isInt() const;
	bool isDouble() const;
	bool isChar() const;
	bool isString() const;
};

