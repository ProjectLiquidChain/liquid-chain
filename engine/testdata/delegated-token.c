#include <string.h>
#include <vertex.h>

extern void chain_args_write(byte_t*, byte_t*, uint32_t);
extern void chain_args_hash(byte_t*, byte_t*);
extern int chain_ed25519_verify(address, byte_t*, byte_t*);
extern int chain_get_contract_address(address);


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

int from_bytes(uint8_t* bytes){
  if (bytes==0) {
    return 0;
  }
  return *(uint32_t*)bytes;
}

int change_balance(address to, uint64_t amount, int sign){
  uint64_t to_balance = from_bytes(sdk_storage_get(to, ADDRESS_SIZE));
  if (sign < 0) {
    if (to_balance < amount) {
      return -1;
    }
    to_balance -= amount;
  } else {
    to_balance += amount;
  }
  sdk_storage_set(to, ADDRESS_SIZE, (uint8_t*)&to_balance, 8);
  return 0;
}

//////////PUBLIC

int mint(uint64_t amount) {
  byte_t caller[ADDRESS_SIZE];
  chain_get_caller(caller);
  int success = change_balance(caller, amount, 1);
  return success;
}

int get_balance(address address){
  return from_bytes(sdk_storage_get(address, ADDRESS_SIZE));
}

int transfer(address to, uint64_t amount){
  byte_t from[ADDRESS_SIZE];
  chain_get_caller(from);
  int success = change_balance(from, amount, -1);
  if (success != -1) {
    change_balance(to, amount, 1);
  }
  return success;
}

int nonce_exists(address caller, uint32_t nonce) {
    // TODO nonce management
    return 0;
}

int delegated_transfer(address to, uint64_t amount, address caller, uint32_t nonce, uint8_t* signature) {
    if (nonce_exists(caller, nonce)) {
        return -1;
    }
    address contract;
    chain_get_contract_address(contract);
    byte_t writer[4*4];
    chain_args_write(writer, to, ADDRESS_SIZE);
    chain_args_write(writer, contract, ADDRESS_SIZE);
    chain_args_write(writer, &amount, sizeof(uint64_t));
    chain_args_write(writer, &nonce, sizeof(uint32_t));
    byte_t hasher[32];
    chain_args_hash(writer, hasher);

    if (!chain_ed25519_verify(caller, hasher, signature)) {
        return -1;
    }
    int success = change_balance(caller, amount, -1);
    if (success != -1) {
        change_balance(to, amount, 1);
    }
    return success;
}