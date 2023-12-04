#include "node.h"
#include "./compiler/memory.h"

Node *newNode() {
        Node *node = ALLOCATE(Node, 1);

        node->connections = NULL;
        node->n_Connections = 0;
        node->key[0] = createKey(TYPE_STRING, "name");
        node->n_Keys = 1;
        initValues(node->values);

        return node;
}

void freeNode(Node *n) {
        freeValues(n->values);
        freeConnections(n->connections, n->n_Connections);
        freeKeys(n->key, n->n_Keys);
        FREE(n);
}

void freeNodes(Node **n, int n_Nodes) {
        for (int i = 0; i < n_Nodes; i++) {
                freeNode(n[i]);
        }
        FREE_ARRAY(n, n_Nodes);
}

void printNode(Node *n) {
        printf("Keys:\n");
        for (int i = 0; i < n->n_Keys; i++) {
                printf("Key %d: ", i);
                printKey(n->key[i]);
                printf("\n");
        }
        printf("Connections:\n");
        for (int i = 0; i < n->n_Connections; i++) {
                printf("Connection %d:\n", i);
                printConnection(n->connections[i]);
        }
        printf("Values:\n");
        printValues(n->values);
}
