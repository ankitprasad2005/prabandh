from fastapi import APIRouter, HTTPException, Depends
from sqlalchemy.orm import Session
from typing import List
import hashlib
from app.models import File as FileModel, Keyword as KeywordModel
from .database import get_db
from .models import File, Keyword
import csv
from pathlib import Path

router = APIRouter()

@router.post("/index")
async def add_entry(file: FileModel, db: Session = Depends(get_db)):
    try:
        file_hash = hashlib.sha256(file.file_path.encode()).hexdigest()
        if file_hash != file.hash:
            raise HTTPException(status_code=400, detail="Hash mismatch")
        
        new_file = File(
            file_name=file.file_name,
            file_path=file.file_path,
            extension=file.extension,
            created=file.created,
            modified=file.modified,
            size=file.size,
            hash=file.hash,
            keywords=[Keyword(keyword=keyword.keyword) for keyword in file.keywords]
        )
        db.add(new_file)
        db.commit()
        db.refresh(new_file)
        return {"message": "File added successfully", "file": new_file}
    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/index/all")
async def add_bulk_entries_from_csv(db: Session = Depends(get_db)):
    try:
        csv_file_path = Path.home() / ".cache" / "prabandh_cache.csv"
        if not csv_file_path.exists():
            raise HTTPException(status_code=404, detail="CSV file not found")

        new_entries = []
        with open(csv_file_path, mode='r') as file:
            reader = csv.DictReader(file)
            for row in reader:
                new_entry = File(
                    file_name=row["FILE_NAME"],
                    file_path=row["DIRECTORY_PATH"],
                    extension=row["EXTENSION"],
                    created=row["CREATED_DATE"],
                    modified=row["MODIFIED_DATE"],
                    size=int(row["SIZE_BYTES"]),
                    hash=row["SHA256_HASH"],
                    keywords=[]  # Assuming no keywords in CSV
                )
                db.add(new_entry)
                new_entries.append(new_entry)
        db.commit()
        return {"message": "Bulk entries added successfully", "entries": new_entries}
    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/search")
async def search_entries(query: str, db: Session = Depends(get_db)):
    try:
        entries = db.query(File).filter(File.file_path.contains(query)).all()
        return entries
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))