#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/types.h> 
#include <sys/wait.h>
#include <time.h>

#define KB 1024
#define MB 1024 * KB
#define WORKLOAD 16 * MB

// sleep time of worker
#define INTERVAL 100 * 1000 * 1000

int main() {
    u_int8_t *data = (u_int8_t*)malloc(sizeof(u_int8_t) * WORKLOAD);
    //initialize data
    for (int i = 0; i < WORKLOAD; i++) {
        data[i] = rand() % 128;
    }
    printf("initialized worker\n");

    int work_finished = 0;
    while (1) {
        // modify data
        for (int i = 0; i < WORKLOAD; i++) {
            data[i] = data[i] + 1;
        }
        work_finished++;

        printf("update from worker\n");
        printf("finished job: %d\n", work_finished);
        time_t t;
        time(&t);
        printf("timestamp: %s", ctime(&t));

        // prevent thrashing
        usleep(INTERVAL);
    }

    free(data);
    return 0;
}