import sys
import traceback

import cosim
import afbd


WRITE_FIFO_PATH = sys.argv[1]
READ_FIFO_PATH = sys.argv[2]
REG_JSON = sys.argv[3]

iface = cosim.Iface(WRITE_FIFO_PATH, READ_FIFO_PATH)

try:
    Main, _ = afbd.generate(iface, REG_JSON)

    max_val = 2 ** Main.WIDTH - 1

    print("\ntesting set()")
    print("\ntesting settting all bits")
    Main.Mask.set()
    got = Main.Mask.read()
    assert got == max_val, f"got {got}, want {max_val}"
    got = Main.St.read()
    assert got == max_val, f"got {got}, want {max_val}"
    Main.Mask.set([])
    print("\ntesting settting one bit")
    Main.Mask.set(4)
    got = Main.Mask.read()
    assert got == 1 << 4, f"got {got}, want 4"
    got = Main.St.read()
    assert got == 1 << 4, f"got {got}, want 4"
    Main.Mask.set([])
    print("\ntesting settting two bits")
    Main.Mask.set([0, 2])
    got = Main.Mask.read()
    assert got == 5, f"got {got}, want 5"
    got = Main.St.read()
    assert got == 5, f"got {got}, want 5"
    Main.Mask.set([])

    print("\n\ntesting clear()")
    print("\ntesting clearing all bits")
    Main.Mask.set()
    Main.Mask.clear()
    got = Main.Mask.read()
    assert got == 0, f"got {got}, want 0"
    got = Main.St.read()
    assert got == 0, f"got {got}, want 0"
    print("\ntesting clearing single bit")
    Main.Mask.set()
    Main.Mask.clear(0)
    got = Main.Mask.read()
    assert got == 0b1111110, f"got {got:#b}, want 0b1111110"
    got = Main.St.read()
    assert got == 0b1111110, f"got {got:#b}, want 0b1111110"
    print("\ntesting clearing two bits")
    Main.Mask.set()
    Main.Mask.clear([0, 3])
    got = Main.Mask.read()
    assert got == 0b1110110, f"got {got:#b}, want 0b1110110"
    got = Main.St.read()
    assert got == 0b1110110, f"got {got:#b}, want 0b1110110"

    print("\n\ntesting toggle()")
    Main.Mask.clear()
    print("\ntesting toggle all bits")
    Main.Mask.toggle()
    got = Main.Mask.read()
    assert got == 0b1111111, f"got {got:#b}, want 0b1111111"
    got = Main.St.read()
    assert got == 0b1111111, f"got {got:#b}, want 0b1111111"
    print("\ntesting toggle one bit")
    Main.Mask.toggle(3)
    got = Main.Mask.read()
    assert got == 0b1110111, f"got {got:#b}, want 0b1110111"
    got = Main.St.read()
    assert got == 0b1110111, f"got {got:#b}, want 0b1110111"
    print("\ntesting toggle two bits")
    Main.Mask.toggle([1, 2])
    got = Main.Mask.read()
    assert got == 0b1110001, f"got {got:#b}, want 0b1110001"
    got = Main.St.read()
    assert got == 0b1110001, f"got {got:#b}, want 0b1110001"


    print("\n\ntesting update_set()")
    Main.Mask.clear()
    print("\ntesting update set one bit")
    Main.Mask.update_set(2)
    got = Main.Mask.read()
    assert got == 0b0000100, f"got {got:#b}, want 0b0000100"
    got = Main.St.read()
    assert got == 0b0000100, f"got {got:#b}, want 0b0000100"
    print("\ntesting update set two bits")
    Main.Mask.update_set([4 ,6])
    got = Main.Mask.read()
    assert got == 0b1010100, f"got {got:#b}, want 0b1010100"
    got = Main.St.read()
    assert got == 0b1010100, f"got {got:#b}, want 0b1010100"

    print("\n\ntesting update_clear()")
    Main.Mask.set()
    print("\ntesting update clear one bit")
    Main.Mask.update_clear(2)
    got = Main.Mask.read()
    assert got == 0b1111011, f"got {got:#b}, want 0b1111011"
    got = Main.St.read()
    assert got == 0b1111011, f"got {got:#b}, want 0b1111011"
    print("\ntesting update clear two bits")
    Main.Mask.update_clear([0, 6])
    got = Main.Mask.read()
    assert got == 0b0111010, f"got {got:#b}, want 0b0111010"
    got = Main.St.read()
    assert got == 0b0111010, f"got {got:#b}, want 0b0111010"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
