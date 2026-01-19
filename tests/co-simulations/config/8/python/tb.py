from random import randint
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

    print("\n\nList Test")
    data = []
    for _ in range(len(main.Cfgs)):
        data.append(randint(0, 2**main.Cfgs.width - 1))
    main.Cfgs.write(data)
    rdata = main.Cfgs.read()
    assert rdata == data, f"invalid data read, got {rdata}, want {data}"

    # Clear data
    data = [0 for _ in range(len(main.Cfgs))]
    main.Cfgs.write(data)

    print("\n\nDictionary Test")
    data = {0: 12, 5: 31, 8:79}
    main.Cfgs.write(data)
    rdata = main.Cfgs.read()
    for i in range(len(main.Cfgs)):
        if i in data:
            assert rdata[i] == data[i], f"{i}: got {rdata[0]}, want {data[0]}"
        else:
            assert rdata[i] == 0, f"{i}: got {rdata[0]}, want 0"

    # Clear data
    data = [0 for _ in range(len(main.Cfgs))]
    main.Cfgs.write(data)

    print("\n\nOffset Test")
    offset = 3
    data = []
    for _ in range(len(main.Cfgs) - offset):
        data.append(randint(0, 2**main.Cfgs.width - 1))

    main.Cfgs.write(data, offset)
    rdata = main.Cfgs.read()
    for i in range(len(main.Cfgs)):
        if i < offset:
            assert rdata[i] == 0, f"{i}: got {rdata[i]}, want 0"
        else:
            assert rdata[i] == data[i - offset], f"{i}: got {rdata[i]}, want {data[i - offset]}"

    iface.end(0)

except Exception as E:
    print(traceback.format_exc())
    iface.end(1)
