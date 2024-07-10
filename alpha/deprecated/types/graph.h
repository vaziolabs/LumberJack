#ifndef GRAPH_H
#define GRAPH_H

#include "key.h"
#include "connection.h"
#include "node.h"

typedef struct {
    Key *key;
    Node **nodes;
    int n_nodes;
    Connection **connections;
    int n_connections;
} Graph;

typedef struct {
    int n_graphs;
    Graph **graph;
} Graphs;

Graph *newGraph();
void initGraph(Graph *graph);
void initGraphs(Graphs *graphs);
void destroyGraph(Graph *graph);
void destroyGraphs(Graphs *graphs);
void addGraph(Graphs *graphs, Graph *graph);
void addNode(Graph *graph, Node *node);
void addConnection(Graph *graph, Connection *connection);
void printGraph(Graph *graph);

#endif