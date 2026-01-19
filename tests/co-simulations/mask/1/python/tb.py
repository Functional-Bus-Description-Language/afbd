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

    print("\ntesting mask set()")
    print("\ntesting setting all bits")
    main.Mask.set()
    got = main.Mask.read()
    assert got == max_val, f"got {got}, want {max_val}"
    got = main.St.read()
    assert got == max_val, f"got {got}, want {max_val}"
    main.Mask.set([])
    print("\ntesting setting one bit")
    main.Mask.set(39)
    got = main.Mask.read()
    assert got == 0x8000000000, f"got {got:#x}, want 0x8000000000"
    got = main.St.read()
    assert got == 0x8000000000, f"got {got:#x}, want 0x8000000000"
    print("\ntesting setting two bits")
    main.Mask.set([0, 32])
    got = main.Mask.read()
    assert got == 0x0100000001, f"got {got:#x}, want 0x0100000001"
    got = main.St.read()
    assert got == 0x0100000001, f"got {got:#x}, want 0x0100000001"

    print("\n\ntesting mask clear()")
    print("\ntesting clearing all bits")
    main.Mask.set()
    main.Mask.clear()
    got = main.Mask.read()
    assert got == 0x0, f"got {got:#x}, want 0x0000000000"
    got = main.St.read()
    assert got == 0x0, f"got {got:#x}, want 0x0000000000"
    print("\ntesting clearing one bit")
    main.Mask.set()
    main.Mask.clear(32)
    got = main.Mask.read()
    assert got == 0xfeffffffff, f"got {got:#x}, want 0fefffffffff"
    got = main.St.read()
    assert got == 0xfeffffffff, f"got {got:#x}, want 0xfeffffffff"
    print("\ntesting clearing two bits")
    main.Mask.set()
    main.Mask.clear([0, 39])
    got = main.Mask.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    got = main.St.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"

    print("\n\ntesting mask toggle()")
    print("\ntesting clearing all bits")
    main.Mask.clear()
    print("\n\ntesting toggling all bits")
    main.Mask.set(0)
    main.Mask.toggle()
    got = main.Mask.read()
    assert got == 0xfffffffffe, f"got {got:#x}, want 0xfffffffffe"
    got = main.St.read()
    assert got == 0xfffffffffe, f"got {got:#x}, want 0xfffffffffe"
    print("\n\ntesting toggling one bit")
    main.Mask.toggle(39)
    got = main.Mask.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    got = main.St.read()
    assert got == 0x7ffffffffe, f"got {got:#x}, want 0x7ffffffffe"
    print("\n\ntesting toggling two bits")
    main.Mask.toggle([4, 32])
    got = main.Mask.read()
    assert got == 0x7effffffee, f"got {got:#x}, want 0x7effffffee"
    got = main.St.read()
    assert got == 0x7effffffee, f"got {got:#x}, want 0x7effffffee"

    print("\n\ntesting mask update_set()")
    main.Mask.set([0])
    print("\n\ntesting update_set one bit")
    main.Mask.update_set([39])
    got = main.Mask.read()
    assert got == 0x8000000001, f"got {got:#x}, want 0x8000000001"
    got = main.St.read()
    assert got == 0x8000000001, f"got {got:#x}, want 0x8000000001"
    print("\n\ntesting update_set two bits")
    main.Mask.update_set([4, 32])
    got = main.Mask.read()
    assert got == 0x8100000011, f"got {got:#x}, want 0x8000000011"
    got = main.St.read()
    assert got == 0x8100000011, f"got {got:#x}, want 0x8000000011"

    print("\n\ntesting mask update_clear()")
    main.Mask.set()
    print("\n\ntesting update_clear one bit")
    main.Mask.update_clear([39])
    got = main.Mask.read()
    assert got == 0x7fffffffff, f"got {got:#x}, want 0x7fffffffff"
    got = main.St.read()
    assert got == 0x7fffffffff, f"got {got:#x}, want 0x7fffffffff"
    print("\n\ntesting update_clear two bits")
    main.Mask.update_clear([8, 32])
    got = main.Mask.read()
    assert got == 0x7efffffeff, f"got {got:#x}, want 0x7efffffeff"
    got = main.St.read()
    assert got == 0x7efffffeff, f"got {got:#x}, want 0x7efffffeff"

    iface.end(0)

except Exception as E:
    iface.end(1)
    print(traceback.format_exc())
