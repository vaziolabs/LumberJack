#include "value.h"

bool Value::isBool() const { return std::holds_alternative<bool>(this->v); }
bool Value::isInt() const { return std::holds_alternative<int>(this->v); }
bool Value::isDouble() const { return std::holds_alternative<double>(this->v); }
bool Value::isChar() const { return std::holds_alternative<char>(this->v); }
bool Value::isString() const { return std::holds_alternative<std::string>(this->v); }

const type_info& Value::type_info() {
	if (this->isBool()) return typeid(bool);
	else if (this->isInt()) return typeid(int);
	else if (this->isDouble()) return typeid(double);
	else if (this->isChar()) return typeid(char);
	else if (this->isString()) return typeid(std::string);
	else return typeid(void);
}

std::string Value::type() {
	if (this->isBool()) return "bool";
	else if (this->isInt()) return "int";
	else if (this->isDouble()) return "double";
	else if (this->isChar()) return "char";
	else if (this->isString()) return "string";
	else return "void";
}

auto Value::value() -> decltype(v) { return this->v; }

std::ostream& operator<<(std::ostream& os, const Value& obj) {
	std::visit([&os](const auto& val) { os << val; }, obj.v);
	return os;
}

