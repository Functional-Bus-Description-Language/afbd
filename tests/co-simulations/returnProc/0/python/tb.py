import random
import traceback
import sys

import cosim
import afbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]


try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    main, _ = afbd.generate(iface, REG_JSON)

    for i in range(3):
        want = random.randint(0, 2 ** 32 - 1)

        print(f"writing want = {want}")
        main.cfg.write(want)

        print("calling proc")
        got = main.my_proc()
        assert got == want, f"got {got}, want {want}"

        exit_cnt = main.exit_cnt.read()
        assert exit_cnt == i + 1, f"exit_cnt = {exit_cnt}, want {i+1}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
