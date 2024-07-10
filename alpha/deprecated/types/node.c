#include "node.h"
#include "memory.h"


static void initNode(Node *n) {
        initKeys(n->keys);
		initConnections(n->connections);
		initValues(n->values);
}

void initNodes(Nodes *n) {
		n->n_Nodes = 0;
		n->node = NULL;
}

Node *newNode() {
        Node *node = ALLOCATE(Node, 1);

        initNode(node);

        return node;
}

void addNode(Nodes *nodes, Node *node) {
        nodes->node = EXPAND(Node, nodes->node, nodes->n_Nodes, 1);

        nodes->node[nodes->n_Nodes] = node;
        nodes->n_Nodes++;
}

void destroyNode(Node *n) {
        destroyValues(n->values);
        destroyConnections(n->connections);
        destroyKeys(n->keys);

        FREE(Node, n);
}

void destroyNodes(Nodes *n) {
        for (int i = 0; i < n->n_Nodes; i++) {
            destroyNode(n->node[i]);
        }

        FREE_ARRAY(Node, n, n->n_Nodes);
}

void printNode(Node *n) {
        printf("Node:\n");
        printf("Keys: %d\n", n->keys->count);
        printf("Connections: %d\n", n->connections->count);
        printf("Values: %d\n", n->values->count);
        printf("Values:");
        
        for (int i = 0; i < n->values->count; i++) {
            printValue(*n->values->value[i]);
        }
}
