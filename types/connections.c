#include <stdio.h>
#include "connection.h"
#include "./compiler/memory.h"

Connection *newConnection() {
        Connection *c = ALLOCATE(Connection, 1);

        c->key = NULL;
        c->from = NULL;
        c->to = NULL;
        c->lambda = NULL;
        
        return c;
}

void freeConnection(Connection *c) {
        freeKey(c->key);

        FREE(Connection, c);
        return;
}

void freeConnections(Connection **c, int n) {
        for (int i = 0; i < n; i++) {
                freeConnection(c[i]);
        }

        FREE_ARRAY(Connection, c, n);
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
