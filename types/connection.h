// Contents: C struct that represents a unidirectional connection between two nodes
#ifndef CONNECTION_H
#define CONNECTION_H

#include "key.h"
#include "node.h"

// C struct that represents a unidirectional connection between two nodes
// and contains the ability to introduce an anonymous lambda function.
typedef struct {
    Key *key;
    Node *from;
    Node *to;
    void (*lambda)(void);
} Connection;

// Functions to create and destroy connections.
Connection *newConnection();
void freeConnection(Connection *c);
void freeConnections(Connection c[], int n);
void printConnection(Connection *c);

#endif
