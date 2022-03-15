#include "vertex.h"

extern size_t chain_storage_size_get(byte_t*, size_t);
extern byte_t* chain_storage_get(byte_t*, size_t, byte_t*);
extern void chain_storage_set(byte_t*, size_t, byte_t*, size_t);
extern void chain_method_bind(address, byte_t*, size_t, byte_t*, size_t);
extern int chain_arg_size_get(void*);
extern int chain_arg_size_set(void*, size_t);
extern int chain_arg_size_set(void*, size_t);

extern int sum_of_squares(int*);
extern int get_mean(int*);
extern double sqroot(int);
extern int get_average(int*);
extern int address_xor(address);
extern int matched_parity(int);


extern Event address_checked(address contract, uint8_t checksum);

const char MATH_KEY[] = "math";

const char MATH_MEAN[] = "mean";
const char LOCAL_MEAN[] = "get_mean";

const char MATH_SS[] = "sum_of_squares";
const char LOCAL_SS[] = "sum_of_squares";
const char MATH_SR[] = "square_root";
const char LOCAL_SR[] = "sqroot";
const char MATH_XOR[] = "address_xor";
const char LOCAL_XOR[] = "address_xor";
const char MATH_PARITY[] = "matched_parity";
const char LOCAL_PARITY[] = "matched_parity";


// sdk_storage_set
void sdk_storage_set(byte_t* key, size_t key_size, byte_t* value, size_t value_size) {
    chain_storage_set(key, key_size, value, value_size);
}

// sdk_storage_get
byte_t* sdk_storage_get(byte_t* key, size_t key_size) {
  int size = chain_storage_size_get(key, key_size);
  if (size == 0) {
      return 0;
  }
  byte_t* ret = (byte_t*) malloc(size * sizeof(byte_t));
  chain_storage_get(key, key_size, ret);
  return ret;
}

int sdk_caller_is_creator() {
  address caller = (address) malloc(ADDRESS_SIZE * sizeof(byte_t));
  address creator = (address) malloc(ADDRESS_SIZE * sizeof(byte_t));
  chain_get_caller(caller);
  chain_get_creator(creator);
  int n = memcmp(creator, caller, ADDRESS_SIZE);
  free(caller);
  free(creator);
  if (n == 0) {
    return 1;
  }
  return 0;
}

int init(address math_contract) {
  if (sdk_caller_is_creator()) {  
    sdk_storage_set(MATH_KEY, sizeof(MATH_KEY), math_contract, ADDRESS_SIZE);
  }
}

int variance(int32_t* a) {
  int len = chain_arg_size_get(a)/4;

  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_MEAN, sizeof(MATH_MEAN), LOCAL_MEAN, sizeof(LOCAL_MEAN));
  chain_method_bind(math_contract, MATH_SS, sizeof(MATH_SS), LOCAL_SS, sizeof(LOCAL_SS));
  chain_method_bind(math_contract, MATH_SR, sizeof(MATH_SR), LOCAL_SR, sizeof(LOCAL_SR));

  int b[len];  
  int mean = get_mean(a);
  for(int i=0; i< len; i++){
    b[i] = a[i] - mean;
  }
  chain_arg_size_set(b, len * 4);
  int squared_sum = sum_of_squares(b);
  return squared_sum / len;
}

double hypotenuse(int32_t a, int32_t b) {
    address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
    chain_method_bind(math_contract, MATH_SS, sizeof(MATH_SS), LOCAL_SS, sizeof(LOCAL_SS));
    chain_method_bind(math_contract, MATH_SR, sizeof(MATH_SR), LOCAL_SR, sizeof(LOCAL_SR));

    int sides[] = {a, b};
    chain_arg_size_set(sides, 2 * 4);
    int squared_sum = sum_of_squares(sides);
    return sqroot(squared_sum);
}

uint8_t xor_checksum(address a){
  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_XOR, sizeof(MATH_XOR), LOCAL_XOR, sizeof(LOCAL_XOR));
  int ret = address_xor(a);
  address_checked(a, ret);
  return ret;
}

uint8_t mod_invoke(address a){
  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_XOR, sizeof(MATH_XOR), LOCAL_XOR, sizeof(LOCAL_XOR));
  a[0] = 0;
  int ret = address_xor(a);
  address_checked(a, ret);
  return ret;
}

uint8_t mod_emit(address a){
  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_XOR, sizeof(MATH_XOR), LOCAL_XOR, sizeof(LOCAL_XOR));
  int ret = address_xor(a);
  a[0] = 0;
  address_checked(a, ret);
  return ret;
}


// overflow test
int mean(int32_t* a){
  int len = chain_arg_size_get(a)/4;
  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_MEAN, sizeof(MATH_MEAN), LOCAL_MEAN, sizeof(LOCAL_MEAN));
  return get_mean(a);
}

// unknown import test
int average(int32_t* a) {
  return get_average(a);
}

// invalid params
int parity(int a, int b) {
  address math_contract = sdk_storage_get(MATH_KEY, sizeof(MATH_KEY));
  chain_method_bind(math_contract, MATH_PARITY, sizeof(MATH_PARITY), LOCAL_PARITY, sizeof(LOCAL_PARITY));
  return matched_parity(a);
}