#ifndef NODE_H
#define NODE_H

#include "key.h"
#include "connection.h"
#include "value.h"

// C struct that represents a node in a network.
typedef struct {
    Connections *connections;
    Keys *keys;
    Values *values;
} Node;

typedef struct {
    int n_Nodes;
    Node **node;
} Nodes;

// Functions to create and destroy nodes.
void initNodes(Nodes *nodes);
Node *newNode();
void addNode(Nodes *nodes, Node *n);
void destroyNode(Node *n);
void destroyNodes(Nodes *n);
void printNode(Node *n);

#endif
