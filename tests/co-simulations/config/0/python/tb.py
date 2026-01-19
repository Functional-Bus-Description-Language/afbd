import sys
import random

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    main, _ = afbd.generate(iface, REG_JSON)

    val = random.randint(0, 2 ** 7 - 1)

    print(f"Generated random value: {val}")

    print("Writing Cfg")
    main.Cfg.write(val)

    print("Reading Cfg")
    read_val = main.Cfg.read()
    if read_val != val:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    print("Reading St")
    read_val = main.St.read()
    if read_val != val:
        raise Exception(f"Read wrong value form St {read_val}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(E)
