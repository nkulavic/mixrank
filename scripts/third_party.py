#!/usr/bin/env python3
"""Collect license notices from the pinned Go module cache for binary releases."""
import json,pathlib,subprocess
root=pathlib.Path(__file__).resolve().parents[1]
data=subprocess.check_output(['go','list','-m','-json','all'],cwd=root,text=True,encoding="utf-8");dec=json.JSONDecoder();modules=[]
while data.strip():
 data=data.lstrip();obj,n=dec.raw_decode(data);modules.append(obj);data=data[n:]
parts=['Third-party license notices for the MixRank native binaries.\nOriginal MixRank code is MIT licensed; provider data rights are separate.\n']
for m in modules:
 if m.get('Main'):continue
 
 if 'Dir' not in m:
  downloaded=json.loads(subprocess.check_output(['go','mod','download','-json',m['Path']+'@'+m['Version']],text=True,encoding="utf-8"));m['Dir']=downloaded['Dir']
 d=pathlib.Path(m['Dir']);license_files=[p for p in d.iterdir() if p.is_file() and p.name.upper().startswith(('LICENSE','COPYING','NOTICE'))]
 if not license_files:raise SystemExit('No license found for '+m['Path'])
 parts.append('\n===== '+m['Path']+' '+m.get('Version','')+' =====\n')
 for p in sorted(license_files):parts.append(p.name+'\n'+p.read_text(encoding="utf-8"))
goroot=pathlib.Path(subprocess.check_output(['go','env','GOROOT'],text=True,encoding="utf-8").strip());parts.append('\n===== Go runtime =====\n'+(goroot/'LICENSE').read_text(encoding="utf-8"))
(root/'THIRD_PARTY_NOTICES.txt').write_text('\n'.join(parts)+'\n',encoding='utf-8',newline='\n')
