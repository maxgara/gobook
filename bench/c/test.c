#include "stdio.h"
#include "time.h"
//O3 got 1 second for a billion iseven calls
//O0 got 6.7 seconds
#define COUNT 1000000
#define LOOPS 1000
int iseven (int x);
long long countevens(int x);

int main(){
    clock_t start, end;
    start = clock();
    int counts[LOOPS];
    for(int i=0;i<LOOPS;i++){
    counts[i] += countevens(i)/COUNT;
    }
    long long cc = 0;
    for(int i=0;i<LOOPS;i++){
    cc += counts[i] - LOOPS/2;
    }
    end = clock();
    printf("%lld\n",cc);
    printf("even counter(%d)\n",COUNT);
    printf("time elapsed:%lu\n", end-start);

}

long long countevens(int x){
    int evens[COUNT];
    for(int i=0;i<COUNT;i++){
        evens[i] = iseven(x+i);
    }
    long long counter = 0;
    for(int i=0;i<COUNT;i++){
        counter += evens[i];
    }
    return counter;
}
int iseven (int x){
    int y;
    if (x<10){
        y = x+4;
    }
    else{
        y = x;
    }
    int z = y%2;
    if (z == 0){
        return 0;
    }
    return 1;
}
