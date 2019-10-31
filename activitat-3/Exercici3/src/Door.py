from datetime import timedelta
from typing import List

from src.Entry import Entry


class Door:

    def __init__(self, id: int):
        self.id = id
        self.entries: List[Entry] = []

    def get_average_entry_time(self):
        total_time = timedelta(0)
        for i in range(len(self.entries)):
            if i > 0:
                total_time += (self.entries[i].time - self.entries[i - 1].time)
        return total_time / (len(self.entries) - 1)
