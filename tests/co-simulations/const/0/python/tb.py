import sys
import traceback

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main = afbd.Main(iface)


    print("\n\nTesting int constant")
    print("Reading St register")
    read = Main.St.read()
    assert (
        read == afbd.mainPkg.C
    ), f"read value {read} differs from constant value {afbd.mainPkg.C}"


    print("\n\nTesting int list constants")
    print("Reading Stl register")
    read = Main.Stl.read()
    for i, v in enumerate(afbd.mainPkg.CL):
        assert (
            read[i] == v
        ), f"read value {read[i]} differs from constant value {v}"


    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
