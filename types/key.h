#ifndef TYPES_KEY_H
#define TYPES_KEY_H


typedef enum {
        TYPE_INT,
        TYPE_STRING,
} KeyType;

typedef union {
        int int_v;
        char *string_v;
} KeyValue;

typedef struct {
        KeyType type;
        KeyValue value;
} Key;

// functions
Key *createKey(KeyType type, void *value);
void printKey(Key key);
void freeKey(Key *key);

#endif