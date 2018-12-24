#!/usr/bin/env python3

# chip-8.py - main executable file for the chip-8 emulator

import sys

import numpy as np
import cv2 as cv


# 4K memory used by the emulator
# +---------------+= 0xFFF (4095) End of Chip-8 RAM
# |               |
# |               |
# |               |
# |               |
# |               |
# | 0x200 to 0xFFF|
# |     Chip-8    |
# | Program / Data|
# |     Space     |
# |               |
# |               |
# |               |
# +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
# |               |
# |               |
# |               |
# +---------------+= 0x200 (512) Start of most Chip-8 programs
# | 0x000 to 0x1FF|
# | Reserved for  |
# |  interpreter  |
# +---------------+= 0x000 (0) Start of Chip-8 RAM
mem = np.empty(4096, dtype=np.uint8)

registers = np.zeros(16, dtype=np.uint8)


class CPU:
    def __init__(self):
        self.PC = 0  # program counter
        self.SP = 0  # stack pointer
        self.DT = 0  # delay timer
        self.ST = 0  # sound timer
        self.V = np.zeros(16, dtype=np.uint8)  # general purpose registers


