# Gian Lucas Martín Chamorro
import threading
import time

TURNS = 6  # number of times Santa has to answer the elves' questions
LOADS = 1  # number of times Santa has to load the toys
wakeUp = threading.Semaphore(0)


REINDEERS = 9  # number of reindeers
DELIVERIES = 1  # number of deliveries to be made
reindeersWaiting = 0  # number of reindeers waiting
reindeers = threading.Semaphore(0)  # permission to wait to be hitched by santa
reindeersMutex = threading.Lock()  # mutual exclusion for reindeersWaiting
reindeers_list = {"RUDOLPH", "BLITZEN", "DONDER", "CUPID", "COMET",
                  "VIXEN", "PRANCER", "DANCER", "DASHER"}


ELVES = 9  # number of elfs
ELF_GROUP = 3  # size of waiting group
QUESTIONS = 2  # number of questions to ask to santa
elvesWaiting = 0  # number of elves waiting
elves = threading.Semaphore(0)  # permission to wait to be attended by santa
elvesMutex = threading.Lock()  # mutual exclusion for elvesWaiting
elfGroup = threading.Semaphore(1)  # permission to wait in a group
elves_list = {"Chaenath", "Elrond", "Hycis", "Imryll", "Galadriel",
              "Arwen", "Tauriel", "Esiyae", "Legolas"}


def santa():
    global reindeersWaiting, elvesWaiting
    tu=0
    lo=0

    print("--------> Santa says: I'm tired")
    while tu < TURNS or lo < LOADS:
        print("--------> Santa says: I'm going to sleep")
        # sleepSanta()
        wakeUp.acquire()
        print("--------> Santa says: I'm awake ho ho ho!")
        if reindeersWaiting == REINDEERS:
            with reindeersMutex:
                reindeersWaiting = 0
            # load()
            for i in range(REINDEERS):
                reindeers.release()
            lo += 1
        else:
            with elvesMutex:
                elvesWaiting = 0
            elfGroup.release()
            print("--------> Santa says: What is the problem?")
            for i in range(ELF_GROUP):
                print("--------> Santa helps the elf {} of 3".format(i+1))
                # help()
                elves.release()
            tu += 1
            print("--------> Santa ends turn {}".format(tu))
    print("--------> Santa ends")


def elf():
    global elvesWaiting
    
    print("Hi I am the elf {}".format(threading.currentThread().getName()))
    for i in range(QUESTIONS):
        elfGroup.acquire()
        with elvesMutex:
            elvesWaiting += 1
            if elvesWaiting == 3:
                print("Elf {} says: I have a question, I'm the 3 waiting SANTAAA!".format(
                    threading.currentThread().getName()))
                wakeUp.release()
            else:
                print("Elf {} says: I have a question, I'm the {} waiting...".format(
                    threading.currentThread().getName(), elvesWaiting))
                elfGroup.release()
        elves.acquire()
        print("The elf {} got help from Santa".format(
            threading.currentThread().getName()))
        # work()

    print("Elf {} ends".format(threading.currentThread().getName()))


def reindeer():
    global reindeersWaiting

    print("\t\t{} here!".format(threading.currentThread().getName()))
    for i in range(DELIVERIES):
        # holidays()
        with reindeersMutex:
            reindeersWaiting += 1
            if reindeersWaiting == 9:
                print("\t\tReindeer {} I'm the 9".format(
                    threading.currentThread().getName()))
                wakeUp.release()
            else:
                print("\t\tReindeer {} arrives".format(
                    threading.currentThread().getName()))
        reindeers.acquire()
        print("\t\t{} ready and hitched".format(
            threading.currentThread().getName()))
        # deliverToys()
    print("\t\tReindeer {} ends".format(threading.currentThread().getName()))


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
