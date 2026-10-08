#ifndef AppVersion
  #define AppVersion "0.0.0"
#endif
#ifndef Arch
  #define Arch "amd64"
#endif
#ifndef BuildDir
  #define BuildDir "dist\installer-stage\amd64"
#endif

#if Arch == "arm64"
  #define AllowedArch "arm64"
#else
  #define AllowedArch "x64compatible and not arm64"
#endif

[Setup]
AppId=andmetoo.dellkvm
AppName=dellkvm
AppVersion={#AppVersion}
AppPublisher=andmetoo
AppPublisherURL=https://github.com/andmetoo/dellkvm
AppSupportURL=https://github.com/andmetoo/dellkvm/issues
DefaultDirName={localappdata}\Programs\dellkvm
DefaultGroupName=dellkvm
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
ArchitecturesAllowed={#AllowedArch}
OutputBaseFilename=dellkvm_{#AppVersion}_windows_{#Arch}_setup
Compression=lzma2
SolidCompression=yes
UninstallDisplayIcon={app}\dellkvm-tray.exe
CloseApplications=yes

[Files]
Source: "{#BuildDir}\dellkvm.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#BuildDir}\dellkvm-tray.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#BuildDir}\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#BuildDir}\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#BuildDir}\config.toml.default"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: "autostart"; Description: "Start dellkvm when I sign in"; GroupDescription: "Startup:"

[Icons]
Name: "{group}\dellkvm"; Filename: "{app}\dellkvm-tray.exe"; WorkingDir: "{app}"
Name: "{userstartup}\dellkvm"; Filename: "{app}\dellkvm-tray.exe"; WorkingDir: "{app}"; Tasks: autostart

[Run]
Filename: "{app}\dellkvm-tray.exe"; Description: "Run dellkvm now"; Flags: postinstall nowait skipifsilent unchecked
