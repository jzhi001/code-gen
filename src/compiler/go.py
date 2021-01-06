import strconv
from typing import Tuple
from types.general import Token
from ds import BufferedStrList, CharIterator, TokenIterator
from types.go import GoField, GoFieldName, GoFieldType, GoType, GoTypeName

def isWs(c: str):
    return c.isspace() and c != '\n'

# struct dec := type <name> struct {[field-dec]}
# field dec := <name> <type> [`<desc>`] [//<comment>] \n
def tokenize(structDec: str) -> list[Token]:

    """
        tokenize one or more Go type declarations
    """

    tkBuf: BufferedStrList = BufferedStrList()
    cit: CharIterator = CharIterator(structDec)

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

    return tkBuf.getItems()

def parse(tokens: list[Token]) -> list[GoType]:
    
    tkIter = TokenIterator(tokens, isWs)

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

def parseTypeDec(tkIter: TokenIterator) -> Tuple[GoFieldName, GoFieldType]:
    name = tkIter.nextNonWs()
    value = tkIter.nextNonWs()
    return name, value


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
        fdec += ftype + ' ' + strconv.camel2Snake(fname) + ' = ' + str(i + 1) + ';'
        fieldDecList.append(fdec)

    ans += str.join('\n', fieldDecList)

    ans += '\n}\n'

    return ans

if __name__ == "__main__":
    with open('../test/type.txt', 'r') as f:
        structDec = f.read()
        tokens = tokenize(structDec)
        print(tokens)
        goTypes = parse(tokens)

        with open('./result.proto', 'w') as resFile:
            for gt in goTypes:
                resFile.write(protobuf(gt))