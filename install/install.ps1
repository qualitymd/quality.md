param(
  [string]$Version = $env:QUALITYMD_VERSION,
  [string]$InstallDir = $env:QUALITYMD_HOME,
  [switch]$NonInteractive
)

$ErrorActionPreference = "Stop"
if ([string]::IsNullOrWhiteSpace($Version)) { $Version = "latest" }
if ([string]::IsNullOrWhiteSpace($InstallDir)) { $InstallDir = Join-Path $HOME ".qualitymd" }

$repo = "qualitymd/quality.md"
$binDir = Join-Path $InstallDir "bin"
if ($Version -eq "latest") {
  $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -Headers @{ "User-Agent" = "qualitymd-installer" }
  $Version = $release.tag_name
}
if ([string]::IsNullOrWhiteSpace($Version)) { throw "could not resolve qualitymd version" }

$arch = if ([Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture -eq "Arm64") { "arm64" } else { "amd64" }
$archive = "qualitymd_$($Version.TrimStart('v'))_windows_$arch.zip"
$baseUrl = "https://github.com/$repo/releases/download/$Version"
$tmp = Join-Path ([IO.Path]::GetTempPath()) ("qualitymd-install-" + [Guid]::NewGuid())
$stage = Join-Path (Join-Path $InstallDir "releases") $Version
New-Item -ItemType Directory -Force -Path $tmp, $stage, $binDir | Out-Null

try {
  $archivePath = Join-Path $tmp $archive
  Invoke-WebRequest -Uri "$baseUrl/$archive" -OutFile $archivePath
  $checksumsPath = Join-Path $tmp "checksums.txt"
  try {
    Invoke-WebRequest -Uri "$baseUrl/checksums.txt" -OutFile $checksumsPath
    $expected = (Select-String -Path $checksumsPath -Pattern " $([regex]::Escape($archive))$" | Select-Object -First 1).Line.Split(" ")[0]
    if ($expected) {
      $actual = (Get-FileHash -Algorithm SHA256 $archivePath).Hash.ToLowerInvariant()
      if ($expected.ToLowerInvariant() -ne $actual) { throw "checksum mismatch for $archive" }
    }
  } catch {
    if ($_.Exception.Message -like "checksum mismatch*") { throw }
  }

  Expand-Archive -Force -Path $archivePath -DestinationPath $stage
  $binary = Get-ChildItem -Path $stage -Filter "qualitymd.exe" -Recurse | Select-Object -First 1
  if (-not $binary) { throw "archive did not contain qualitymd.exe" }
  Copy-Item -Force $binary.FullName (Join-Path $binDir "qualitymd.exe")
  @(
    "layoutVersion=1",
    "version=$Version",
    "channel=github"
  ) | Set-Content -Encoding ascii (Join-Path $InstallDir ".qualitymd-managed-install")

  & (Join-Path $binDir "qualitymd.exe") --version | Out-Null
  Write-Output "Installed qualitymd $Version to $(Join-Path $binDir "qualitymd.exe")"
} finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
