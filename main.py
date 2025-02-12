#!/usr/bin/env python3

import argparse
import os
import subprocess
import sys
from pathlib import Path

def run_indexer(directory: str = None) -> None:
    """
    Run the indexer.sh script with the specified directory.
    If no directory is specified, uses the current directory.
    
    Args:
        directory (str, optional): Directory path to index. Defaults to None.
    """
    # Get the directory where this Python script is located
    script_dir = Path(__file__).parent.absolute()
    
    # Path to the indexer.sh script (expecting it to be in the same directory)
    indexer_path = script_dir / "indexer.sh"
    
    # Check if indexer.sh exists
    if not indexer_path.exists():
        print(f"Error: Could not find indexer.sh at {indexer_path}", file=sys.stderr)
        sys.exit(1)
    
    # Make sure indexer.sh is executable
    try:
        indexer_path.chmod(indexer_path.stat().st_mode | 0o111)
    except Exception as e:
        print(f"Error: Could not make indexer.sh executable: {e}", file=sys.stderr)
        sys.exit(1)
    
    # Build the command
    cmd = [str(indexer_path)]
    if directory:
        # Convert to absolute path
        abs_directory = str(Path(directory).absolute())
        cmd.append(abs_directory)
    
    # Run the indexer script
    try:
        process = subprocess.run(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Print stdout
        if process.stdout:
            print(process.stdout, end='')
        
        # Print stderr to stderr
        if process.stderr:
            print(process.stderr, end='', file=sys.stderr)
        
        # Check return code
        if process.returncode != 0:
            print(f"Error: indexer.sh exited with code {process.returncode}", file=sys.stderr)
            sys.exit(process.returncode)
            
    except subprocess.SubprocessError as e:
        print(f"Error running indexer.sh: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error: {e}", file=sys.stderr)
        sys.exit(1)

def main():
    # Create argument parser
    parser = argparse.ArgumentParser(
        description="Run file indexer on specified directory",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s                   # Index current directory
  %(prog)s /path/to/dir      # Index specified directory
  %(prog)s -h                # Show this help message
        """
    )
    
    # Add arguments
    parser.add_argument(
        'directory',
        nargs='?',
        help='Directory to index (defaults to current directory)',
        default=None
    )
    
    # Parse arguments
    args = parser.parse_args()
    
    # Run indexer with specified directory
    run_indexer(args.directory)

if __name__ == "__main__":
    main()