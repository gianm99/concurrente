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


class BearBees(object):
    def __init__(self):
        self.honeyJar = int
        self.mutex = threading.Lock()
        self.notFull = threading.Condition(self.mutex)
        self.Full = threading.Condition(self.mutex)

    def produce(self):
        with self.mutex:
            while self.honeyJar < H:
                self.notFull.wait()
            self.honeyJar += 1
    
    def eat(self):
        with self.mutex:
            while self