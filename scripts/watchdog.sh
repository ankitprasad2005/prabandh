#!/bin/bash

# Set the API endpoint
API_URL="http://localhost:8000/api/update"  # Replace with your actual API endpoint

# Load ignore patterns from .indexignore
IGNORE_FILE="$(dirname "$(readlink -f "$0")")/.indexignore"
IGNORE_PATTERNS=()
if [ -f "$IGNORE_FILE" ]; then
    while IFS= read -r line; do
        # Skip empty lines and comments
        [[ "$line" =~ ^#.*$ ]] && continue
        [[ -z "$line" ]] && continue
        IGNORE_PATTERNS+=("$line")
    done < "$IGNORE_FILE"
fi

# Function to check if a file should be ignored
should_ignore() {
    local file="$1"
    for pattern in "${IGNORE_PATTERNS[@]}"; do
        if [[ "$file" =~ $pattern ]]; then
            return 0
        fi
    done
    return 1
}

# Function to get file components
get_file_components() {
    local fullpath="$1"
    local dirpath=$(dirname "$fullpath")
    local fullname=$(basename "$fullpath")
    local filename="${fullname%.*}"
    local ext="${fullname##*.}"
    if [ "$ext" = "$fullname" ]; then
        ext="no_extension"
        filename="$fullname"
    fi
    echo "$dirpath|$filename|$ext"
}

# Function to format date
format_date() {
    date -d "@$1" "+%Y-%m-%d %H:%M:%S"
}

# Function to calculate SHA256 hash
get_sha256() {
    sha256sum "$1" | cut -d' ' -f1
}

# Function to send API call
send_api_call() {
    local dirpath="$1"
    local filename="$2"
    local extension="$3"
    local created_date="$4"
    local modified_date="$5"
    local size="$6"
    local sha256="$7"

    curl -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d '{
            "DIRECTORY_PATH": "'"$dirpath"'",
            "FILE_NAME": "'"$filename"'",
            "EXTENSION": "'"$extension"'",
            "CREATED_DATE": "'"$created_date"'",
            "MODIFIED_DATE": "'"$modified_date"'",
            "SIZE_BYTES": '"$size"',
            "SHA256_HASH": "'"$sha256"'"
        }'
}

# Function to process file
process_file() {
    local file="$1"
    if should_ignore "$file"; then
        echo "Ignoring $file"
        return
    fi

    file_components=$(get_file_components "$file")
    IFS='|' read -r dirpath filename extension <<< "$file_components"
    
    created_date=$(stat -c %W "$file" 2>/dev/null || stat -c %Y "$file")
    modified_date=$(stat -c %Y "$file")
    size=$(stat -c %s "$file")
    sha256=$(get_sha256 "$file")
    
    created_date=$(format_date "$created_date")
    modified_date=$(format_date "$modified_date")
    
    echo "$dirpath|$filename|$extension|$created_date|$modified_date|$size|$sha256"
    send_api_call "$dirpath" "$filename" "$extension" "$created_date" "$modified_date" "$size" "$sha256"
}

# Monitor the home directory for changes
inotifywait -m -r -e create -e moved_to --format '%w%f' "$HOME" | while read file; do
    process_file "$file"
done