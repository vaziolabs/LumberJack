#include <cstdio>
#include <iostream>
#include "key.h"

int main() {
	Object int_type = Object(1);
	Object double_type = Object(1.0);
	Object char_type = Object('a');
	Object string_type = Object("Hello, World!");

	std::cout << int_type.getTypeName() << ": " << int_type.getValue() << std::endl;



	return 0;
}
