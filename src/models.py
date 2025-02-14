from dataclasses import dataclass
from datetime import datetime
from uuid import UUID, uuid4


@dataclass
class File:
    id: UUID
    name: str
    path: str
    extension: str
    created: datetime
    updated: datetime
    size: int
    hash: str


@dataclass
class Keyword:
    id: UUID
    file_id: str
    keyword: str
