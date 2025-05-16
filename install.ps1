$ErrorActionPreference = 'Stop'

$repo = 'valcinei/jiboia-tunnel'
$binary = 'jiboia'
$latest = (Invoke-RestMethod https://api.github.com/repos/$repo/releases/latest).tag_name
$arch = 'amd64'
$url = "https://github.com/$repo/releases/download/$latest/${binary}-windows-$arch.zip"

Write-Host "Downloading $url..."
Invoke-WebRequest -Uri $url -OutFile "$binary.zip"

Expand-Archive "$binary.zip" -DestinationPath . -Force
Rename-Item -Path "${binary}-windows-$arch" -NewName "$binary.exe"

$installPath = "$Env:ProgramFiles\JiboiaTunnel"
New-Item -ItemType Directory -Path $installPath -Force | Out-Null
Move-Item "$binary.exe" "$installPath\$binary.exe" -Force

$env:Path += ";$installPath"
Write-Host "✅ $binary installed to $installPath"
& "$installPath\$binary.exe" --help
