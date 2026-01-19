import sys
import traceback
import random

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    main, _ = afbd.generate(iface, REG_JSON)

    lower = random.randint(0, 2 ** 30  - 1)
    upper = random.randint(0, 2 ** 20  - 1)

    print(f"Generated random values: lower = {lower}, upper = {upper}")

    print("Writing Lower")
    main.Lower.write(lower)

    print("Writing Upper")
    main.Upper.write(upper)

    print("Reading St")
    st = main.St.read()
    if st != (upper << 30) | lower:
        raise Exception(f"Read wrong value form St {st}, expects {(upper << 30) | lower}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
