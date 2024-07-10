#include <stdio.h>
#include "connection.h"
#include "memory.h"


void initConnections(Connections *c) {
		c->connection = NULL;
		c->count = 0;
		c->capacity = 0;

		return;
}

void initConnection(Connection* c) {
    c->key = NULL;
    c->from = NULL;
    c->to = NULL;
    c->lambda = NULL;

    return;
}

Connection *newConnection() {
        Connection *c = ALLOCATE(Connection, 1);

        initConnection(c);
        
        return c;
}

void addConnection(Connections *c, Connection *connection) {
		if (c->count + 1 > c->capacity) {
				int oldCapacity = c->capacity;
				c->capacity = EXPAND_CAPACITY(oldCapacity);
				c->connection = EXPAND(Connection, c->connection, oldCapacity, c->capacity);
		}

		c->connection[c->count] = connection;
		c->count++;

		return;
}

void destroyConnection(Connection *c) {
        freeKey(c->key);

        FREE(Connection, c);
        return;
}


void destroyConnections(Connections* c) {
    for (int i = 0; i < c->count; i++) {
        destroyConnection(c->connection[i]);
    }

    FREE_ARRAY(Connection, c->connection, c->capacity);
    initConnections(c);

    return;
}

void printConnection(Connection *c) {
        printf("Key: ");
        printKey(*c->key);
        printf("\n");
        printf("From: ");
        printNode(c->from);
        printf("\n");
        printf("To: ");
        printNode(c->to);
        printf("\n");
        printf("Lambda: TODO");
        //TODO: printLambda(c->lambda);

        return;
}

void printConnections(Connections *c) {
		for (int i = 0; i < c->count; i++) {
				printConnection(c->connection[i]);
				printf("\n");
		}

		return;
}
