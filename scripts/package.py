#!/usr/bin/env python3
"""Generate synchronized agent packages; --release adds native binaries and archives."""
import argparse,hashlib,json,os,pathlib,shutil,subprocess,tarfile,zipfile
root=pathlib.Path(__file__).resolve().parents[1];version='0.1.0'
ap=argparse.ArgumentParser();ap.add_argument('--release',action='store_true');ap.add_argument('--reuse-binaries',action='store_true',help='Repackage already-built binaries without recompiling');ap.add_argument('--check',action='store_true');args=ap.parse_args()
def write(path,data):
 p=root/path;p.parent.mkdir(parents=True,exist_ok=True);b=json.dumps(data,indent=2)+'\n'
 if args.check:
  if not p.exists() or p.read_text(encoding="utf-8")!=b:raise SystemExit('Stale package file: '+str(path))
 else:p.write_text(b,encoding="utf-8",newline="\n")
base={'name':'mixrank','version':version,'description':'MixRank API tools, CLI recipes and evidence-based research workflows','author':{'name':'Nick Kulavic','url':'https://github.com/nkulavic'},'homepage':'https://github.com/nkulavic/mixrank','repository':'https://github.com/nkulavic/mixrank','license':'MIT'}
codex={k:v for k,v in base.items() if k not in ['homepage','repository','license']};codex.update({'skills':'./skills/','mcpServers':'./.mcp.json','interface':{'displayName':'MixRank','shortDescription':'Research companies and people with MixRank.','longDescription':'Company and person discovery, enrichment, validation, technology research and private product context through a Go CLI and MCP server.','developerName':'Nick Kulavic','category':'Productivity','capabilities':[],'defaultPrompt':'Research companies and people with MixRank using my criteria.'}})
claude={**base,'skills':'./skills','mcpServers':'./.mcp.json','userConfig':{'api_key':{'type':'string','title':'MixRank API key','description':'Your MixRank API credential. Claude manages secure storage.','sensitive':True,'required':True}}}
# Source and GitHub packages use a pinned launcher; release bundles also include binaries.
cmcp={'mcpServers':{'mixrank':{'command':'${CLAUDE_PLUGIN_ROOT}/scripts/run-mixrank','args':['mcp'],'env':{'MIXRANK_API_KEY':'${user_config.api_key}'}}}}
xmcp={'mcpServers':{'mixrank':{'command':'./scripts/run-mixrank','cwd':'.','args':['mcp'],'env_vars':['MIXRANK_API_KEY']}}}
write('.codex-plugin/plugin.json',codex);write('.mcp.json',xmcp)
portable={**base,'$schema':'https://agent-plugins.org/schemas/1.0.0/plugin.schema.json','extensions':{'com.openai':{'interface':codex['interface']}}}
pmcp={'$schema':'https://agent-plugins.org/schemas/1.0.0/mcp.schema.json','mcpServers':{'mixrank':{'type':'stdio','command':'./scripts/run-mixrank','args':['mcp']}}}
write('plugin.json',portable);write('mcp.json',pmcp)
write('plugins/mixrank/plugin.json',portable);write('plugins/mixrank/mcp.json',pmcp)
write('.claude-plugin/plugin.json',claude)
write('.claude-plugin/marketplace.json',{'name':'mixrank-toolkit','description':'Native MixRank tools and research skills for coding agents','owner':{'name':'Nick Kulavic'},'plugins':[{'name':'mixrank','source':'./plugins/claude/mixrank','version':version,'description':base['description']}]})
write('.agents/plugins/marketplace.json',{'name':'mixrank-toolkit','interface':{'displayName':'MixRank Toolkit'},'plugins':[{'name':'mixrank','source':{'source':'local','path':'./plugins/mixrank'},'policy':{'installation':'AVAILABLE','authentication':'ON_INSTALL'},'category':'Productivity'}]})
for path,manifest,mcp in [('plugins/mixrank',codex,xmcp),('plugins/claude/mixrank',claude,cmcp)]:
 sub=root/path
 write(path+('/.codex-plugin/plugin.json' if path=='plugins/mixrank' else '/.claude-plugin/plugin.json'),manifest);write(path+'/.mcp.json',mcp)
 for source in list((root/'skills').rglob('*'))+[root/'scripts/run-mixrank',root/'scripts/run-mixrank.ps1']:
  if not source.is_file():continue
  target=sub/source.relative_to(root)
  if args.check:
   if not target.exists() or target.read_bytes()!=source.read_bytes():raise SystemExit('Stale packaged content: '+str(target))
  else:target.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(source,target)
