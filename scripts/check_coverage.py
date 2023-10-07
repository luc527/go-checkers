import sys
import re

for line in sys.stdin:
    line = line.rstrip()
    match = re.search(r'coverage: (\d+\.\d+)%', line)
    if match is None:
        continue
    percentage = float(match.group(1))
    if percentage < 80:
        print(r'Quality gate failed: all packages need >=80% coverage')
        exit(1)

exit(0)