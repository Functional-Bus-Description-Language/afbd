import sys

import cosim
import afbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main, _ = afbd.generate(iface, REG_JSON)

    values = main.status_array.read()
    assert len(values) == 9
    for i, v in enumerate(values):
        assert v == i

    idx = [2, 7]
    values = main.status_array.read(idx)
    assert values[0] == 2
    assert values[1] == 7

    value = main.status_array.read(5)
    assert value == 5

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
