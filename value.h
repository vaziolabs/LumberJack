#pragma once
#include "function.h"
#include <string>
#include <typeinfo>
#include <variant>
#include <ostream>

// TODO: Are we going to use other types 
//		 (e.g. float, long, etc.)?
//		 (enum, struct, class, etc.)?
//		 (vector, list, map, etc.)?
//		 (function, pointer, etc.)?
//		 (void, null, etc.)?
//       or do we want to use any/void or a generic <T>?

// ARE We GOING TO MAKE EVERYTHING A STRING, OR ARE WE GOING TO USE A NUMBER TYPE for Hashing Str->Int?

using ValueType = std::variant<bool, int, double, char, std::string>;

class Value {
	 ValueType v;

public:
	template <typename T>
	Value(T val) : v(val) {}

	const type_info& type_info();
	std::string type();
	auto value() -> decltype(v);

	friend std::ostream& operator<<(std::ostream& os, const Value& obj);

private:
	bool isBool() const;
	bool isInt() const;
	bool isDouble() const;
	bool isChar() const;
	bool isString() const;
};

