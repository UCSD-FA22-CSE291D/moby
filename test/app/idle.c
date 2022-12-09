#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/types.h> 
#include <sys/wait.h>
#include <time.h>

#define KB 1024
#define MB 1024 * KB
#define WORKLOAD 16 * MB
#define STABLESIZE 99 * WORKLOAD

// in second
#define HIBERNATE sleep(10000)


int main() {
    u_int8_t *data = (u_int8_t*)malloc(sizeof(u_int8_t) * STABLESIZE);
    //initialize data
    for (int i = 0; i < STABLESIZE; i++) {
        data[i] = rand() % 128;
    }
    printf("stable workload is on\n");
    HIBERNATE;
    free(data);
    return 0;
}