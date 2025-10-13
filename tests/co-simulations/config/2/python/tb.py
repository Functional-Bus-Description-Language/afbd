import sys

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]
CONST_JSON = sys.argv[4]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main, const = afbd.generate(iface, REG_JSON, CONST_JSON)

    print(f"Writing VALID_VALUE ({const['main']['VALID_VALUE']}) to Cfg register")
    Main.Cfg.write(const['main']['VALID_VALUE'])

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != const['main']['VALID_VALUE']:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(E)
