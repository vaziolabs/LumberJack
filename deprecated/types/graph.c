#include <stdio.h>
#include "graph.h"
#include "memory.h"

void initGraphs(Graphs *gs) {
		gs->graph = NULL;
		gs->n_graphs = 0;
		return;
}

void initGraph(Graph *g) {
		g->key = NULL;
		g->nodes = NULL;
		g->n_nodes = 0;
		g->connections = NULL;
		g->n_connections = 0;
		return;
}

Graph *newGraph() {
        Graph *g = ALLOCATE(Graph, 1);
        initGraph(g);
        return g;
}

void destroyGraph(Graph *g) {
        freeNodes(g->nodes, g->n_nodes);
        freeConnections(g->connections, g->n_connections);

        FREE(Graph, g);
        return;
}

void destroyGraphs(Graphs *gs) {
		for (int i = 0; i < gs->n_graphs; i++) {
				destroyGraph(gs->graph[i]);
		}

		FREE(Graphs, gs);
		return;
}

void addGraph(Graphs *gs, Graph *g) {
		gs->graph = EXPAND(Graph, gs->graph, gs->n_graphs, gs->n_graphs + 1);
		gs->graph[gs->n_graphs] = g;
		gs->n_graphs++;
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
