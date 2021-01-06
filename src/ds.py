
from types.general import Token
from typing import Callable

class TokenIterator:

    def __init__(self, tokens: list[Token], isWs: Callable[[str], bool]) -> None:
        self.tokens = tokens
        self.i = 0
        # TODO TokenIterator is not responsible for this
        self.isWs = isWs

    def hasNext(self) -> bool:
        return self.i < len(self.tokens)

    def next(self) -> Token:
        tk = self.tokens[self.i]
        self.i += 1
        return tk
    
    def nextNonWs(self) -> Token:
        tk =  None
        while self.hasNext() and (tk is None or self.isWs(tk)):
            tk = self.next()
        if tk is None:
            raise RuntimeError('no more non-ws charactor')
        return tk

class CharIterator:
    def __init__(self, string: str) -> None:
        self.s = string
        self.n = len(string)
        self.i = 0

    def hasNext(self) -> bool:
        return self.i < self.n

    def next(self) -> str:
        if not self.hasNext():
            raise RuntimeError('no more character!')
        c = self.s[self.i]
        self.i += 1
        return c

    def jumpTo(self, tar: str) -> None:
        if not self.hasNext():
            raise RuntimeError('no more character!')
        self.i = self.s.find(tar, self.i)
        if self.i < 0:
            raise RuntimeError('target character not found!')
        self.i += 1

    def skipUntil(self, tar: str) -> None:
        """
            i = nextPosition(tar) - 1
        """
        self.jumpTo(tar)
        self.i -= 1

class BufferedStrList:
    def __init__(self) -> None:
        self.items = []
        self.cur = ''

    def appendBuffer(self, s: str) -> None:
        self.cur += s

    def flush(self) -> None:
        if len(self.cur) == 0:
            return
        self.items.append(self.cur)
        self.cur = ''

    def getItems(self) -> list:
        return self.items.copy()

    def appendItem(self, item: str) -> None:
        if len(self.cur) > 0:
            raise RuntimeError('current buffer is not empty') 
        self.items.append(item)