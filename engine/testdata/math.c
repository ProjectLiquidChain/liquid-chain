#include <math.h>
#include "vertex.h"

extern int chain_arg_size_get(void*);

const int ADDR_SIZE=35;

int mean(int32_t* a){
   int len = chain_arg_size_get(a)/4;
   int sum = 0;
   for(int i=0; i<len; i++) {
      sum += a[i];
   }
   return sum / len;
}

int sum_of_squares(int32_t* a) {
   int len = chain_arg_size_get(a)/4;
   int sum = 0;
   for(int i=0; i<len; i++) {
      sum += a[i]*a[i];
   }
   return sum;
}

double square_root(int32_t a) {
    return sqrt(a);
}

int address_xor(address a){
   int ret = 0;
   for (int i = 0; i < ADDR_SIZE; i++){
      ret = ret ^ a[i];
   }
   return ret;
}

int matched_parity(int a, int b) {
  return a%2 == b%2;
}