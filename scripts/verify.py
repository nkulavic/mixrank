#!/usr/bin/env python3
"""Offline invariants for catalog, skills and public release source."""
import json,pathlib,re,sys
root=pathlib.Path(__file__).resolve().parents[1];cat=json.loads((root/'catalog/endpoints.json').read_text(encoding="utf-8"))
ops=cat['operations'];assert len(ops)==102
assert len({(o['method'],o['path']) for o in ops})==len(ops)
methods=(root/'operations_gen.go').read_text(encoding="utf-8")
for op in ops:
 assert op['id'] in methods,op['id']
 assert op['sources'] and op['verified']
for skill in (root/'skills').iterdir():
 if not skill.is_dir():continue
 text=(skill/'SKILL.md').read_text(encoding="utf-8");assert text.startswith('---\nname: '+skill.name+'\n');assert '\ndescription:' in text
 for target in re.findall(r'\]\(([^)]+)\)',text):
  if '://' not in target:assert (skill/target).exists(),target
# Check secret-like API URLs, not random hashes or synthetic fixtures.
for p in root.rglob('*'):
 if not p.is_file() or any(x in p.parts for x in ['.git','dist','bin']):continue
 try:text=p.read_text(encoding="utf-8")
 except UnicodeDecodeError:continue
 assert not re.search(r'/v2/json/[a-f0-9]{32}(?:/|["\s])',text,re.I),'credential URL in '+str(p)
 assert ('REDACTED_'+'API_KEY') not in text,'raw documentation snapshot in '+str(p)
print('Catalog, generated SDK and portable skill invariants passed')
