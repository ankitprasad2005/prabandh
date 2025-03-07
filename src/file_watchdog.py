import time
import sys
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
from pathlib import Path
from sqlalchemy.orm import Session
from sqlalchemy.dialects.postgresql import insert
from app.database import get_db
from app.models import File as FileModel, Keyword as KeywordModel
from fastapi.encoders import jsonable_encoder
from datetime import datetime
import csv

class CSVEventHandler(FileSystemEventHandler):
    def __init__(self, db: Session):
        self.db = db

    def on_modified(self, event):
        if event.src_path.endswith("prabandh_cache.csv"):
            self.process_csv(event.src_path)

    def process_csv(self, csv_file_path):
        try:
            with open(csv_file_path, mode='r') as file:
                reader = csv.DictReader(file)
                for row in reader:
                    self.upsert_data(row)
                print("CSV processed and database updated.")
        except Exception as e:
            self.db.rollback()
            print(f"Error processing CSV: {e}")

    def upsert_data(self, data: dict):
        try:
            file_data = FileModel(
                file_name=data["FILE_NAME"],
                file_path=data["DIRECTORY_PATH"],
                extension=data["EXTENSION"],
                created=datetime.strptime(data["CREATED_DATE"], "%Y-%m-%d %H:%M:%S"),
                modified=datetime.strptime(data["MODIFIED_DATE"], "%Y-%m-%d %H:%M:%S"),
                size=int(data["SIZE_BYTES"]),
                hash=data["SHA256_HASH"]
            )

            file_dict = jsonable_encoder(file_data)
            stmt = insert(FileModel).values(file_dict).on_conflict_do_update(
                index_elements=['file_path'],
                set_=file_dict
            )
            self.db.execute(stmt)
            self.db.commit()
            print("Data upserted successfully.")
        except Exception as e:
            self.db.rollback()
            print(f"Error during upsert: {e}")

def main():
    if len(sys.argv) != 8:
        print("Usage: python watchdog.py <DIRECTORY_PATH> <FILE_NAME> <EXTENSION> <CREATED_DATE> <MODIFIED_DATE> <SIZE_BYTES> <SHA256_HASH>")
        sys.exit(1)

    data = {
        "DIRECTORY_PATH": sys.argv[1],
        "FILE_NAME": sys.argv[2],
        "EXTENSION": sys.argv[3],
        "CREATED_DATE": sys.argv[4],
        "MODIFIED_DATE": sys.argv[5],
        "SIZE_BYTES": sys.argv[6],
        "SHA256_HASH": sys.argv[7]
    }

    db = next(get_db())
    try:
        event_handler = CSVEventHandler(db)
        event_handler.upsert_data(data)
    except Exception as e:
        print(f"Error: {e}")
    finally:
        db.close()

if __name__ == "__main__":
    main()