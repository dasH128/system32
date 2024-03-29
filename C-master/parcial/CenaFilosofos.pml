#define wait(s) atomic {s>0 -> s--}
#define signal(s) s++
#define N 5

byte fork[N] = {1,1,1,1,1}

active[N-1] proctype p1(){
    do
    ::
    printf("Pensando\n")
    wait(fork[_pid])
    wait(fork[_pid+1])
    printf("Comiendo\n")
    signal(fork[_pid])
    signal(fork[_pid+1])
    od
}

active proctype p2(){
    do
    ::
        printf("L Pensando\n")
        wait(fork[0])
        wait(fork[N-1])
        printf("L Comiendo\n")
        signal(fork[0])
        signal(fork[N-1])
    od
}