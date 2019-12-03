import threading

# No intercala els fils, fa primer un bloc de fils i després l'laltre

MAXCOUNT = 1000
THREADS = 2
SUMA = 0
threads = []
WANT = []

target = int(MAXCOUNT / THREADS)

def thread():
    global SUMA
    global WANT

    ident = int(threading.current_thread().name)
    otro = (ident + 1) % THREADS

    for j in range(target):
        while WANT[otro]:
            pass
        WANT[ident] = True
        SUMA += 1
        print("Hi, I'm the thread %d and I added 1 to %d" % (ident,SUMA))
        WANT[ident] = False


for i in range(THREADS):
    # Create new threads
    t = threading.Thread(target=thread, name=str(i))
    threads.append(t)
    WANT.append(False)

for t in threads:
    t.start()  # start the thread

for i in threads:
    i.join()

print(SUMA)
