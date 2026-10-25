# SPDX-FileCopyrightText: none
#
# SPDX-License-Identifier: CC0-1.0

# NOTE(bouncepaw): Ported from Betula with some changes.

import os, re, subprocess, sys

ROOT: str = os.path.dirname(os.path.abspath(__file__))

def guess_version() -> str:
    with open(os.path.join(ROOT, 'help', 'en.myco'), encoding='utf-8') as f:
        content = f.read()
    # Extract version.
    m = re.search(r'This is documentation for Mycorrhiza Wiki (\d+\.\d+(?:\.\d+)*(?:-rc\d+)?)\. ', content)
    if m is None:
        sys.exit(
            "Could not find the version in help/en.myco. "
            "Make sure it contains 'This is documentation for Mycorrhiza Wiki 1.15.1. '."
        )
    return 'v' + m.group(1)

def pick_version() -> str:
    result = subprocess.run(
        ['git', 'tag', '--sort=-creatordate'],
        capture_output=True, text=True, cwd=ROOT, check=True
    )
    tags = [t for t in result.stdout.strip().splitlines() if t]
    if not tags:
        sys.exit('No tags found to compare to.')
    for i, tag in enumerate(tags):
        print(f'  {i + 1}. {tag}')
    while True:
        try:
            n = int(input('Enter number: ').strip())
            if 1 <= n <= len(tags):
                return tags[n - 1]
        except ValueError:
            pass

def shortlog(version: str) -> str:
    result = subprocess.run(
        ['git', 'shortlog', f'{version}..HEAD'],
        capture_output=True, text=True, cwd=ROOT, check=True
    )
    return result.stdout.strip()

if __name__ == '__main__':
    print('genlog.py will help you prepare a Mycorrhiza release.')

    while True:
        print('\nStep 1. New version')
        guessed = guess_version()
        answer = input(f'Is it {guessed}? [Enter to agree, or type a version]: ').strip()
        new_version = guessed if answer == '' else answer
        if not new_version.startswith('v'):
            new_version = 'v' + new_version
        release = re.sub(r'-rc\d+$', '', new_version[1:])

        print('\nStep 2. Pick old version to compare to')
        old_version = pick_version()

        print('\nStep 3. The message')
        optional_text = input('Optional text (empty to skip): ').strip()
        log = shortlog(old_version)

        parts = [f'Mycorrhiza {new_version}']
        if optional_text:
            parts.append('')
            parts.append(optional_text)
        parts.append('')
        parts.append('Release notes:')
        parts.append('')
        parts.append(f'=> https://mycorrhiza.wiki/hypha/release/{release}')
        parts.append('')
        parts.append(log)
        msg = '\n'.join(parts)

        print()
        print(msg)

        print('\nStep 4. Tag creation')
        action = input('Create tag? [y]es / [r]estart / [q]uit: ').strip().lower()
        if action == 'y':
            subprocess.run(['git', 'tag', '-a', new_version, '-m', msg], cwd=ROOT, check=True)
            print(f'Tag {new_version} created.')
            break
        elif action == 'r':
            continue
        else:
            sys.exit(0)
