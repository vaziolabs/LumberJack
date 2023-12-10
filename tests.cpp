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
	KeyNode parent_node = KeyNode("Parent");
	KeyNode child_1 = KeyNode("First Child", &parent_node);	
	KeyNode child_2 = KeyNode("Second Child");

	parent_node.children.push_back(&child_1);
	parent_node.children.push_back(&child_2);

	child_2.parent = &parent_node;

	std::cout << "Parent: \t\t" << parent_node.key << std::endl;
	std::cout << "\tChild[0]: \t" << parent_node.children[0]->key << std::endl;
	std::cout << "\tChild[1]: \t" << parent_node.children[1]->key << std::endl;

	std::cout << "Child 1's Parent: \t" << child_1.parent->key << std::endl;
	std::cout << "Child 2's Parent: \t" << child_2.parent->key << std::endl;
}

void keyTreeTest() {
	KeyTree* key_tree = new KeyTree(new KeyNode("Root"));

	key_tree->insert(new KeyNode("First Child"));
	key_tree->insert(new KeyNode("Second Child"));

	key_tree->insert(new KeyNode("First Grandchild"));
	key_tree->insert(new KeyNode("Second Grandchild"));

	key_tree->insert(new KeyNode("First Great Grandchild"));

	key_tree->insert(new KeyNode("Third Child"));

	key_tree->print();

	printf("\n");

	KeyNode* cursor = key_tree->root;
	std::cout << "Root: " << cursor->key << std::endl;
	std::cout << "Root's Children: " << std::endl;
	for (int i = 0; i < cursor->children.size(); i++) {
		std::cout << "\t" << cursor->children[i]->key << std::endl;
	}

	cursor = cursor->children[0];
	printf("\n");

	std::cout << "First Child: " << cursor->key << std::endl;
	std::cout << "First Child's Parent: " << cursor->parent->key << std::endl;
	std::cout << "First Child's Children: " << std::endl;
	for (int i = 0; i < cursor->children.size(); i++) {
		std::cout << "\t" << cursor->children[i]->key << std::endl;
	}

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

	//printf("\nKey tests:\n");
	//keyTest();

	//printf("\nKeyNode tests:\n");
	//keyNodeTest();

	printf("\nKeyTree tests:\n");
	keyTreeTest();


	// printf("\nValue tests:\n");
	// valueTest();


}