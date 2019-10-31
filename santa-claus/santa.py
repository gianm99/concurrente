# Gian Lucas Martín Chamorro
import threading
import time

REINDEERS = 9
ELVES = 9


wakeUp = threading.Semaphore(0)

reindeersWaiting = 0  # number of reindeers waiting
reindeers = threading.Semaphore(0)  # permission to wait to be hitched by santa
reindeersMutex = threading.Lock()  # mutual exclusion for reindeersWaiting
reindeers_list = {"RUDOLPH", "BLITZEN", "DONDER", "CUPID", "COMET",
                  "VIXEN", "PRANCER", "DANCER", "DASHER"}


elvesWaiting = 0  # number of elves waiting
elves = threading.Semaphore(0)  # permission to wait to be attended by santa
elvesMutex = threading.Lock()  # mutual exclusion for elvesWaiting
elfGroup = threading.Semaphore(1)  # permission to wait in a group
elves_list = {"Chaenath", "Elrond", "Hycis", "Imryll", "Galadriel",
              "Arwen", "Tauriel", "Esiyae", "Legolas"}


def santa():
    global reindeersWaiting, elvesWaiting

    while True:
        wakeUp.acquire()
        if reindeersWaiting == REINDEERS:
            with reindeersMutex:
                reindeersWaiting = 0
            for i in range(REINDEERS):
                reindeers.release()
        else:
            with elvesMutex:
                elvesWaiting = 0
            for i in range(ELVES):
                elves.release()
            elfGroup.release()


def elf():
    global elvesWaiting

    print(threading.currentThread().getName())
    while True:
        elfGroup.acquire()
        with elvesMutex:
            elvesWaiting += 1
            if elvesWaiting == 3:
                wakeUp.release()
            else:
                elfGroup.release()
        elves.acquire()


def reindeer():
    global reindeersWaiting

    print(threading.currentThread().getName())
    while True:
        # holidays()
        with reindeersMutex:
            reindeersWaiting += 1
            if reindeersWaiting == 9:
                wakeUp.release()
        reindeers.acquire()
        # deliverToys()


def main():
    threads = []

    # Santa
    s = threading.Thread(target=santa)
    threads.append(s)

    # Reindeers
    for i in reindeers_list:
        r = threading.Thread(target=reindeer, name=i)
        threads.append(r)

    # Elves
    for i in elves_list:
        e = threading.Thread(target=elf, name=i)
        threads.append(e)

    # Start all threads
    for t in threads:
        t.start()

    # Wait for all threads to complete
    for t in threads:
        t.join()


if __name__ == "__main__":
    main()
