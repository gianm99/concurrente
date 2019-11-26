# Gian Lucas Martín Chamorro
import threading
import time
import random

TURNS = 6  # number of times Santa has to answer the elves' questions
LOADS = 1  # number of times Santa has to load the toys
wakeUp = threading.Semaphore(0)


REINDEERS = 9  # number of reindeers
DELIVERIES = 1  # number of deliveries to be made
reindeersWaiting = 0  # number of reindeers waiting
reindeers = threading.Semaphore(0)  # permission to wait to be hitched by santa
reindeersGroup = threading.Semaphore(1) # mutual exclusion for reindeersWaiting
reindeers_list = {"RUDOLPH", "BLITZEN", "DONDER", "CUPID", "COMET",
                  "VIXEN", "PRANCER", "DANCER", "DASHER"}


ELVES = 9  # number of elfs
ELF_GROUP = 3  # size of waiting group
QUESTIONS = 2  # number of questions to ask to santa
elvesWaiting = 0  # number of elves waiting
elves = threading.Semaphore(0)  # permission to wait to be attended by santa
elfGroup = threading.Semaphore(1)  # permission to wait in a group
elves_list = {"Chaenath", "Elrond", "Hycis", "Imryll", "Galadriel",
              "Arwen", "Tauriel", "Esiyae", "Legolas"}

TIME_SLEEP_SANTA = 1
TIME_LOAD = 1
TIME_HOLIDAYS = 5
TIME_DELIVER = 4
TIME_WORK = 1
TIME_HELP = 1

def santa():
    global reindeersWaiting, elvesWaiting
    turns = 0
    loads = 0

    print("--------> Santa says: I'm tired")
    while turns < TURNS or loads < LOADS:
        print("--------> Santa says: I'm going to sleep")
        wakeUp.acquire()
        sleepSanta()
        print("--------> Santa says: I'm awake ho ho ho!")
        if reindeersWaiting == REINDEERS:
            print("--------> Santa says: Toys are ready!")
            print("--------> Santa loads the toys")
            loadToys()
            # reindeersGroup.acquire()
            reindeersWaiting = 0
            # reindeersGroup.release()
            print("--------> Santa says: Until next Christmas!")
            for i in range(REINDEERS):
                reindeers.release()
            loads += 1
        else:
            print("--------> Santa says: What is the problem?")
            for i in range(ELF_GROUP):
                print("--------> Santa helps the elf {} of 3".format(i+1))
                helpElf()
                elves.release()
            # elvesMutex.acquire()
            elvesWaiting = 0
            # elvesMutex.release()
            turns += 1
            print("--------> Santa ends turn {}".format(turns))
            elfGroup.release()
    print("--------> Santa ends")


def elf():
    global elvesWaiting

    name = threading.currentThread().getName()
    print("Hi I am the elf {}".format(name))
    for i in range(QUESTIONS):
        work()
        elfGroup.acquire()  # acquire
        elvesWaiting += 1
        if elvesWaiting < 3:
            elfGroup.release()  #release
            print("Elf {} says: I have a question, I'm the {} waiting...".format(name, elvesWaiting))
        else:
            print("Elf {} says: I have a question, I'm the 3 waiting SANTAAA!".format(name))
            wakeUp.release()
        elves.acquire()
        print("The elf {} is getting help".format(name))
        getHelp()
    print("Elf {} ends".format(name))


def reindeer():
    global reindeersWaiting

    name = threading.currentThread().getName()
    print("\t\t{} here!".format(name))
    for i in range(DELIVERIES):
        holidays()
        reindeersGroup.acquire()
        reindeersWaiting += 1
        if reindeersWaiting < 9:
            reindeersGroup.release()
            print("\t\tReindeer {} arrives".format(name))
        else:
            # reindeersGroup.release()
            print("\t\tReindeer {} I'm the 9".format(name))
            wakeUp.release()
        reindeers.acquire()
        print("\t\t{} ready and hitched".format(name))
        deliverToys()
    print("\t\tReindeer {} ends".format(name))

def sleepSanta():
    time.sleep(random.uniform(TIME_SLEEP_SANTA,TIME_SLEEP_SANTA+2))

def deliverToys():
    time.sleep(random.uniform(TIME_DELIVER,TIME_DELIVER+2))


def work():
    time.sleep(random.uniform(TIME_WORK,TIME_WORK+2))


def getHelp():
    time.sleep(TIME_HELP)

def helpElf():
    time.sleep(TIME_HELP)

def loadToys():
    time.sleep(random.uniform(TIME_LOAD,TIME_LOAD+2))


def holidays():
    time.sleep(random.uniform(TIME_HOLIDAYS,TIME_HOLIDAYS+2))


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