if not args.release:raise SystemExit(0)
dist=root/'dist';dist.mkdir(exist_ok=True)
for goos in ['darwin','linux','windows']:
 for arch in ['amd64','arm64']:
  ext='.exe' if goos=='windows' else '';folder=dist/f'mixrank_{version}_{goos}_{arch}';folder.mkdir(exist_ok=True);binary=folder/('mixrank'+ext)
  if not args.reuse_binaries:subprocess.run(['go','build','-trimpath','-ldflags','-s -w -buildid=','-o',str(binary),'./cmd/mixrank'],cwd=root,env={**os.environ,'GOOS':goos,'GOARCH':arch,'CGO_ENABLED':'0'},check=True)
  if goos=='windows':
   with zipfile.ZipFile(str(folder)+'.zip','w',zipfile.ZIP_DEFLATED) as z:
    z.write(binary,binary.name)
    for notice in ['LICENSE','THIRD_PARTY_NOTICES.txt']:z.write(root/notice,notice)
  else:
   with tarfile.open(str(folder)+'.tar.gz','w:gz') as t:
    t.add(binary,arcname='mixrank')
    for notice in ['LICENSE','THIRD_PARTY_NOTICES.txt']:t.add(root/notice,arcname=notice)
  # Native Desktop extension; one architecture per package avoids host ambiguity.
  bundle=dist/f'mcpb-{goos}-{arch}';bundle.mkdir(exist_ok=True);shutil.copy2(binary,bundle/binary.name)
  for notice in ['LICENSE','THIRD_PARTY_NOTICES.txt']:shutil.copy2(root/notice,bundle/notice)
  mf={'manifest_version':'0.3','name':'mixrank','display_name':'MixRank','version':version,'description':base['description'],'author':base['author'],'homepage':base['homepage'],'license':'MIT','server':{'type':'binary','entry_point':binary.name,'mcp_config':{'command':'${__dirname}/'+binary.name,'args':['mcp'],'env':{'MIXRANK_API_KEY':'${user_config.api_key}'}}},'compatibility':{'platforms':['win32' if goos=='windows' else goos]},'user_config':{'api_key':{'type':'string','title':'MixRank API key','description':'Your API key, stored by Claude Desktop','sensitive':True,'required':True}}}
  (bundle/'manifest.json').write_text(json.dumps(mf,indent=2)+'\n',encoding='utf-8',newline='\n')
  with zipfile.ZipFile(dist/f'mixrank_{goos}_{arch}.mcpb','w',zipfile.ZIP_DEFLATED) as z:
   for f in bundle.iterdir():z.write(f,f.name)
# Cache-independent universal plugin bundles.
for label,source in [('codex',root/'plugins/mixrank'),('claude',root/'plugins/claude/mixrank')]:
 with zipfile.ZipFile(dist/f'mixrank-{label}.zip','w',zipfile.ZIP_DEFLATED) as z:
  for notice in ['LICENSE','THIRD_PARTY_NOTICES.txt']:z.write(root/notice,notice)
  for f in source.rglob('*'):
   if f.is_file() and 'bin' not in f.relative_to(source).parts:z.write(f,f.relative_to(source))
  for goos in ['darwin','linux','windows']:
   for arch in ['amd64','arm64']:
    name='mixrank'+('.exe' if goos=='windows' else '');z.write(dist/f'mixrank_{version}_{goos}_{arch}'/name,f'bin/{goos}_{arch}/{name}')
# Direct-native plugin bundles for every supported OS/architecture.
for label,source in [('codex',root/'plugins/mixrank'),('claude',root/'plugins/claude/mixrank')]:
 for goos in ['darwin','linux','windows']:
  for arch in ['amd64','arm64']:
   name='mixrank'+('.exe' if goos=='windows' else '')
   variable='.' if label=='codex' else '${CLAUDE_PLUGIN_ROOT}'
   config={'mcpServers':{'mixrank':{'command':variable+'/bin/'+name,'args':['mcp']}}}
   if label=='claude':config['mcpServers']['mixrank']['env']={'MIXRANK_API_KEY':'${user_config.api_key}'}
   else:config['mcpServers']['mixrank'].update({'env_vars':['MIXRANK_API_KEY'],'cwd':'.'})
   with zipfile.ZipFile(dist/f'mixrank-{label}-{goos}-{arch}.zip','w',zipfile.ZIP_DEFLATED) as z:
    for notice in ['LICENSE','THIRD_PARTY_NOTICES.txt']:z.write(root/notice,notice)
    for f in source.rglob('*'):
     if f.is_file() and f.name not in ['.mcp.json','mcp.json'] and 'bin' not in f.relative_to(source).parts:z.write(f,f.relative_to(source))
    z.writestr('.mcp.json',json.dumps(config,indent=2));
    if label=='codex':z.writestr('mcp.json',json.dumps({'$schema':pmcp['$schema'],'mcpServers':{'mixrank':{'type':'stdio','command':'./bin/'+name,'args':['mcp']}}},indent=2))
    z.write(dist/f'mixrank_{version}_{goos}_{arch}'/name,'bin/'+name)

# Cowork gets portable skills, not a falsely host-bound Desktop MCP launcher.
with zipfile.ZipFile(dist/'mixrank-cowork.zip','w',zipfile.ZIP_DEFLATED) as z:
 z.write(root/'LICENSE','LICENSE')
 manifest={k:v for k,v in base.items()};manifest['skills']='./skills';z.writestr('.claude-plugin/plugin.json',json.dumps(manifest,indent=2))
 for f in (root/'skills').rglob('*'):
  if f.is_file():z.write(f,f.relative_to(root))
 z.writestr('CONNECTOR.md','Import this plugin in Cowork. Research requires the hosted MixRank connector once OAuth deployment is complete. Configure it in Cowork; host keychains and local Desktop MCP configuration are not available in the VM. General research/personalization skills can operate with supplied files meanwhile.\n')
for name in ['install.sh','install.ps1']:shutil.copy2(root/name,dist/name)
assets=sorted(p for p in dist.iterdir() if p.is_file() and p.name!='checksums.txt')
(dist/'checksums.txt').write_text(''.join(hashlib.sha256(p.read_bytes()).hexdigest()+'  '+p.name+'\n' for p in assets),encoding='utf-8',newline='\n')
print('Release artifacts:',len(assets))
