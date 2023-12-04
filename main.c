#include "./types/network.h"
#include "./types/graph.h"
#include "./types/node.h"
#include "./types/value.h"

int main(int argc, char **argv) {
    Network *network = createNetwork();

    Graph *graph = createGraph();

    Node *node1 = newNode();
    Node *node2 = newNode();
    Connection *connection1 = newConnection();
    Connection *connection2 = newConnection();

    node1->key[0] = createKey(TYPE_STRING, "name");
    node1->values[0] = createValue(TYPE_STRING, "John");
    
    node1->key[1] = createKey(TYPE_STRING, "age");
    node1->values[1] = createValue(TYPE_INT, 20);
    
    node2->key[0] = createKey(TYPE_STRING, "name");
	node2->key[1] = createKey(TYPE_STRING, "age");

    connection1->key = createKey(TYPE_STRING, "name");
    connection1->from = node1;
    connection1->to = node2;

    connection2->key = createKey(TYPE_STRING, "age");
    connection2->from = node2;
    connection2->to = node1;

    addNode(graph, node1);
	addNode(graph, node2);

    addConnection(graph, connection1);
    addConnection(graph, connection2);

    addGraph(network, graph);

    printNetwork(network);

    freeNetwork(network);
    return 0;
}
