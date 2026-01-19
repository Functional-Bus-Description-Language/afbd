import random
import traceback
import sys

import cosim
import afbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main, _ = afbd.generate(iface, REG_JSON)

    N = random.randint(1, 32)
    print(f"Reading fifo stream {N} times")
    vals = main.fifo.read(N)

    for i, v in enumerate(vals):
        assert i == v[0], f"read {v[0]}, expecting {i}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
