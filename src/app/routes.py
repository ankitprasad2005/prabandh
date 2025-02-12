from fastapi import APIRouter, HTTPException, Depends
from prisma import Prisma
from typing import List
import hashlib
from .models import File as FileModel, Keyword as KeywordModel

router = APIRouter()
prisma = Prisma()

@router.on_event("startup")
async def startup():
    await prisma.connect()

@router.on_event("shutdown")
async def shutdown():
    await prisma.disconnect()

@router.post("/index")
async def add_entry(file: FileModel):
    try:
        file_hash = hashlib.sha256(file.file_path.encode()).hexdigest()
        if file_hash != file.hash:
            raise HTTPException(status_code=400, detail="Hash mismatch")
        
        new_file = await prisma.file.create(
            data={
                "file_name": file.file_name,
                "file_path": file.file_path,
                "extension": file.extension,
                "created": file.created,
                "modified": file.modified,
                "size": file.size,
                "hash": file.hash,
                "keywords": {
                    "create": [{"keyword": keyword.keyword} for keyword in file.keywords]
                }
            }
        )
        return {"message": "File added successfully", "file": new_file}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/index/all")
async def add_bulk_entries(entries: List[FileModel]):
    try:
        new_entries = []
        for entry in entries:
            new_entry = await prisma.file.create(
                data={
                    "file_name": entry.file_name,
                    "file_path": entry.file_path,
                    "extension": entry.extension,
                    "created": entry.created,
                    "modified": entry.modified,
                    "size": entry.size,
                    "hash": entry.hash,
                    "keywords": {
                        "create": [{"keyword": keyword.keyword} for keyword in entry.keywords]
                    }
                }
            )
            new_entries.append(new_entry)
        return {"message": "Bulk entries added successfully", "entries": new_entries}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/search")
async def search_entries(query: str):
    try:
        entries = await prisma.file.find_many(
            where={
                "file_path": {
                    "contains": query
                }
            },
            include={
                "keywords": True
            }
        )
        return entries
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))