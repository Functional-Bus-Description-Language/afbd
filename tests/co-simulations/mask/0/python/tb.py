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

    max_val = 2 ** main.WIDTH - 1

    print("\ntesting set()")
    print("\ntesting settting all bits")
    main.Mask.set()
    got = main.Mask.read()
    assert got == max_val, f"got {got}, want {max_val}"
    got = main.St.read()
    assert got == max_val, f"got {got}, want {max_val}"
    main.Mask.set([])
    print("\ntesting settting one bit")
    main.Mask.set(4)
    got = main.Mask.read()
    assert got == 1 << 4, f"got {got}, want 4"
    got = main.St.read()
    assert got == 1 << 4, f"got {got}, want 4"
    main.Mask.set([])
    print("\ntesting settting two bits")
    main.Mask.set([0, 2])
    got = main.Mask.read()
    assert got == 5, f"got {got}, want 5"
    got = main.St.read()
    assert got == 5, f"got {got}, want 5"
    main.Mask.set([])

    print("\n\ntesting clear()")
    print("\ntesting clearing all bits")
    main.Mask.set()
    main.Mask.clear()
    got = main.Mask.read()
    assert got == 0, f"got {got}, want 0"
    got = main.St.read()
    assert got == 0, f"got {got}, want 0"
    print("\ntesting clearing single bit")
    main.Mask.set()
    main.Mask.clear(0)
    got = main.Mask.read()
    assert got == 0b1111110, f"got {got:#b}, want 0b1111110"
    got = main.St.read()
    assert got == 0b1111110, f"got {got:#b}, want 0b1111110"
    print("\ntesting clearing two bits")
    main.Mask.set()
    main.Mask.clear([0, 3])
    got = main.Mask.read()
    assert got == 0b1110110, f"got {got:#b}, want 0b1110110"
    got = main.St.read()
    assert got == 0b1110110, f"got {got:#b}, want 0b1110110"

    print("\n\ntesting toggle()")
    main.Mask.clear()
    print("\ntesting toggle all bits")
    main.Mask.toggle()
    got = main.Mask.read()
    assert got == 0b1111111, f"got {got:#b}, want 0b1111111"
    got = main.St.read()
    assert got == 0b1111111, f"got {got:#b}, want 0b1111111"
    print("\ntesting toggle one bit")
    main.Mask.toggle(3)
    got = main.Mask.read()
    assert got == 0b1110111, f"got {got:#b}, want 0b1110111"
    got = main.St.read()
    assert got == 0b1110111, f"got {got:#b}, want 0b1110111"
    print("\ntesting toggle two bits")
    main.Mask.toggle([1, 2])
    got = main.Mask.read()
    assert got == 0b1110001, f"got {got:#b}, want 0b1110001"
    got = main.St.read()
    assert got == 0b1110001, f"got {got:#b}, want 0b1110001"


    print("\n\ntesting update_set()")
    main.Mask.clear()
    print("\ntesting update set one bit")
    main.Mask.update_set(2)
    got = main.Mask.read()
    assert got == 0b0000100, f"got {got:#b}, want 0b0000100"
    got = main.St.read()
    assert got == 0b0000100, f"got {got:#b}, want 0b0000100"
    print("\ntesting update set two bits")
    main.Mask.update_set([4 ,6])
    got = main.Mask.read()
    assert got == 0b1010100, f"got {got:#b}, want 0b1010100"
    got = main.St.read()
    assert got == 0b1010100, f"got {got:#b}, want 0b1010100"

    print("\n\ntesting update_clear()")
    main.Mask.set()
    print("\ntesting update clear one bit")
    main.Mask.update_clear(2)
    got = main.Mask.read()
    assert got == 0b1111011, f"got {got:#b}, want 0b1111011"
    got = main.St.read()
    assert got == 0b1111011, f"got {got:#b}, want 0b1111011"
    print("\ntesting update clear two bits")
    main.Mask.update_clear([0, 6])
    got = main.Mask.read()
    assert got == 0b0111010, f"got {got:#b}, want 0b0111010"
    got = main.St.read()
    assert got == 0b0111010, f"got {got:#b}, want 0b0111010"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
