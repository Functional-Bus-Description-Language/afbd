import sys
import traceback

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

CLK_PERIOD = 10

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    main, _ = afbd.generate(iface, REG_JSON)

    id = main.ID.read()
    assert id == main.ID.value, f"Read wrong ID {id}, expecting {main.ID.value}"
    print(f"ID: {id}\n")

    ts = main.TIMESTAMP.read()
    assert ts == main.TIMESTAMP.value, f"Read wrong TIMESTAMP {ts}, expecting {main.TIMESTAMP.value}"
    print(f"Timestamp: {ts}\n")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
