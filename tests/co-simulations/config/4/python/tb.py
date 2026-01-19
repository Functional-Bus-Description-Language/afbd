import sys
import traceback

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]
CONST_JSON = sys.argv[4]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    main, const = afbd.generate(iface, REG_JSON, CONST_JSON)

    value = 2 ** const['main']['WIDTH'] - 1

    print(f"Writing VALID_VALUE ({value}) to Cfg register")
    main.Cfg.write(value)

    print("Reading Cfg")
    read_val = main.Cfg.read()
    if read_val != value:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    iface.end(0)

except Exception as E:
    print(traceback.format_exc())
    iface.end(1)
