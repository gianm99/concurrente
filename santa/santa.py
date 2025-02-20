# Autor: Gian Lucas Martín Chamorro
# Vídeo: https://youtu.be/u2MhrgO5pbw
import threading
import time
import random

TURNS = 6  # number of times Santa has to answer the elves' questions
LOADS = 1  # number of times Santa has to load the toys
wakeUp = threading.Semaphore(0) # sleep until being woken up


REINDEERS = 9  # number of reindeers
DELIVERIES = 1  # number of deliveries to be made
reindeersWaiting = 0  # number of reindeers waiting
reindeers = threading.Semaphore(0)  # permission to wait for santa
reindeersGroup = threading.Semaphore(1) # permission to wait in a group
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

TIME_SLEEP_SANTA = 1.5
TIME_LOAD = 0.5
TIME_HOLIDAYS = 8
TIME_DELIVER = 1
TIME_WORK = 0.5
TIME_HELP = 0.5

def santa():
    global reindeersWaiting, elvesWaiting
    turns = 0
    loads = 0

    print("--------> Santa says: I'm tired")
    while turns < TURNS or loads < LOADS:
        print("--------> Santa says: I'm going to sleep")
        sleepSanta()
		# Espera a ser despertado
        wakeUp.acquire()
        print("--------> Santa says: I'm awake ho ho ho!")
        if reindeersWaiting == REINDEERS:
            print("--------> Santa says: Toys are ready!")
            print("--------> Santa loads the toys")
			# Libera a cada reno y le carga los juguetes
            for i in range(REINDEERS):
                reindeers.release()
                loadToys()
			# No queda ningún reno esperando
            reindeersWaiting = 0
            print("--------> Santa says: Until next Christmas!")
			# Se libera al grupo para que se puedan poner a esperar otra vez
            reindeersGroup.release()
            loads += 1
        else:
            print("--------> Santa says: What is the problem?")
			# Libera a los tres elfos y los ayuda
            for i in range(ELF_GROUP):
                print("--------> Santa is helping the elf {} of 3".format(i+1))
                elves.release()
                helpElf()
			# No queda ningun elfo del grupo esperando
            elvesWaiting = 0
            print("--------> Santa ends turn {}".format(turns+1))
			# Se libera al grupo para que se puedan poner a esperar otros 3 elfos
            elfGroup.release()
            turns += 1
    print("--------> Santa ends")


def elf():
    global elvesWaiting

    name = threading.currentThread().getName()
    print("Hi I am the elf {}".format(name))
    for i in range(QUESTIONS):
        work()
        elfGroup.acquire()
        elvesWaiting += 1
        if elvesWaiting < 3:
			# Si no es el tercero libera el semáforo
            print("Elf {} says: I have a question, I'm the {} waiting...".format(name, elvesWaiting))
            elfGroup.release()
        else:
			# Si es el tercero no libera el semáforo y despierta a Santa
            print("Elf {} says: I have a question, I'm the 3 waiting SANTAAA!".format(name))
            wakeUp.release()
		# Espera a que le atienda Santa
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
			# Si no es el noveno libera el semáforo
            reindeersGroup.release()
            print("\t\tReindeer {} arrives".format(name))
        else:
			# Si es el noveno no libera el semáforo y despierta a Santa 
            print("\t\tReindeer {} I'm the 9".format(name))
            wakeUp.release()
		# Espera a que le atienda Santa
        reindeers.acquire()
        loadToys()
        print("\t\t{} ready and hitched".format(name))
    print("\t\tReindeer {} ends".format(name))

def sleepSanta():
    time.sleep(random.uniform(TIME_SLEEP_SANTA,TIME_SLEEP_SANTA+1))

def work():
    time.sleep(random.uniform(TIME_WORK,TIME_WORK+1))

def getHelp():
    time.sleep(TIME_HELP)

def helpElf():
    time.sleep(TIME_HELP)

def loadToys():
    time.sleep(TIME_LOAD)

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

main()
