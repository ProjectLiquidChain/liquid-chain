#include <math.h>
#include <vertex.h>

extern uint32_t chain_block_height();
extern uint64_t chain_block_time();

uint32_t block_height() {
   return chain_block_height();
}

uint64_t block_time() {
   return chain_block_time();
}