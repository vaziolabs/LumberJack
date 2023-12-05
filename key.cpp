#include "key.h"

bool Key::isInt() const { return std::holds_alternative<int>(this->value); }
bool Key::isChar() const { return std::holds_alternative<char>(this->value); }
bool Key::isString() const { return std::holds_alternative<std::string>(this->value); }

const type_info& Key::getType() {
	if (this->isInt()) return typeid(int);
	else if (this->isChar()) return typeid(char);
	else if (this->isString()) return typeid(std::string);
	else return typeid(void);
}

std::string Key::getTypeName() {
	if (this->isInt()) return "int";
	else if (this->isChar()) return "char";
	else if (this->isString()) return "string";
	else return "void";
}

auto Key::getValue() -> decltype(value) { return this->value; }

