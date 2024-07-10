#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include "memory.h"
#include "value.h"
#include "object.h"

void initValues (Values* arr) {
    arr->value = NULL;
    arr->count = 0;
    arr->capacity = 0;

    return;
}

void destroyValues (Values* arr) {
    FREE_ARRAY(Value, arr->value, arr->count);
    initValues(arr);

    return;
}

Value *createValue(ValueT type, void *value) {
    Value *val = ALLOCATE(Value, 1);
    val->type = type;
    switch (type) {
        case V_BOOLEAN:
            val->as.boolean = *((bool *) value);
            break;
        case V_STRING:
			val->as.string = (char *) value;
			break;
        case V_NUMERAL:
            val->as.numeral = *((double *) value);
            break;
        case V_OBJECT:
            val->as.object = (Obj *) value;
            break;
    }
    return val;
}

// Todo, Compare 3 Values
bool equalValues (Value a, Value b) {
    if (a.type != b.type) { return false; }
    switch (a.type) {
        case V_BOOLEAN:     return AS_BOOLEAN(a) == AS_BOOLEAN(b);
        case V_NONE:        return true;
        case V_NUMERAL:     return AS_NUMERAL(a) == AS_NUMERAL(b);
        case V_STRING:      return strcmp(AS_STRING(a), AS_STRING(b)) == 0;
        case V_OBJECT:      return AS_OBJECT(a) == AS_OBJECT(b);
        default:
            return false;
    }
}

// TODO: Move this to Memory ADD_ITEM(t, c, i, v) as a macro definition
void addValue (Values* arr, Value *val) {
    if (arr->capacity < arr->count + 1) {
        int current_limit = arr->capacity;
        arr->capacity = EXPAND_CAPACITY(current_limit);
        arr->value = EXPAND(Value, arr->value, current_limit, arr->capacity);
    }

    arr->value[arr->count] = val;
    arr->count++;

    return;
}

void printValue (Value val) {
    switch(val.type) {
        case V_BOOLEAN:
            printf(AS_BOOLEAN(val) ? "true" : "false");
            break;
        case V_NONE:
            printf("none");
            break;
        case V_NUMERAL:
            printf("%g", AS_NUMERAL(val));
            break;
        case V_OBJECT:
            printObject(val);
            break;
    }
}
