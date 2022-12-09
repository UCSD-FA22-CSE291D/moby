#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/types.h> 
#include <sys/wait.h>
#include <time.h>

// number of workers
#define N 128

// the size of working data in each
#define KB 1024
#define MB 1024 * KB
#define WORKLOAD 16 * MB

// sleep time of worker in nanosecond
#define INTERVAL 100 * 1000 * 1000

// in second
#define HIBERNATE sleep(100)

void worker(int id) {
    u_int8_t *data = (u_int8_t*)malloc(sizeof(u_int8_t) * WORKLOAD);
    //initialize data
    for (int i = 0; i < WORKLOAD; i++) {
        data[i] = rand() % 128;
    }
    printf("initialized worker %d\n", id);

    u_int32_t work_finished = 0;
    // working indefinitely
    while (1) {
        if (id == 0) {
            // modify data
            for (int i = 0; i < WORKLOAD; i++) {
                data[i] = data[i] + 1;
            }
            work_finished++;
            printf("update from worker %d\n", id);
            printf("finished job: %d\n", work_finished);
            time_t t;
            time(&t);
            printf("timestamp: %s", ctime(&t));

            // prevent thrashing
            usleep(INTERVAL);
        }
        else { // otherwise loop for nothing
            HIBERNATE;
        }
    }

    free(data);
    return;
}

void master() {
    pid_t pid;
    int count = 0;
    for (int i = 0; i < N; i++) {
        if((pid = fork()) < 0) {
            perror("Failed at fork\n");
            exit(1);
        }
        
        if (pid == 0) { // child
            worker(i);
        } else { // parent
            count++;
        }
    }
    while (count > 0) {
        wait(NULL);
        count--;
    }
    return;
}

int main() {
    master();
    return 0;
}
