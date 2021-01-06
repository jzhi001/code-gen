
import logging
import re
import argparse
from typing import Iterator, Tuple

# type aliases
GoTypeName = str
GoFieldType = str
GoFieldName = str
Token = str
GoField = Tuple[GoFieldName, GoFieldType]
GoType = Tuple[GoTypeName, list[GoField]]

argParser = argparse.ArgumentParser()

logging.basicConfig(level=logging.INFO, format='%(asctime)s %(message)s')

protobufTypes = ['string', 'int32', 'int64', 'bool']

class TokenIterator:
    def __init__(self, tokens: list[Token]) -> None:
        self.tokens = tokens
        self.i = 0

    def hasNext(self) -> bool:
        return self.i < len(self.tokens)

    def next(self) -> Token:
        tk = self.tokens[self.i]
        self.i += 1
        return tk
    
    def nextNonWs(self) -> Token:
        tk =  None
        while self.hasNext() and (tk is None or tk == '\n'):
            tk = self.next()
        if tk is None:
            raise RuntimeError('no more non-ws charactor')
        return tk

def isWs(c: str):
    return c.isspace() and c != '\n'

def camel2Snake(s: str) -> str:
    ans = ''
    buf = ''

    # TODO refactor with CharIterator

    for c in s:
        if c.isupper():
            if len(buf) > 0:
                if len(ans) > 0:
                    ans += '_'
                ans += buf
                buf = ''
        buf += c.lower()
    if len(buf) > 0:
        if len(ans) > 0:
            ans += '_'
        ans += buf
    return ans

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
        logging.debug('jumpTo(%s) %s', re.escape(tar), re.escape(self.s[self.i: self.i + 20]))
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

def parseTypeDec(tkIter: TokenIterator) -> Tuple[GoFieldName, GoFieldType]:
    name = tkIter.nextNonWs()
    value = tkIter.nextNonWs()
    return name, value

def parse(tokens: list[Token]) -> list[GoType]:
    
    tkIter = TokenIterator(tokens)

    goTypes: list[GoType] = []

    while tkIter.hasNext():
        tk = tkIter.next()

        if tk == '\n':
            continue

        if tk != 'type':
            raise RuntimeError('expect "type" keyword')

        typeName = tkIter.nextNonWs()
        tkIter.nextNonWs() # skip 'struct'
        tkIter.nextNonWs() # skip '{'

        fields: list[GoField] = []

        fieldName = tkIter.nextNonWs()
        while fieldName != '}':
            fieldType = tkIter.nextNonWs()
            fields.append((fieldName, fieldType))
            fieldName = tkIter.nextNonWs()

        goTypes.append((typeName, fields))

    return goTypes

# struct dec := type <name> struct {[field-dec]}
# field dec := <name> <type> [`<desc>`] [//<comment>] \n
def tokenize(structDec: str) -> list[Token]:

    class TokenBuffer:
        def __init__(self) -> None:
            self.tokens = []
            self.cur = ''

        def add(self, char: str) -> None:
            self.cur += char

        def flush(self) -> None:
            if len(self.cur) == 0:
                return
            self.tokens.append(self.cur)
            self.cur = ''

        def tokenList(self) -> list:
            return self.tokens.copy()

        def addToken(self, tk: str) -> None:
            if len(self.cur) > 0:
               raise RuntimeError('current buffer is not empty') 
            self.tokens.append(tk)

    tkBuf = TokenBuffer()
    cit = CharIterator(structDec)

    while cit.hasNext():
        c = cit.next()
        if isWs(c):
            tkBuf.flush()
        elif c == '/':
            tkBuf.flush()
            cit.skipUntil('\n')
        elif c == '`':
            tkBuf.flush()
            cit.jumpTo('`')
        elif c == '{' or c == '}' or c == '\n':
            tkBuf.flush()
            tkBuf.addToken(c)
        else:
            tkBuf.add(c)

    return tkBuf.tokenList()

def protobuf(goType: GoType) -> str:
    
    structName, fields = goType

    ans = 'message ' + structName + '{\n'

    fieldDecList = []

    for i, field in enumerate(fields):
        fname, ftype = field
        fdec = '\t'
        if ftype.startswith('[]'):
            fdec += 'repeated '
            ftype = ftype[2:]
        if ftype.startswith('*'):
            ftype = ftype[1:]
        fdec += ftype + ' ' + camel2Snake(fname) + ' = ' + str(i + 1) + ';'
        fieldDecList.append(fdec)

    ans += str.join('\n', fieldDecList)

    ans += '\n}\n'

    return ans

if __name__ == "__main__":
    with open('./type.txt', 'r') as f:
        structDec = f.read()
        tokens = tokenize(structDec)
        print(tokens)
        goTypes = parse(tokens)

        with open('./result.proto', 'w') as resFile:
            for gt in goTypes:
                resFile.write(protobuf(gt))
        