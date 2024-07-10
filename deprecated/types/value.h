#ifndef VALUE_H 
#define VALUE_H
#include <stdbool.h>

typedef struct Obj Obj;
typedef struct OString OString;
//  TODO: typedef struct Numeral Numeral;

typedef enum {
    V_BOOLEAN,
    V_NUMERAL,
    V_STRING,
    V_OBJECT,
    V_PAIR,
    V_NONE,
} ValueT;

// TODO: add numeral type conversion to/from binary, hex, octal
typedef struct {
    ValueT type;
    union {
        bool boolean;
        double numeral;
        char* string;
        Obj* object;
    } as;
} Value;

typedef struct {
    int capacity;
    int count;
    Value **value;
} Values;

#define NONE_VALUE              ((Value){V_NONE, {.numeral = 0}})
#define BOOLEAN_VALUE(value)    ((Value){V_BOOLEAN, {.boolean = value}})
#define NUMERAL_VALUE(value)    ((Value){V_NUMERAL, {.numeral = value}})
#define STRING_VALUE(value)     ((Value){V_STRING, {.string = value}})
#define OBJECT_VALUE(obj)       ((Value){V_OBJECT, {.object = (Obj*)obj}})
// ADD HEXIDECIMAL, BINARY, OCTAL CONVERSIONS to NUMERAL TYPE
#define AS_BOOLEAN(value)       ((value).as.boolean)
#define AS_NUMERAL(value)       ((value).as.numeral)
#define AS_STRING(value)        ((value).as.string)
#define AS_OBJECT(value)        ((value).as.object)
#define IS_NONE(value)          ((value).type == V_NONE)
#define IS_BOOLEAN(value)       ((value).type == V_BOOLEAN)
#define IS_NUMERAL(value)       ((value).type == V_NUMERAL)
#define IS_OBJECT(value)        ((value).type == V_OBJECT)

Value *createValue(ValueT type, void *value);
void initValues (Values* arr);
void addValue (Values* arr, Value *val);
bool equalValues (Value a, Value b);
void destroyValues (Values* arr);
void printValue (Value value);

#endif
