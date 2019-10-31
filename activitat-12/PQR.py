# Sincronia amb semàfors: Implementar la simulació consistent en tres
# processos que van imprimint "P", "Q" i "R" respectivament, un nombre
# determinat de vegades. S'ha de verificar sempre que el nombre de "R"
# ha de ser inferior a la suma del nombre de "P" + "Q".

# Gian Lucas Martín Chamorro

import threading
import time

TO_PRINT_P = 10
TO_PRINT_Q = 10
TO_PRINT_R = 19

TIME_TO_P = 0.006
TIME_TO_Q = 0.005
TIME_TO_R = 0.001

lessThanSum = threading.Semaphore(0)
mutex = threading.Lock()
sum = 0
n_R = 0


def P():
    global sum, n_R
    for i in range(TO_PRINT_P):
        time.sleep(TIME_TO_P)
        print('P')
        with mutex:
            sum += 1
            if n_R + 1 < sum:
                print("sum is {} and n_R is {}".format(sum,n_R))
                lessThanSum.release()


def Q():
    global sum, n_R
    for i in range(TO_PRINT_Q):
        time.sleep(TIME_TO_Q)
        print('Q')
        with mutex:
            sum += 1
            if n_R + 1 < sum:
                print("sum is {} and n_R is {}".format(sum,n_R))
                lessThanSum.release()


def R():
    global n_R
    for i in range(TO_PRINT_R):
        lessThanSum.acquire()
        time.sleep(TIME_TO_R)
        print('R')
        n_R+=1


def main():
    threads = []

    #P
    p = threading.Thread(target=P)
    threads.append(p)
    #Q
    q = threading.Thread(target=Q)
    threads.append(q)
    #R
    r = threading.Thread(target=R)
    threads.append(r)

    # Start all threads
    for t in threads:
        t.start()
    
    # Wait for all threads to complete
    for t in threads:
        t.join()
    
    print("End")


if __name__ == "__main__":
    main()
