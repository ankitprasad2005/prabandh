import sqlite3
from sqlite3 import Connection


class Database:
    def __init__(self):
        self.database = sqlite3.connect("../local_database.db")
        self.cursor = self.database.cursor()

        # Enable foreign key support
        self.cursor.execute("PRAGMA foreign_keys = ON;")

        # Create database tables
        self.cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS file (
                id TEXT PRIMARY KEY,
                name TEXT NOT NULL,
                path TEXT NOT NULL,
                extension TEXT NOT NULL,
                created DATETIME NOT NULL, 
                update DATETIME NOT NULL,
                size INT NOT NULL,
                hash TEXT
            ) 
            """
        )

        self.cursor.execute(
            """
            CREATE TABLE IF NOT EXISTS keyword (
                id TEXT PRIMARY KEY,
                file_id TEXT NOT NULL,
                keyword TEXT NOT NULL,
                FOREIGN KEY (file_id) REFERENCES file(id) ON DELETE CASCADE
            ) 
            """
        )

        self.database.commit()
        self.database.execute

    def get_database(self) -> Connection:
        return self.database
