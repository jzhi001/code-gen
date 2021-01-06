from dataclasses import dataclass, fields
from enum import Enum
from typing import List

FieldType = Enum('string', 'bool', 'int32', 'int64')

@dataclass
class Field:
    """single protobuf field"""
    id: int
    name: str
    type: str
    repeated: bool
    required: bool

    def __repr__(self) -> str:
        f = '{type} {name} = {id};'.format(self.type, self.name, self.id)
        if self.repeated:
            f = 'repeated ' + f
        return f

@dataclass
class Message:
    """Data structure for a Protobuf message"""
    name: str
    fields: List[Field]

    def __repr__(self) -> str:
        fdesc = str.join('\n', map(lambda f: str(f), fields))
        return 'message {name}{{\n{fdesc}}}'.format(self.name, fdesc)