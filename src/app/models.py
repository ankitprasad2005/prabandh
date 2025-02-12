from sqlalchemy import Column, String, Integer, DateTime, ForeignKey, UniqueConstraint
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship
from uuid import uuid4

Base = declarative_base()

class File(Base):
    __tablename__ = 'file'

    id = Column(String, primary_key=True, default=lambda: str(uuid4()))
    file_name = Column(String, nullable=False)
    file_path = Column(String, unique=True, nullable=False)
    extension = Column(String, nullable=False)
    created = Column(DateTime, nullable=False)
    modified = Column(DateTime, nullable=False)
    size = Column(Integer, nullable=False)
    hash = Column(String, nullable=False)
    keywords = relationship("Keyword", back_populates="file")

    # __table_args__ = (
    #     UniqueConstraint('hash', 'size', 'extension', name='_hash_size_extension_uc'),
    # )

class Keyword(Base):
    __tablename__ = 'keyword'

    id = Column(String, primary_key=True, default=lambda: str(uuid4()))
    keyword = Column(String, nullable=False)
    file_id = Column(String, ForeignKey('file.id'), nullable=False)
    file = relationship("File", back_populates="keywords")