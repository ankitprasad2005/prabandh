import csv
from pathlib import Path
from sqlalchemy.orm import Session
from sqlalchemy.dialects.postgresql import insert
from app.database import get_db, engine
from app.models import File as FileModel, Keyword as KeywordModel
from fastapi.encoders import jsonable_encoder
from datetime import datetime

def bulk_upsert_from_csv(db: Session, csv_file_path: Path):
    if not csv_file_path.exists():
        raise FileNotFoundError(f"CSV file not found at {csv_file_path}")

    with open(csv_file_path, mode='r') as file:
        reader = csv.DictReader(file)
        for row in reader:
            file_data = FileModel(
                file_name=row["FILE_NAME"],
                file_path=row["DIRECTORY_PATH"],
                extension=row["EXTENSION"],
                created=datetime.strptime(row["CREATED_DATE"], "%Y-%m-%d %H:%M:%S"),
                modified=datetime.strptime(row["MODIFIED_DATE"], "%Y-%m-%d %H:%M:%S"),
                size=int(row["SIZE_BYTES"]),
                hash=row["SHA256_HASH"]
            )

            file_dict = jsonable_encoder(file_data)
            stmt = insert(FileModel).values(file_dict).on_conflict_do_update(
                index_elements=['file_path'],
                set_=file_dict
            )
            db.execute(stmt)
        db.commit()

def main():
    csv_file_path = Path.home() / ".cache" / "prabandh_cache.csv"
    db = next(get_db())
    try:
        bulk_upsert_from_csv(db, csv_file_path)
        print("Bulk upsert completed successfully.")
    except Exception as e:
        db.rollback()
        print(f"Error during bulk upsert: {e}")
    finally:
        db.close()

if __name__ == "__main__":
    main()