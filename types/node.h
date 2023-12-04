#ifndef NODE_H
#define NODE_H

#include "key.h"
#include "connection.h"
#include "value.h"

// C struct that represents a node in a network.
typedef struct {
    Connection **connections;
    int n_Connections;
    Key **key;
    int n_Keys;
    Values **values;
    int n_Values;
} Node;

// Functions to create and destroy nodes.
Node *newNode();
void freeNode(Node *n);
void freeNodes(Node **n, int n_Nodes);
void printNode(Node *n);

#endif
