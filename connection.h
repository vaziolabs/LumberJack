#pragma once
#include "key.h"
#include "value.h"
#include "node.h"
#include "function.h"

typedef struct Connection {
	Key key;
	Node* from = nullptr;
	Node* to = nullptr;
	Function<Value> function;
} Connection;