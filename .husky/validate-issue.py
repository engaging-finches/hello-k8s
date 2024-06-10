#!/usr/bin/env python3
import os
import re
import requests
import sys
import argparse
from dotenv import load_dotenv

# Determine the path to the .env file in the parent directory
parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
env_path = os.path.join(parent_dir, '.env')

# Load environment variables from .env file
load_dotenv(dotenv_path=env_path)


# Configuration
# THIS NEEDS TO BE PRESENT IN A .env FILE IN ROOT DIR
GITHUB_TOKEN = os.getenv("PAT_TOKEN")

REPO_OWNER = "engaging-finches"
REPO_NAME = "hello-k8s"
ISSUE_REGEX = r"#(\d+)"

def get_commit_message():
    parser = argparse.ArgumentParser()
    parser.add_argument('--message', help='Commit message')
    args = parser.parse_args()
    return args.message

def is_valid_issue(issue_number):

    print(f"In is_valid_issue()")
    query = """
    query {
        repository(owner: "%s", name: "%s") {
            issue(number: %d) {
                number
            }
        }
    }
    """ % (REPO_OWNER, REPO_NAME, issue_number)

    headers = {
        "Authorization": f"Bearer {GITHUB_TOKEN}",
        "Content-Type": "application/json"
    }

    response = requests.post("https://api.github.com/graphql", json={"query": query}, headers=headers)
    print(response)
    if response.status_code == 200:
        data = response.json()
        print(data)
        # Check for errors in the response
        if 'errors' in data and data['errors']:
            # If there are errors, print the first error message and return False
            print(f"Error!!!: {data['errors'][0]['message']}")
            return False

        if "data" in data and "repository" in data["data"] and "issue" in data["data"]["repository"]:
            return True
    return False

def main():
    # print(f"In validate-issue main()")
    commit_message = get_commit_message()
    # print(f"Commit message: {commit_message}")
    # prevent commit from going through for testing
    match = re.search(ISSUE_REGEX, commit_message)

    if match:
        issue_number = int(match.group(1))
        print("\nISSUE NUMBER:", issue_number ,"\n")
        if is_valid_issue(issue_number):
            print(f"Valid issue reference found: #{issue_number}")
            sys.exit(0)
        else:
            print(f"Error: Invalid issue reference: #{issue_number}")
            sys.exit(1)
    else:
        print("Error: No valid issue reference found in commit message.")
        sys.exit(1)

if __name__ == "__main__":
    main()
