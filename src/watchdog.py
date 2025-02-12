import time
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
                    file_data = FileModel(
                        file_name=row["FILE_NAME"],
                        file_path=row["DIRECTORY_PATH"],
                        extension=row["EXTENSION"],
                        created=datetime.strptime(row["CREATED_DATE"], "%Y-%m-%d %H:%M:%S"),
                        modified=datetime.strptime(row["MODIFIED_DATE"], "%Y-%m-%d %H:%M:%S"),
                        size=int(row["SIZE_BYTES"]),
                        hash=row["SHA256_HASH"],
                        keywords=[]  # Assuming no keywords in CSV
                    )

                    file_dict = jsonable_encoder(file_data)
                    stmt = insert(FileModel).values(file_dict).on_conflict_do_update(
                        index_elements=['file_path'],
                        set_=file_dict
                    )
                    self.db.execute(stmt)
                self.db.commit()
                print("CSV processed and database updated.")
        except Exception as e:
            self.db.rollback()
            print(f"Error processing CSV: {e}")

def main():
    csv_file_path = Path.home() / ".cache" / "prabandh_cache.csv"
    db = next(get_db())
    event_handler = CSVEventHandler(db)
    observer = Observer()
    observer.schedule(event_handler, path=str(csv_file_path.parent), recursive=False)
    observer.start()
    print(f"Watching for changes in {csv_file_path}")

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
    observer.join()
    db.close()

if __name__ == "__main__":
    main()