#include "key.h"

bool Key::isInt() const { return std::holds_alternative<int>(key); }
bool Key::isChar() const { return std::holds_alternative<char>(key); }
bool Key::isString() const { return std::holds_alternative<std::string>(key); }

const type_info& Key::type_info() const {
	if (this->isString()) return typeid(std::string);
	else if (this->isInt()) return typeid(int);
	else if (this->isChar()) return typeid(char);
	else return typeid(void);
}

std::string Key::type() const {
	if (this->isString()) return "string";
	else if (this->isInt()) return "int";
	else if (this->isChar()) return "char";
	else return "void";
}

auto Key::value() const -> decltype(key) { return key; }