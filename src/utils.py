from .models import File, Keyword

from sqlite3 import Connection
from dataclasses import asdict


# Singular creation of a new file
def insert_file(database: Connection, file: File) -> None:
    query = """
        INSERT INTO file (id, name, path, extension, created, update, size, hash)
        VALUES (:id, :name, :path, :extension, :created, :update, :size, :hash)    
    """

    database.execute(query, asdict(file))
    database.commit()


# Singular updating of a file
def upsert_file(database: Connection, file: File) -> None:
    query = """
        INSERT INTO file (id, name, path, extension, created, update, size, hash)
        VALUES (:id, :name, :path, :extension, :created, :update, :size, :hash)    
        ON CONFLICT(id) DO UPDATE SET 
            name = excluded.name,
            path = excluded.path,
            extension = excluded.extension,
            created = excluded.created,
            update = excluded.update,
            size = excluded.size,
            hash = excluded.hash
    """

    database.execute(query, asdict(file))
    database.commit()


# Batch uploading of files
def insert_batch_files(database: Connection, files: list[File]) -> None:
    query = """
        INSERT INTO file (id, name, path, extension, created, update, size, hash)
        VALUES (:id, :name, :path, :extension, :created, :update, :size, :hash)
    """

    database.executemany(query, [asdict(file) for file in files])
    database.commit()


# Batch update and insertion of files
def upsert_batch_files(database: Connection, files: list[File]) -> None:
    query = """
        INSERT INTO file (id, name, path, extension, created, update, size, hash)
        VALUES (:id, :name, :path, :extension, :created, :update, :size, :hash)
        ON CONFLICT(id) DO UPDATE SET 
            name = excluded.name,
            path = excluded.path,
            extension = excluded.extension,
            created = excluded.created,
            update = excluded.update,
            size = excluded.size,
            hash = excluded.hash
    """

    database.executemany(query, [asdict(file) for file in files])
    database.commit()
