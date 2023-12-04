#include "key.h"
#include <iostream>

int main() {
	Object int_type = Object(1);
	Object double_type = Object(1.2);
	Object char_type = Object('a');
	Object string_type = Object("Hello, World!");

	std::cout << int_type.getTypeName() << ": " << int_type << std::endl;
	std::cout << double_type.getTypeName() << ": " << double_type << std::endl;
	std::cout << char_type.getTypeName() << ": " << char_type << std::endl;
	std::cout << string_type.getTypeName() << ": " << string_type << std::endl;

	return 0;
}
