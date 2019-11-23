import threading

# No és concurrent, es bloqueja

MAXCOUNT = 10
THREADS = 2
SUMA = 0
threads = []
WANT = []
target = int(MAXCOUNT / THREADS)

def thread():
    global SUMA
    global WANT

    ident = int(threading.current_thread().name)

    for j in range(target):
        WANT[ident] = True
        while WANT[(ident + 1) % THREADS]:
            pass
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