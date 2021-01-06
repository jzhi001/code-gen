from typing import Tuple

# type aliases
GoTypeName = str
GoFieldType = str
GoFieldName = str
GoField = Tuple[GoFieldName, GoFieldType]
GoType = Tuple[GoTypeName, list[GoField]]

mapping = {'int64': 'int64'}