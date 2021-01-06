import ds

def camelToSnake(s: str) -> str:

    bufList = ds.BufferedStrList()
    iter = ds.CharIterator(s)

    while iter.hasNext():
        c = iter.next()

        if c.isupper():
            bufList.flush()
        bufList.appendBuffer(c.lower())
    bufList.flush()

    return str.join('_', bufList.getItems())

