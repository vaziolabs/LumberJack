#include "network.h"
#include "memory.h"
#include <stdio.h>

void initNetworks(Networks *networks) {
	networks->network = NULL;
	networks->n_networks = 0;

	return;
}

static void initNetwork(Network* network) {
    network->graphs = createGraphs();
    network->n_graphs = 0;

    return;
}

Network *createNetwork() {
    Network *network = ALLOCATE(Network, 1);
    
    initNetwork(network);
    
    return network;
}

void addNetwork(Networks *networks, Network *network) {
	networks->network = EXPAND(Network, networks->network, networks->n_networks, networks->n_networks + 1);

    networks->network[networks->n_networks] = network;
	networks->n_networks++;
}

void printNetwork(Network *network) {
    printf("Network:\n");
    for (int i = 0; i < network->graphs->n_graphs; i++) {
        printf("Graph %d:\n", i);
        printGraph(network->graphs->graph[i]);
    }
}


void destroyNetwork(Network* network) {
    destroyGraph(network->graphs);
    FREE(Network, network);
}

void destroyNetworks(Networks *networks) {
    for (int i = 0; i < networks->n_networks; i++) {
        freeNetwork(networks[i]);
    }

    FREE_ARRAY(Network *, networks, networks->n_networks);
}

void printNetworks(Networks *networks) {
	for (int i = 0; i < networks->n_networks; i++) {
		printf("Network %d:\n", i);
		printNetwork(networks->network[i]);
	}
}
