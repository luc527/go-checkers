import sys
import re

passed = True
packages_below = []

for line in sys.stdin:
    line = line.rstrip()

    match_cov = re.search(r'coverage: (\d+\.\d+)%', line)
    if match_cov is None:
        continue

    match_pkg = re.search(r'github.com/luc527/go_checkers/([a-zA-Z0-9_]+)', line)
    if match_pkg is None:
        print('Failed to match with package name')
        exit(1)

    percentage = float(match_cov.group(1))
    package    = match_pkg.group(1)

    print(f'Package {package} has {percentage:.1f}% coverage')

    if percentage < 80:
        passed = False

print()
if not passed:
    print(r'Quality gate failed: all packages need >=80% coverage')
    exit(1)
else:
    print('Quality gate passed')
    exit(0)
