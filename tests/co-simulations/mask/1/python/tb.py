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

    print("\ntesting mask set()")
    print("\ntesting setting all bits")
    Main.Mask.set()
    got = Main.Mask.read()
    assert got == max_val, f"got {got}, want {max_val}"
    got = Main.St.read()
    assert got == max_val, f"got {got}, want {max_val}"
    Main.Mask.set([])
    print("\ntesting setting one bit")
    Main.Mask.set(39)
    got = Main.Mask.read()
    assert got == 0x8000000000, f"got {got:#x}, want 0x8000000000"
    got = Main.St.read()
    assert got == 0x8000000000, f"got {got:#x}, want 0x8000000000"
    print("\ntesting setting two bits")
    Main.Mask.set([0, 32])
    got = Main.Mask.read()
    assert got == 0x0100000001, f"got {got:#x}, want 0x0100000001"
    got = Main.St.read()
    assert got == 0x0100000001, f"got {got:#x}, want 0x0100000001"

    print("\n\ntesting mask clear()")
    print("\ntesting clearing all bits")
    Main.Mask.set()
    Main.Mask.clear()
    got = Main.Mask.read()
    assert got == 0x0, f"got {got:#x}, want 0x0000000000"
    got = Main.St.read()
    assert got == 0x0, f"got {got:#x}, want 0x0000000000"
    print("\ntesting clearing one bit")
    Main.Mask.set()
    Main.Mask.clear(32)
    got = Main.Mask.read()
    assert got == 0xfeffffffff, f"got {got:#x}, want 0fefffffffff"
    got = Main.St.read()
    assert got == 0xfeffffffff, f"got {got:#x}, want 0xfeffffffff"
    print("\ntesting clearing two bits")
    Main.Mask.set()
    Main.Mask.clear([0, 39])
    got = Main.Mask.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    got = Main.St.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"

    print("\n\ntesting mask toggle()")
    print("\ntesting clearing all bits")
    Main.Mask.clear()
    print("\n\ntesting toggling all bits")
    Main.Mask.set(0)
    Main.Mask.toggle()
    got = Main.Mask.read()
    assert got == 0xfffffffffe, f"got {got:#x}, want 0xfffffffffe"
    got = Main.St.read()
    assert got == 0xfffffffffe, f"got {got:#x}, want 0xfffffffffe"
    print("\n\ntesting toggling one bit")
    Main.Mask.toggle(39)
    got = Main.Mask.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    got = Main.St.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    print("\n\ntesting toggling two bits")
    Main.Mask.toggle([4, 32])
    got = Main.Mask.read()
    assert got == 0x7effffffee, f"got {got:#x}, want 0x7effffffee"
    got = Main.St.read()
    assert got == 0x7effffffee, f"got {got:#x}, want 0x7effffffee"

    print("\n\ntesting mask update_set()")
    Main.Mask.set([0])
    print("\n\ntesting update_set one bit")
    Main.Mask.update_set([39])
    got = Main.Mask.read()
    assert got == 0x8000000001, f"got {got:#x}, want 0x8000000001"
    got = Main.St.read()
    assert got == 0x8000000001, f"got {got:#x}, want 0x8000000001"
    print("\n\ntesting update_set two bits")
    Main.Mask.update_set([4, 32])
    got = Main.Mask.read()
    assert got == 0x8100000011, f"got {got:#x}, want 0x8000000011"
    got = Main.St.read()
    assert got == 0x8100000011, f"got {got:#x}, want 0x8000000011"

    print("\n\ntesting mask update_clear()")
    Main.Mask.set()
    print("\n\ntesting update_clear one bit")
    Main.Mask.update_clear([39])
    got = Main.Mask.read()
    assert got == 0x7fffffffff, f"got {got:#x}, want 0x7fffffffff"
    got = Main.St.read()
    assert got == 0x7fffffffff, f"got {got:#x}, want 0x7fffffffff"
    print("\n\ntesting update_clear two bits")
    Main.Mask.update_clear([8, 32])
    got = Main.Mask.read()
    assert got == 0x7efffffeff, f"got {got:#x}, want 0x7efffffeff"
    got = Main.St.read()
    assert got == 0x7efffffeff, f"got {got:#x}, want 0x7efffffeff"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
