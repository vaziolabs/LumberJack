#include <stdio.h>
#include "graph.h"
#include "./compiler/memory.h"

Graph *newGraph() {
        Graph *g = ALLOCATE(Graph, 1);
        g->nodes = NULL;
        g->n_nodes = 0;
        g->connections = NULL;
        g->n_connections = 0;
        return g;
}

void freeGraph(Graph *g) {
        freeNodes(g->nodes, g->n_nodes);
        freeConnections(g->connections, g->n_connections);

        FREE(Graph, g);
        return;
}

void freeGraphs(Graph **g, int n_Graphs) {
        for (int i = 0; i < n_Graphs; i++) {
                freeGraph(g[i]);
        }
        FREE_ARRAY(Graph *, g, n_Graphs);
        return;
}

void addNode(Graph *g, Node *n) {
        g->nodes = EXPAND(Node *, g->nodes, g->n_nodes, g->n_nodes + 1);
        g->nodes[g->n_nodes] = n;
        g->n_nodes++;
        return;
}

void addConnection(Graph *g, Connection *c) {
        g->connections = EXPAND(Connection *, g->connections, g->n_connections, g->n_connections + 1);
        g->connections[g->n_connections] = c;
        g->n_connections++;
        return;
}

void printGraph(Graph *g) {
        printf("Graph:\n");
        
        for (int i = 0; i < g->n_nodes; i++) {
                printf("Node %d:\n", i);
                printNode(g->nodes[i]);
        }

        for (int i = 0; i < g->n_connections; i++) {
                printf("Connection %d:\n", i);
                printConnection(g->connections[i]);
        }
        return;
}
