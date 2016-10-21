$ErrorActionPreference = 'Stop';

$packageName= 'unicreds' # arbitrary name for the package, used in messages
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = '{{DownloadUrl}}' # download url for 32-bit
$url64      = 'https://github.com/Versent/unicreds/releases/download/1.5.1/unicreds_1.5.1_windows_amd64.tar.gz' #

Install-ChocolateyZipPackage -PackageName $packageName `
  -Url64 $url64 `
  -UnzipLocation $toolsDir

$filename = [System.IO.Path]::GetFileNameWithoutExtension($(Split-Path $url64 -Leaf))

Get-ChocolateyUnzip `
  -FileFullPath "$toolsDir\$filename" `
  -Destination $toolsDir
