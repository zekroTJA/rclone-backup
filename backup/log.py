__FMT_ERROR = "38;5;196"
__FMT_HIGHLIGHT = "38;5;87"
__FMT_SUCCESS = "38;5;10"

no_colors = False


def disable_colors():
    global no_colors
    no_colors = True


def __format(v, fmt) -> str:
    if no_colors:
        return v
    return f"\x1b[{fmt}m{v}\x1b[0m"


def error(v):
    print(f"{__format('ERROR:', __FMT_ERROR)} {v}")


def highlight(v):
    print(__format(v, __FMT_HIGHLIGHT))


def success(v):
    print(__format(v, __FMT_SUCCESS))


def fail(v):
    print(__format(v, __FMT_ERROR))
