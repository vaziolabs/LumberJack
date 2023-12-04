#include "key.h"

bool Object::isInt() const { return std::holds_alternative<int>(this->value); }
bool Object::isDouble() const { return std::holds_alternative<double>(this->value); }
bool Object::isChar() const { return std::holds_alternative<char>(this->value); }
bool Object::isString() const { return std::holds_alternative<std::string>(this->value); }

const type_info& Object::getType() {
	if (this->isInt()) return typeid(int);
	else if (this->isDouble()) return typeid(double);
	else if (this->isChar()) return typeid(char);
	else if (this->isString()) return typeid(std::string);
	else return typeid(void);
}

std::string Object::getTypeName() {
	if (this->isInt()) return "int";
	else if (this->isDouble()) return "double";
	else if (this->isChar()) return "char";
	else if (this->isString()) return "string";
	else return "void";
}

auto Object::getValue() -> decltype(value) { return this->value; }