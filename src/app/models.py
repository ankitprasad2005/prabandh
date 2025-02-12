from pydantic import BaseModel, Field
from typing import List
from datetime import datetime
from uuid import UUID, uuid4

class Keyword(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    keyword: str
    file_id: UUID

    class Config:
        orm_mode = True

class File(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    file_name: str
    file_path: str
    extension: str
    created: datetime
    modified: datetime
    size: int
    hash: str
    keywords: List[Keyword] = []

    class Config:
        orm_mode = True