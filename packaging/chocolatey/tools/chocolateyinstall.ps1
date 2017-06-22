$ErrorActionPreference = 'Stop';

$packageName= 'unicreds'
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64      = 'https://github.com/Versent/unicreds/releases/download/1.5.1/unicreds_1.5.1_windows_amd64.tar.gz'

Install-ChocolateyZipPackage -PackageName $packageName `
  -Url64 $url64 `
  -UnzipLocation $toolsDir `
  -ChecksumType 'md5' `
  -Checksum64 '75D8B2496675A2AFDC41F865472AB18A'

$filename = [System.IO.Path]::GetFileNameWithoutExtension($(Split-Path $url64 -Leaf))

Get-ChocolateyUnzip `
  -FileFullPath "$toolsDir\$filename" `
  -Destination $toolsDir
