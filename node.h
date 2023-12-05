#pragma once
#include "key.h"
#include "value.h"
#include "setmap.h"

struct Node {
	Key key;
	SetMap values;
};
