$ErrorActionPreference = 'Stop'
$Version = '0.2.0'
$Arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLower()
if ($Arch -eq 'x64') { $Arch = 'amd64' }
if ($Arch -notin @('amd64','arm64')) { throw 'Unsupported architecture' }
$Temp = Join-Path ([System.IO.Path]::GetTempPath()) ([guid]::NewGuid().ToString())
New-Item -ItemType Directory $Temp | Out-Null
try {
  $Asset = "mixrank_${Version}_windows_$Arch.zip"
  $Base = "https://github.com/nkulavic/mixrank/releases/download/v$Version"
  Invoke-WebRequest "$Base/$Asset" -OutFile (Join-Path $Temp $Asset)
  Invoke-WebRequest "$Base/checksums.txt" -OutFile (Join-Path $Temp 'checksums.txt')
  $Line = Get-Content (Join-Path $Temp 'checksums.txt') | Where-Object { ($_ -split '\s+')[1] -eq $Asset }
  if (@($Line).Count -ne 1) { throw 'Missing release checksum' }
  $Expected = ($Line -split '\s+')[0]
  if ((Get-FileHash (Join-Path $Temp $Asset) -Algorithm SHA256).Hash -ne $Expected) { throw 'Release checksum mismatch' }
  Expand-Archive (Join-Path $Temp $Asset) -DestinationPath $Temp
  & (Join-Path $Temp 'mixrank.exe') setup @args
  $Result = $LASTEXITCODE
} finally { Remove-Item -Recurse -Force $Temp }
exit $Result
