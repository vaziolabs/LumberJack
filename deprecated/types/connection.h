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

typedef struct {
    int capacity;
    int count;
    Connection **connection;
} Connections;

// Functions to create and destroy connections.
Connection *newConnection();
void initConnection(Connection *c);
void initConnections(Connections *c);
void addConnection(Connections *c, Connection *connection);
void destroyConnection(Connection *c);
void destroyConnections(Connections *c);
void printConnection(Connection *c);

#endif
