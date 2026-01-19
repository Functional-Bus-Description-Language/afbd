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

    vec = []
    for _ in range(10):
        vec.append(random.randint(0, 2 ** 17 - 1))
    sum = sum(vec)

    print(f"Calling add function, vec = {vec}")
    main.add(vec)

    print(f"Reading result")
    result = main.result.read()

    if result != sum:
        print(f"Wrong result, got {result}, expecting {sum}")
        iface.end(1)

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
