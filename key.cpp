#include "key.h"

template <typename T>
Object::Object(T value) { this->value = value; }



const type_info& Object::getType() const { return value.type(); }
const char* Object::getTypeName() const { return value.type().name(); }
bool Object::isInt() const { return value.type() == typeid(int); }
bool Object::isDouble() const { return value.type() == typeid(double); }
bool Object::isChar() const { return value.type() == typeid(char); }
bool Object::isString() const { return value.type() == typeid(std::string); }

std::any Object::getValue() const { return value; }

