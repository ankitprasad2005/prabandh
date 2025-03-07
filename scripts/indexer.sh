#!/bin/bash
# Default to current directory if no argument provided
SEARCH_DIR="${1:-.}"
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"
IGNORE_FILE="$SCRIPT_DIR/.indexignore"
CACHE_FILE="$HOME/.cache/prabandh_cache.csv"

mkdir -p "$(dirname "$CACHE_FILE")"
if [ ! -f "$CACHE_FILE" ]; then
    touch "$CACHE_FILE"
fi

# Function to display usage
usage() {
    echo "Usage: $(basename "$0") [directory]"
    echo "If directory is not specified, current directory will be used"
    exit 1
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

# Check if help is requested
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
fi

# Check if specified directory exists
if [ ! -d "$SEARCH_DIR" ]; then
    echo "Error: Directory '$SEARCH_DIR' does not exist"
    exit 1
fi

# Create or truncate cache file with header
echo "DIRECTORY_PATH,FILE_NAME,EXTENSION,CREATED_DATE,MODIFIED_DATE,SIZE_BYTES,SHA256_HASH" > "$CACHE_FILE"

# Build the find command prune conditions
if [ -f "$IGNORE_FILE" ]; then
    PRUNE_EXPR=""
    while IFS= read -r pattern || [ -n "$pattern" ]; do
        [[ -z "$pattern" || "$pattern" =~ ^[[:space:]]*# ]] && continue
        pattern="${pattern%"${pattern##*[![:space:]]}"}"
        pattern="${pattern#"${pattern%%[![:space:]]*}"}"
        pattern="${pattern%/}"
        if [ -z "$PRUNE_EXPR" ]; then
            PRUNE_EXPR="-name $pattern"
        else
            PRUNE_EXPR="$PRUNE_EXPR -o -name $pattern"
        fi
    done < "$IGNORE_FILE"

    # Process files and collect metadata
    if [ -n "$PRUNE_EXPR" ]; then
        find "$SEARCH_DIR" -type d \( $PRUNE_EXPR \) -prune -o -type f -print0 | 
        while IFS= read -r -d '' file; do
            file_components=$(get_file_components "$file")
            IFS='|' read -r dirpath filename extension <<< "$file_components"
            
            created_date=$(stat -c %W "$file" 2>/dev/null || stat -c %Y "$file")
            modified_date=$(stat -c %Y "$file")
            size=$(stat -c %s "$file")
            sha256=$(get_sha256 "$file")
            
            created_date=$(format_date "$created_date")
            modified_date=$(format_date "$modified_date")
            
            echo "$dirpath,$filename,$extension,$created_date,$modified_date,$size,$sha256" >> "$CACHE_FILE"
        done
    else
        find "$SEARCH_DIR" -type f -print0 |
        while IFS= read -r -d '' file; do
            file_components=$(get_file_components "$file")
            IFS='|' read -r dirpath filename extension <<< "$file_components"
            
            created_date=$(stat -c %W "$file" 2>/dev/null || stat -c %Y "$file")
            modified_date=$(stat -c %Y "$file")
            size=$(stat -c %s "$file")
            sha256=$(get_sha256 "$file")
            
            created_date=$(format_date "$created_date")
            modified_date=$(format_date "$modified_date")
            
            echo "$dirpath,$filename,$extension,$created_date,$modified_date,$size,$sha256" >> "$CACHE_FILE"
        done
    fi
else
    echo "Warning: indexer.ignore file not found at $IGNORE_FILE"
    echo "Proceeding without ignore patterns..."
    
    find "$SEARCH_DIR" -type f -print0 |
    while IFS= read -r -d '' file; do
        file_components=$(get_file_components "$file")
        IFS='|' read -r dirpath filename extension <<< "$file_components"
        
        created_date=$(stat -c %W "$file" 2>/dev/null || stat -c %Y "$file")
        modified_date=$(stat -c %Y "$file")
        size=$(stat -c %s "$file")
        sha256=$(get_sha256 "$file")
        
        created_date=$(format_date "$created_date")
        modified_date=$(format_date "$modified_date")
        
        echo "$dirpath,$filename,$extension,$created_date,$modified_date,$size,$sha256" >> "$CACHE_FILE"
    done
fi

echo "Index cache has been created at: $CACHE_FILE"