import random
import sys
import traceback

import cosim
import afbd

WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]
CONST_JSON = sys.argv[4]

try:
    iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

    Main, const = afbd.generate(iface, REG_JSON, CONST_JSON)

    data = []
    for i in range(const['main']['DEPTH']):
        dataset = []
        dataset.append(random.randint(0, 2 ** Main.Add.params[0]['Width'] - 1))
        dataset.append(random.randint(0, 2 ** Main.Add.params[1]['Width'] - 1))
        dataset.append(random.randint(0, 2 ** Main.Add.params[2]['Width'] - 1))
        data.append(dataset)

    print(f"Writing downstream {const['main']['DEPTH']} times")
    Main.Add.write(data)

    results = Main.Result.read(const['main']['DEPTH'])

    for i, dataset in enumerate(data):
        got = results[i][0]
        want = sum(dataset)
        assert got == want, f"{i}: got {got}, want {want}"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
