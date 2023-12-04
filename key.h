#pragma once
#include <any>
#include <string>
#include <typeinfo>

class Object {
public:
	std::any value;

	template <typename T>
	Object(T value);

	std::any getValue() const;
	const type_info& getType() const;
	const char* getTypeName() const;
	bool isInt() const;
	bool isDouble() const;
	bool isChar() const;
	bool isString() const;
};

