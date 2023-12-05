#include "tests.h"
#include <iostream>
#include <cstdio>

inline Value test() {
	return Value("Test Function");
}

void connectionTest() {
	Node node1;
	node1.key = Key(1);
	node1.values.addKV(Key("Name"), Value("Node 1"));

	Node node2;
	node2.key = Key(2);
	node2.values.addKV(Key("Name"), Value("Node 2"));

	Connection connection;

	connection.key = Key(1);
	connection.from = &node1;
	connection.to = &node2;
	connection.function = test;

	std::cout << connection.key << std::endl;
	std::cout << "Connection From: " << connection.from->key << std::endl;
	std::cout << "Connection To: " << connection.to->key << std::endl;
	std::cout << connection.function() << std::endl;
}

void keyTest() {
	Key int_type = Key(1);
	Key char_type = Key('a');
	Key string_type = Key("Hello, Keys!");

	std::cout << int_type.getTypeName() << ": " << int_type << std::endl;
	std::cout << char_type.getTypeName() << ": " << char_type << std::endl;
	std::cout << string_type.getTypeName() << ": " << string_type << std::endl;
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

	printf("\nConnection tests:\n");
	connectionTest();
	// printf("\nKey tests:\n");
	// keyTest();

	// printf("\nValue tests:\n");
	// valueTest();


}