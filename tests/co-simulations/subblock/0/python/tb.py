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

    subblocks = [main.blk0, main.blk1, main.blk1.blk2]

    for i, sb in enumerate(subblocks):
        print(f"Testing access to blk{i}")
        r = random.randrange(0, 2 ** 32 - 1)
        print(f"Writing value {r} to cfg register")
        sb.cfg.write(r)

        print(f"Reading cfg register")
        read = sb.cfg.read()
        assert read == r, f"Read wrong value from cfg register {read}"

        print(f"Reading st register")
        read = sb.st.read()
        assert read == r, f"Read wrong value from st register {read}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
