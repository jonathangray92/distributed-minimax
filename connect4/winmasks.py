# produce a one-hot bitmask at position (row, col) 
def spot(row, col):
    return 1 << (row * 7 + col)

def all_win_masks():
    # horiziontal lines
    for r in xrange(6):
        for c in xrange(4):
            yield spot(r, c) | spot(r, c + 1) | spot(r, c + 2) | spot(r, c + 3)

    # vertical lines
    for r in xrange(3):
        for c in xrange(7):
            yield spot(r, c) | spot(r + 1, c) | spot(r + 2, c) | spot(r + 3, c)

    # diagonal lines (top-left to bottom-right)
    for r in xrange(3):
        for c in xrange(4):
            yield spot(r, c) | spot(r + 1, c + 1) | spot(r + 2, c + 2) | spot(r + 3, c + 3)

    # diagonal lines (top-right to bottom-left)
    for r in xrange(3):
        for c in xrange(6, 2, -1):
            yield spot(r, c) | spot(r + 1, c - 1) | spot(r + 2, c - 2) | spot(r + 3, c - 3)


print ''.join(str(mask)+',' for mask in all_win_masks())
