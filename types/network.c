#include "network.h"
#include "./compiler/memory.h"
#include <cstddef>

Network *createNetwork() {
    Network *network = ALLOCATE(Network, 1);
    network->graphs = NULL;
    network->n = 0;
    return network;
}

void freeNetwork(Network *network) {
    freeGraphs(network->graphs, network->n);
    FREE(Network, network);
}

void addGraph(Network *network, Graph *graph) {
    network->graphs = EXPAND(Graph *, network->graphs, network->n, network->n + 1);
    network->graphs[network->n] = graph;
    network->n++;
}

void printNetwork(Network *network) {
    printf("Network:\n");
    for (int i = 0; i < network->n; i++) {
        printf("Graph %d:\n", i);
        printGraph(network->graphs[i]);
    }
}

void freeNetworks(Network **networks, int n_Networks) {
    for (int i = 0; i < n_Networks; i++) {
        freeNetwork(networks[i]);
    }

    FREE_ARRAY(Network *, networks, n_Networks);
}
