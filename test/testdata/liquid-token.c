#include <string.h>
#include "vertex.h"

extern Event Mint(address to, uint64_t amount);
extern Event Transfer(address from, address to, uint64_t amount, uint64_t memo);

const char OWNER[] = "OWNER";
const char IS_PAUSE[] = "IS_PAUSE";

// sdk_storage_get
void *sdk_storage_get(const void* key, size_t key_size) {
  size_t size = chain_storage_size_get(key, key_size);
  if (size == 0) {
      return 0;
  }
  void *ret = (void *) malloc(size);
  chain_storage_get(key, key_size, ret);
  return ret;
}

int sdk_caller_is_creator() {
  address caller;
  address creator;
  chain_get_caller(caller);
  chain_get_creator(creator);
  int n = memcmp(creator, caller, ADDRESS_SIZE);
  if (n == 0) {
    return 1;
  }
  return 0;
}

int caller_is_owner() {
  address caller;
  chain_get_caller(caller);
  void *owner = sdk_storage_get(OWNER, sizeof(OWNER));
  int n = memcmp(owner, caller, ADDRESS_SIZE);
  free(owner);
  if (n == 0) {
    return 1;
  }
  return 0;
}

int set_owner(address owner) {
  if (!caller_is_owner()) {
    return -1;
  }
  chain_storage_set(OWNER, sizeof(OWNER), owner, ADDRESS_SIZE);
  return 0;
}

int pause() {
  if (!caller_is_owner()) {
    return -1;
  }
  uint8_t pause = 1;
  chain_storage_set(IS_PAUSE, sizeof(IS_PAUSE), &pause, 1);
  return 0;
}

int unpause() {
  if (!caller_is_owner()) {
    return -1;
  }
  uint8_t unpause = 0;
  chain_storage_set(IS_PAUSE, sizeof(IS_PAUSE), &unpause, 1);
  return 0;
}

int is_pausing() {
  void *data = sdk_storage_get(IS_PAUSE, sizeof(IS_PAUSE));
  uint8_t ret = *(uint8_t *)(data);
  free(data);
  return ret;
}

uint64_t get_balance(address address) {
  void *data = sdk_storage_get(address, ADDRESS_SIZE);
  uint64_t ret = *(uint64_t *)data;
  free(data);
  return ret;
}

int change_balance(address to, uint64_t amount, int sign) {
  uint64_t to_balance = get_balance(to);
  if (sign < 0) {
    if (to_balance < amount) {
      return -1;
    }
    to_balance -= amount;
  } else {
    to_balance += amount;
  }
  chain_storage_set(to, ADDRESS_SIZE, &to_balance, 8);
  return 0;
}

int set_owner_to_creator() {
  address creator;
  chain_get_creator(creator);
  chain_storage_set(OWNER, sizeof(OWNER), creator, ADDRESS_SIZE);
  return 0;
}

int mint(uint64_t amount) {
  // set up genesis owner
  if (!caller_is_owner()) {
    void *data = sdk_storage_get(OWNER, sizeof(OWNER));
    if (!data) {
      set_owner_to_creator();
    }
    free(data);
  }

  if (!caller_is_owner()) {
    return -1;
  }

  // minting
  address caller;
  chain_get_caller(caller);
  int success = change_balance(caller, amount, 1);
  if (success != -1) {
    Mint(caller, amount);
  }
  return success;
}

int transfer_with_memo(address to, uint64_t amount, uint64_t memo) {
  if (is_pausing()) {
    return -1;
  }
  address from;
  chain_get_caller(from);
  int success = change_balance(from, amount, -1);
  if (success != -1) {
    success = change_balance(to, amount, 1);
  }
  if (success != -1) {
    Transfer(from, to, amount, memo);
    return success;
  }
  return -1;
}

int transfer(address to, uint64_t amount) {
  return transfer_with_memo(to, amount, 0);
}

int init(uint64_t amount) {
  set_owner_to_creator();
  address caller;
  chain_get_caller(caller);
  int success = change_balance(caller, amount, 1);
  if (success != -1) {
    Mint(caller, amount);
  }
  return success;
}