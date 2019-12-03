import threading
import time
from datetime import datetime
import random
from typing import List

from src.Door import Door
from src.Entry import Entry


# Park variables
NUM_DOORS = 2
ENTRIES_BY_DOOR = 10
MAX_DELAY = 5000  # milliseconds
num_entries = 0
doors: List[Door] = []

# Concurrency variables - Peterson's algorithm
NUM_THREADS = NUM_DOORS
want = [False, False]
last = 0


def enter():
    global num_entries, last, want

    current_thread = int(threading.current_thread().name)
    other_thread = (current_thread + 1) % NUM_THREADS

    current_door = doors[current_thread]
    print("Door {}".format(current_door.id))

    for i in range(ENTRIES_BY_DOOR):

        # Delay
        delay = random.randint(1, MAX_DELAY)
        time.sleep(delay / 1000)

        # Preprotocol
        want[current_thread] = True
        last = current_thread
        while not ((not want[other_thread]) or (last == other_thread)):
            pass

        # CS
        num_entries += 1
        entry = Entry(num_entries, datetime.now())
        current_door.entries.append(entry)
        print("Door {} ({} entries). Entry {} at: {}".format(current_door.id, len(current_door.entries), entry.id, entry.time))

        # Postprotocol
        want[current_thread] = False


def main():
    threads = []

    # Create doors and threads
    for i in range(NUM_DOORS):
        door = Door(i)
        doors.append(door)

        thread = threading.Thread(target=enter)
        thread.name = i
        threads.append(thread)

        thread.start()

    # Wait for threads
    for thread in threads:
        thread.join()

    # Final information
    print("Total entries: {}".format(num_entries))
    for door in doors:
        print("Average entry time door {}: {}".format(door.id, door.get_average_entry_time()))


if __name__ == "__main__":
    main()
