import sys

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = afbd.Main(iface)

    print(f"Writing VALID_VALUE ({afbd.mainPkg.VALID_VALUE}) to Cfg register")
    Main.Cfg.write(afbd.mainPkg.VALID_VALUE)

    print("Reading Cfg")
    read_val = Main.Cfg.read()
    if read_val != afbd.mainPkg.VALID_VALUE:
        raise Exception(f"Read wrong value form Cfg {read_val}")

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(E)
