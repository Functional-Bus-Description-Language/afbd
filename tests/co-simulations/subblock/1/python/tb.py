import random
import sys
import traceback

import cosim
import afbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main, _ = afbd.generate(iface, REG_JSON)

    sum = 0
    for blk in main.Blk:
        x = random.randint(0, 2**blk.X.width-1)
        blk.X.write(x)
        sum += x

    read_sum = main.Sum.read()

    assert read_sum == sum, f"wrong read sum, got {read_sum}, want {sum}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
