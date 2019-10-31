import threading

# És concurrent intercalant els dos fils

MAXCOUNT = 100
THREADS = 2
SUMA = 0
TURNO = 0
threads = []

def thread():
    global TURNO
    global SUMA
    
    target = int(MAXCOUNT / THREADS)
    otroproceso = (TURNO + 1) % THREADS

    for j in range(target):
        while TURNO == otroproceso:
            pass
        SUMA += 1
        print("Hi, I'm the thread %s and I added 1 to %d" % (threading.currentThread().name,SUMA))
        TURNO = otroproceso

for i in range(THREADS):
    # Create new threads
    t = threading.Thread(target=thread)
    threads.append(t)
    t.start() # start the thread

for i in threads:
    i.join()

print(SUMA)
