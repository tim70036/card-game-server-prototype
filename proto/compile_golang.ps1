$compilerVersion = "protoc-21.12-win64"
$sourceDir = ".\src"
$artifactsDir = "..\pkg\grpc"
$compilerPath = "${compilerVersion}\bin\protoc.exe"
$includePath = "${compilerVersion}\include"

[Environment]::CurrentDirectory = $pwd
$artifacts = New-Object -Type System.IO.DirectoryInfo -ArgumentList $artifactsDir

if ($artifacts.Exists) {
    $artifacts.Delete($true)
}

[System.IO.Directory]::CreateDirectory($artifactsDir) | Out-Null

Get-ChildItem -Path $sourceDir -Recurse -Filter "*.proto" | ForEach-Object {

    $protoFile = Resolve-Path -Relative $_

    $protoFile
    & $compilerPath `
        --proto_path=$sourceDir `
        --proto_path=$includePath `
        --go_out=$artifactsDir `
        --go_opt=paths=import,module=card-game-server-prototype/pkg/grpc `
        --go-grpc_out=$artifactsDir `
        --go-grpc_opt=paths=import,module=card-game-server-prototype/pkg/grpc `
        $protoFile
                    
}
