#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

int global = 0;
void* thread_add(){
    int i;
    for(i = 0; i<1000000; i = i+1){
        global=global+1;
    }
    return NULL;
}

void* thread_sub(){
    int j;
    for(j = 0; j<1000000; j= j+1){
    global = global-1;
    }
    return NULL;
}
int main(void)
{
    pthread_t thread1, thread2;
   
    pthread_create( &thread1, NULL, thread_add, NULL);
    pthread_create( &thread2, NULL, thread_sub, NULL);


    printf("%d", global);

    return 0;
}