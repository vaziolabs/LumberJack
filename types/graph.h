
#include "key.h"
#include "value.h"
#include "connection.h"
#include "node.h"

typedef struct {
    Node **nodes;
    int n_nodes;
    Connection **connections;
    int n_connections;
} Graph;

Graph *createGraph();
void destroyGraph(Graph *graph);
void addNode(Graph *graph, Node *node);
void addConnection(Graph *graph, Connection *connection);
void printGraph(Graph *graph);
void printGraphNodes(Graph *graph);
void printGraphConnections(Graph *graph);
