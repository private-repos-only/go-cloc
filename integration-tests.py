import os
import platform
import subprocess
import sys
import shutil

GITHUB_ORGANIZATION = os.getenv('GO_CLOC_GITHUB_ORGANIZATION')
GITHUB_ACCESS_TOKEN = os.getenv('GO_CLOC_GITHUB_ACCESS_TOKEN')
AZURE_DEVOPS_ORGANIZATION = os.getenv('GO_CLOC_AZURE_DEVOPS_ORGANIZATION')
AZURE_DEVOPS_ACCESS_TOKEN = os.getenv('GO_CLOC_AZURE_DEVOPS_ACCESS_TOKEN')
GITLAB_ORGANIZATION = os.getenv('GO_CLOC_GITLAB_ORGANIZATION')
GITLAB_ACCESS_TOKEN = os.getenv('GO_CLOC_GITLAB_ACCESS_TOKEN')
BITBUCKET_ORGANIZATION = os.getenv('GO_CLOC_BITBUCKET_ORGANIZATION')
BITBUCKET_ACCESS_TOKEN = os.getenv('GO_CLOC_BITBUCKET_ACCESS_TOKEN')

def execute_go_cloc(args):
    os_name = platform.system()

    # Set the binary path based on the operating system
    if os_name == "Linux" or os_name == "Darwin":
        binary_name = "go-cloc"
    elif os_name == "Windows":
        binary_name = "go-cloc.exe"
    else:
        print(f"Unsupported OS: {os_name}")
        sys.exit(1)
    
    # Verify that the binary is in the PATH
    binary_path = shutil.which(binary_name)
    if binary_path is None:
        print(f"Error: {binary_name} not found in PATH.")
        sys.exit(1)
    
    # Collect all command-line arguments passed to the script
    args = [binary_path] + args

     # Run the binary with the provided arguments and capture the output
    try:
        process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        last_line = ""
        for line in process.stdout:
            print(line, end='')  # Print each line to standard output
            last_line = line.strip()  # Keep track of the last line

        process.wait()
        if process.returncode != 0:
            print(f"Error: Process exited with code {process.returncode}")
            sys.exit(process.returncode)

        # Parse the desired value from the last line
        if last_line.isdigit():
            totalLoc = int(last_line)
            return totalLoc
        else:
            print("Expected output not found in the last line")
            return None
    except subprocess.CalledProcessError as e:
        print(f"Error: {e.stderr}")
        sys.exit(e.returncode)
    
    
def run_test(name,args,expected):
    print(f"--------Running test: {name}---------")
    result = execute_go_cloc(args)
    did_pass = (result == expected)
    return {
        "name": name,
        "did_pass": did_pass,
        "expected": expected,
        "actual": result
    }

def print_test_results(test_results):
    for test in test_results:
        print(f"Test: {test['name']}")
        print(f"Expected: {test['expected']}")
        print(f"Actual: {test['actual']}")
        print(f"Pass: {test['did_pass']}")
        print("")
    did_all_pass = all(test['did_pass'] for test in test_results)
    return did_all_pass

if __name__ == "__main__":

    # Run the tests
    test_results = []
    test_results.append(
        run_test(name="GitHub", expected=143933,args=["--devops","GitHub","--organization",GITHUB_ORGANIZATION,"--accessToken",GITHUB_ACCESS_TOKEN,"--log-level","DEBUG","--dump-csvs=false"])
    )
    test_results.append(
        run_test(name="AzureDevOps", expected=57888,args=["--devops","AzureDevOps","--organization",AZURE_DEVOPS_ORGANIZATION,"--accessToken",AZURE_DEVOPS_ACCESS_TOKEN,"--log-level","DEBUG","--dump-csvs=false"])
    )
    test_results.append(
        run_test(name="GitLab", expected=162,args=["--devops","GitLab","--organization",GITLAB_ORGANIZATION,"--accessToken",GITLAB_ACCESS_TOKEN,"--log-level","DEBUG","--dump-csvs=false"])
    )
    test_results.append(
        run_test(name="Bitbucket", expected=4317,args=["--devops","Bitbucket","--organization",BITBUCKET_ORGANIZATION,"--accessToken",BITBUCKET_ACCESS_TOKEN,"--log-level","DEBUG","--dump-csvs=false"])
    )

    did_all_pass = print_test_results(test_results)
    if did_all_pass:
        print("All tests passed!")
    else:
        print("Some tests failed. See above for details")
        sys.exit(-1)