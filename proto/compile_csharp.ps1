$sourceDir = "./src"
$artifactsDir = "./.artifacts"
$compilerVersion = "protoc-21.12-win64"

$compilerPath = "${compilerVersion}\bin\protoc.exe"
$pluginPath = "${compilerVersion}\bin\grpc_csharp_plugin.exe"
$includePath = "${compilerVersion}\include"
$fileExtension = ".g.cs"

function LogCreateFile {
    param (
        $inputFile,
        $outputFile
    )

    $inputFile = $inputFile | Resolve-Path -Relative
    $outputFile = $outputFile | Resolve-Path -Relative

    if ($PSVersionTable.PSVersion.Major -ge 7) {
        "$inputFile", $outputFile `
        | Select-Object @{Name = 'String'; Expression = { $_ } } `
        | Format-Wide String -Column 2 | Out-String -NoNewLine
    }
    else {
        "$inputFile", $outputFile -join " ---> "
    }
}

function ParseProto {
    param (
        $protoFile
    )

    $data = (Get-Content $protoFile.FullName) -join "`n"
    
    if ($data -notmatch "option \(csharp_assembly\) = `"([\w\.]+)`";") {
        Write-Error "Assembly name not found in $protoFile"
        return
    }

    $assemblyName = $Matches[1]
    $assemblyNameSegments = $assemblyName -split "\." 

    if ($data -notmatch "option csharp_namespace = `"([\w\.]+)`";") {
        Write-Error "Namespace not found in $protoFile"
        return
    }

    $namespaceSegments = $Matches[1] -split "\."

    $pathSegments = @()
    $pathSegments += $assemblyName

    for ($i = 0; $i -lt $namespaceSegments.Length; $i++) {
        if ($namespaceSegments[$i] -eq $assemblyNameSegments[$i]) {
            continue
        }
        else {
            # add remaining
            $pathSegments += ($namespaceSegments[$i..$namespaceSegments.Length] -join ".")
            break
        }
    }

    $directory = $pathSegments -join "\"
    $publicAccess = $data -match "option \(csharp_access\) = public;"
    $fileName = (Get-Culture).TextInfo.ToTitleCase(($protoFile.Name -replace "_", " ")) -replace " ", "" -replace ".proto", $fileExtension
    
    return $fileName, $directory, $publicAccess
}


[Environment]::CurrentDirectory = $pwd
$artifacts = New-Object -Type System.IO.DirectoryInfo -ArgumentList $artifactsDir

if ($artifacts.Exists) {
    $artifacts.Delete($true)
}

[System.IO.Directory]::CreateDirectory($artifactsDir) | Out-Null

Get-ChildItem -Path $sourceDir -Recurse -Filter "*.proto" | ForEach-Object {
   
    $protoFile = Resolve-Path -Relative $_
    $fileName, $outputDir, $public = ParseProto $_

    $csharpOptions = @()
    $gcrpOptions = @()

    $csharpOptions += "file_extension=$fileExtension"
    $gcrpOptions += "no_server"
    if (-not $public) {
        $csharpOptions += "internal_access"
        $gcrpOptions += "internal_access"
    }

    $outputDir = "$artifactsDir\$outputDir"
    [System.IO.Directory]::CreateDirectory($outputDir) | Out-Null

    $csharpOpt = $csharpOptions -join ','
    $gcrpOpt = $gcrpOptions -join ','

    & $compilerPath `
        --plugin=protoc-gen-grpc=$pluginPath `
        --proto_path=$sourceDir `
        --proto_path=$includePath `
        --grpc_out=$outputDir `
        --grpc_opt=$gcrpOpt `
        --csharp_out=$outputDir `
        --csharp_opt=$csharpOpt `
        $protoFile

    LogCreateFile $protoFile "$outputDir\$fileName"
                    
}
