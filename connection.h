#pragma once
#include "key.h"
#include "value.h"
#include "networknode.h"
#include "function.h"

typedef struct Connection {
	Key* key;
	NetworkNode *from;
	NetworkNode *to;
	//Function<Value> function;
} Connection;