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

typedef struct {
        int capacity;
        int count;
        Key **key;
} Keys;

// functions
void initKeys(Keys *keys);
void addKey(Keys *keys, Key *key);
Key *createKey(KeyType type, void *value);
void printKey(Key key);
void destroyKey(Key *key);
void destroyKeys(Keys *keys);

#endif