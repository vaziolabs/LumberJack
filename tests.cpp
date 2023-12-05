#include "tests.h"
#include <iostream>
#include <cstdio>

static inline Value testFunction() {
	return Value("Test Function");
}

void keyTest() {
	Key int_type = 1;
	Key char_type = 'a';
	Key string_type = "Hello, Keys!";

	std::cout << int_type.getTypeName() << ": " << int_type << std::endl;
	std::cout << char_type.getTypeName() << ": " << char_type << std::endl;
	std::cout << string_type.getTypeName() << ": " << string_type << std::endl;
}

void keyNodeTest() {
	KeyNode key_node1 = KeyNode("Parent");
	KeyNode key_node2 = KeyNode("Child 1");
	KeyNode key_node3 = KeyNode("Child 2");

	key_node1.children.push_back(&key_node2);
	key_node1.children.push_back(&key_node3);

	key_node2.parent = &key_node1;
	key_node3.parent = &key_node1;

	std::cout << key_node1.key << std::endl;
	std::cout << key_node1.children[0]->key << std::endl;
	std::cout << key_node1.children[1]->key << std::endl;

	std::cout << key_node2.parent->key << std::endl;
	std::cout << key_node3.parent->key << std::endl;
}

void keyTreeTest() {
	KeyTree key_tree;

	return;
}


void valueTest() {
	Value bool_type = Value(true);
	Value int_type = Value(1);
	Value double_type = Value(1.0);
	Value char_type = Value('a');
	Value string_type = Value("Hello, Values!");

	std::cout << bool_type.getTypeName() << ": " << bool_type << std::endl;
	std::cout << int_type.getTypeName() << ": " << int_type << std::endl;
	std::cout << double_type.getTypeName() << ": " << double_type << std::endl;
	std::cout << char_type.getTypeName() << ": " << char_type << std::endl;
	std::cout << string_type.getTypeName() << ": " << string_type << std::endl;
}

void testMain() {
	printf("Running tests...\n");

	printf("\nKeyNode tests:\n");
	keyNodeTest();

	// printf("\nKey tests:\n");
	// keyTest();

	// printf("\nValue tests:\n");
	// valueTest();


}