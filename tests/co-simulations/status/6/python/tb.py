import sys
import traceback

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    main, _ = afbd.generate(iface, REG_JSON)

    print("Testing whole array read")
    data = main.Sts.read()
    assert data[0] == main.S0
    assert data[1] == main.S1
    assert data[2] == main.S2

    print("\nTesting index read")
    assert main.Sts.read(0) == main.S0
    assert main.Sts.read(1) == main.S1
    assert main.Sts.read(2) == main.S2

    data = main.Sts.read([2, 0, 1])
    assert data[0] == main.S2
    assert data[1] == main.S0
    assert data[2] == main.S1

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
