// A network is comprised of multiple graphs, each of which has a scope/relationship to the others in the network.

#ifndef NETWORK_H
#define NETWORK_H

#include "graph.h"

// C struct that represents a network of graphs.
typedef struct {
    Graph **graphs;
    int n;
} Network;

// Functions to create and destroy networks.
Network *createNetwork();
void addGraph(Network *n, Graph *g);
void freeNetwork(Network *n);
void freeNetworks(Network **n, int n_n);
void printNetwork(Network *n);

#endif
