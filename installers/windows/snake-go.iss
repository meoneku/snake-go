; Inno Setup Script untuk Snake Go
; Kompilasi dengan Inno Setup Compiler (ISCC) untuk menghasilkan SnakeGoSetup.exe

#define MyAppName "Snake Go"
#define MyAppVersion "1.0.0"
#define MyAppPublisher "Antigravity Games"
#define MyAppExeName "snake-go-gui.exe"

[Setup]
AppId={{C786DF19-53BF-453D-9492-4AEB1D9A7500}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\SnakeGo
DisableProgramGroupPage=yes
OutputDir=..\..\dist
OutputBaseFilename=SnakeGoSetup-{#MyAppVersion}
SetupIconFile=..\..\gui\icon.ico
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=lowest

[Languages]
Name: "indonesian"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "..\..\dist\snake-go-gui-windows-amd64.exe"; DestDir: "{app}"; DestName: "{#MyAppExeName}"; Flags: ignoreversion
Source: "..\..\dist\snake-go-terminal-windows-amd64.exe"; DestDir: "{app}"; DestName: "snake-go.exe"; Flags: ignoreversion; Tasks: 
Source: "..\..\gui\icon.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\VERSION"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\icon.ico"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\icon.ico"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent
