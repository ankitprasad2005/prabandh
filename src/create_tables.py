from app.database import engine, Base
from app.models import File, Keyword

# Create all tables in the database
Base.metadata.create_all(bind=engine)