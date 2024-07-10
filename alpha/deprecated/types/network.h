// A network is comprised of multiple graphs, each of which has a scope/relationship to the others in the network.

#ifndef NETWORK_H
#define NETWORK_H

#include "graph.h"

// C struct that represents a network of graphs.
typedef struct {
    Graphs *graphs;
    int n_graphs;
} Network;

typedef struct {
	Network **network;
	int n_networks;
} Networks;

// Functions to create and destroy networks.
Network *createNetwork();
void initNetwork(Network *n);
void initNetworks(Networks *n);
void addNetwork(Networks *n, Network *net);
void destroyNetwork(Network *n);
void destroyNetworks(Networks *n);
void printNetwork(Network *n);
void printNetworks(Networks *n);

#endif
