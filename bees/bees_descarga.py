# Gian Lucas Martín Chamorro
import threading
import time

BEAR = 1
N_BEES = 10
H = 10
TO_EAT = 10
TO_PRODUCE = 10

TIME_TO_EAT = 0.001
TIME_TO_PRODUCE = 0.002

honeyJar = 0
full = threading.Semaphore(0)
notFull = threading.Semaphore(1)

def eat():
    time.sleep(TIME_TO_EAT)

def produce():
    time.sleep(TIME_TO_PRODUCE)

def bear():
    global honeyJar

    for i in range(TO_EAT):
        full.acquire()
        eat()
        print("Bear: Honey jar is empty")
        honeyJar = 0
        notFull.release()   # Honey jar is empty

def bee():
    global honeyJar

    for i in range(TO_PRODUCE):
        notFull.acquire()
        produce()
        honeyJar += 1
        print ("Honey jar now has {}".format(honeyJar))
        if honeyJar == H:
            print("Bees: Honey jar is full")
            full.release()  # Honey jar is full
        else:
            notFull.release()   # Honey jar is not full

def main():
    threads = []

    # Bees
    for i in range(N_BEES):
        b1 = threading.Thread(target=bee)
        threads.append(b1)
    
    # Bear
    b2 = threading.Thread(target=bear)
    threads.append(b2)

    # Start all threads
    for t in threads:
        t.start()
    
    # Wait for all threads to complete
    for t in threads:
        t.join()
    
    print("End")

if __name__ == "__main__":
    main()