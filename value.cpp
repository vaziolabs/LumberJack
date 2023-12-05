#include "value.h"

bool Value::isBool() const { return std::holds_alternative<bool>(this->value); }
bool Value::isInt() const { return std::holds_alternative<int>(this->value); }
bool Value::isDouble() const { return std::holds_alternative<double>(this->value); }
bool Value::isChar() const { return std::holds_alternative<char>(this->value); }
bool Value::isString() const { return std::holds_alternative<std::string>(this->value); }

const type_info& Value::getType() {
	if (this->isBool()) return typeid(bool);
	else if (this->isInt()) return typeid(int);
	else if (this->isDouble()) return typeid(double);
	else if (this->isChar()) return typeid(char);
	else if (this->isString()) return typeid(std::string);
	else return typeid(void);
}

std::string Value::getTypeName() {
	if (this->isBool()) return "bool";
	else if (this->isInt()) return "int";
	else if (this->isDouble()) return "double";
	else if (this->isChar()) return "char";
	else if (this->isString()) return "string";
	else return "void";
}

auto Value::getValue() -> decltype(value) { return this->value; }

std::ostream& operator<<(std::ostream& os, const Value& obj) {
	std::visit([&os](const auto& val) { os << val; }, obj.value);
	return os;
}

