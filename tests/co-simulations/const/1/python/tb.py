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


    print("\n\nTesting int constant")
    print("Reading st register")
    read = main.st.read()
    assert read == main.C, f"read value {read} differs from constant value {main.C}"


    print("\n\nTesting int list constants")
    print("Reading stl register")
    read = main.stl.read()
    for i, v in enumerate(main.CL):
        assert (
            read[i] == v
        ), f"read value {read[i]} differs from constant value {v}"


    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
