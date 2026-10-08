//go:build linux || windows

package main

// localeData is generated from cmd/dellkvm/locales/active.*.toml.
var localeData = map[string]string{
	"en": `[UsageDetect]
other = "Usage: dellkvm detect"

[UsageCurrent]
other = "Usage: dellkvm current"

[UsageSwitch]
other = "Usage: dellkvm switch <id>"

[UsageRoot]
other = "Usage: dellkvm [tray|tui|init|detect|current|switch <id>|status --json|help]"

[Help]
other = "Usage: dellkvm [command]\n\nCommands:\n  init         create a starter config if missing\n  tray         open the system tray\n  tui          open the terminal UI\n  status --json  print desktop status as JSON\n  detect       list available monitors\n  current      show the current input\n  switch <id>  switch to a configured input\n  help         show this help\n\nRun without a command to open the system tray. Use dellkvm tui for the terminal UI.\nConfig is optional; run dellkvm init to customize inputs. Use dellkvm detect to list monitors, or keep bus = 0 for auto-detect."

[ConfigReadFailed]
other = "failed to read {{.Path}}"

[ConfigNotFound]
other = "config not found. Checked paths: {{.Paths}}\nCopy config.toml.default to config.toml or ~/.config/dellkvm/config.toml and edit it.\nRun dellkvm init to create a starter configuration."

[InvalidBus]
other = "config.toml bus must be 0 for auto-detect or a positive number"

[MissingInputs]
other = "config.toml must contain at least one input"

[InputNotFound]
other = "input with id {{.ID}} was not found"

[SwitchSentNoVerify]
other = "Switch command sent to {{.Name}} ({{.ID}}, {{.Code}}) through bus {{.Bus}}, but the monitor did not answer the getvcp {{.VCP}} verification yet: {{.Error}}"

[SwitchVerifiedCode]
other = "Switched to {{.Name}} ({{.ID}}, {{.Code}}) through bus {{.Bus}}. Verification: current code {{.CurrentCode}}"

[SwitchVerifiedInput]
other = "Switched to {{.Name}} ({{.ID}}, {{.Code}}) through bus {{.Bus}}. Verification: {{.CurrentName}} ({{.CurrentCode}})"

[SwitchVerificationMismatch]
other = "Switch command sent to {{.Name}} ({{.ID}}, {{.Requested}}) through bus {{.Bus}}, but verification reported {{.Observed}}."

[AutoDetectedPrefix]
other = "Auto-detected bus {{.Bus}}. {{.Message}}"

[VCPCodeNotFound]
other = "sl=0x.. was not found in getvcp output:\n{{.Raw}}"

[VCPCodeParseFailed]
other = "failed to parse sl=0x{{.Value}}"

[CurrentCode]
other = "Current input: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Current input: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[MissingDDC]
other = "ddcutil not found. Install it with: sudo dnf install ddcutil i2c-tools"

[I2CAccess]
other = "Add your user to the i2c group: sudo usermod -aG i2c $USER, then log in again"

[NoI2CDevices]
other = "I2C devices were not found. Load the module: sudo modprobe i2c-dev. If that does not help, check DDC/CI in the monitor OSD"

[DDCRetry]
other = "monitor was found, but the DDC/VCP request did not get a stable response. This does not look like a sudo problem. Check DDC/CI in the monitor OSD, cable, dock, or port; switching may still work even when reading the current input is unstable"

[DDCTimeout]
other = "ddcutil command timed out after {{.Seconds}}s: ddcutil {{.Command}}"

[DDCTimeoutWithOutput]
other = "ddcutil command timed out after {{.Seconds}}s: ddcutil {{.Command}}\n{{.Raw}}"

[DDCCommandFailedNoOutput]
other = "ddcutil {{.Command}}: {{.Error}}"

[DDCCommandFailedWithOutput]
other = "ddcutil {{.Command}}: {{.Error}}\n{{.Raw}}"

[AutoDetectFailed]
other = "bus auto-detection failed"

[AutoDetectNoDisplays]
other = "bus auto-detection found no displays in ddcutil detect --brief output"

[AutoDetectNoResponsive]
other = "bus auto-detection checked {{.Buses}}, but no bus answered getvcp {{.VCP}}"

[NoMonitorUnknownBus]
other = "monitor not found on the selected bus. Check available displays: dellkvm detect. Then edit bus in config.toml, or keep bus = 0 for auto-detect"

[NoMonitorBus]
other = "monitor not found on bus {{.Bus}}. Check available displays: dellkvm detect. Then edit bus in config.toml, or keep bus = 0 for auto-detect"

[ListTitle]
other = "Inputs"

[Refreshing]
other = "Refreshing current input..."

[Switching]
other = "Switching to {{.Name}}..."

[ErrorPrefix]
other = "Error: {{.Error}}"

[CurrentRefreshed]
other = "Current input refreshed through bus {{.Bus}}."

[CurrentCodeNotFound]
other = "Current code is not in the input list."

[TuiLastSelected]
other = "Last selected: {{.Name}} (bus {{.Bus}}, unverified)"

[TuiCurrentUnknown]
other = "Current input: unknown"

[TuiCurrentInput]
other = "Current input: {{.Name}} ({{.Code}}, bus {{.Bus}})"

[TuiCurrentCode]
other = "Current input: {{.Code}} (bus {{.Bus}})"

[Keys]
other = "Enter: switch  r: refresh  q/Ctrl+C: quit"
`,
	"ru": `[UsageDetect]
other = "Использование: dellkvm detect"

[UsageCurrent]
other = "Использование: dellkvm current"

[UsageSwitch]
other = "Использование: dellkvm switch <id>"

[UsageRoot]
other = "Использование: dellkvm [tray|tui|init|detect|current|switch <id>|status --json|help]"

[Help]
other = "Использование: dellkvm [command]\n\nКоманды:\n  init         создать стартовый конфиг, если его нет\n  tray         открыть приложение в трее\n  tui          открыть терминальный интерфейс\n  status --json  вывести статус в JSON\n  detect       показать доступные мониторы\n  current      показать текущий вход\n  switch <id>  переключить на input из config.toml\n  help         показать эту справку\n\nЗапуск без команды открывает приложение в трее. dellkvm tui открывает терминальный интерфейс.\nКонфиг необязателен; выполни dellkvm init, чтобы настроить входы. Используй dellkvm detect для списка мониторов, или оставь bus = 0 для автоопределения."

[ConfigReadFailed]
other = "не удалось прочитать {{.Path}}"

[ConfigNotFound]
other = "конфиг не найден. Проверены пути: {{.Paths}}\nСкопируй config.toml.default в config.toml или ~/.config/dellkvm/config.toml и отредактируй.\nВыполни dellkvm init, чтобы создать стартовый конфиг."

[InvalidBus]
other = "в config.toml bus должен быть 0 для auto или положительным номером"

[MissingInputs]
other = "в config.toml должен быть хотя бы один input"

[InputNotFound]
other = "input с id {{.ID}} не найден"

[SwitchSentNoVerify]
other = "Команда переключения отправлена на {{.Name}} ({{.ID}}, {{.Code}}) через bus {{.Bus}}, но монитор временно не отвечает на проверку getvcp {{.VCP}}: {{.Error}}"

[SwitchVerifiedCode]
other = "Переключено на {{.Name}} ({{.ID}}, {{.Code}}) через bus {{.Bus}}. Проверка: текущий код {{.CurrentCode}}"

[SwitchVerifiedInput]
other = "Переключено на {{.Name}} ({{.ID}}, {{.Code}}) через bus {{.Bus}}. Проверка: {{.CurrentName}} ({{.CurrentCode}})"

[SwitchVerificationMismatch]
other = "Команда переключения отправлена на {{.Name}} ({{.ID}}, {{.Requested}}) через bus {{.Bus}}, но проверка показала {{.Observed}}."

[AutoDetectedPrefix]
other = "Автоопределен bus {{.Bus}}. {{.Message}}"

[VCPCodeNotFound]
other = "не найден sl=0x.. в выводе getvcp:\n{{.Raw}}"

[VCPCodeParseFailed]
other = "не удалось распарсить sl=0x{{.Value}}"

[CurrentCode]
other = "Текущий вход: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Текущий вход: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[MissingDDC]
other = "ddcutil не найден. Установи: sudo dnf install ddcutil i2c-tools"

[I2CAccess]
other = "Добавь пользователя в группу i2c: sudo usermod -aG i2c $USER, затем перелогинься"

[NoI2CDevices]
other = "I2C устройства не найдены. Загрузи модуль: sudo modprobe i2c-dev. Если не помогло, проверь DDC/CI в OSD монитора"

[DDCRetry]
other = "монитор найден, но DDC/VCP-запрос не получил стабильный ответ. Это не похоже на проблему sudo. Проверь DDC/CI в OSD монитора, кабель/док/порт; переключение может сработать даже если чтение текущего входа нестабильно"

[DDCTimeout]
other = "команда ddcutil превысила таймаут {{.Seconds}}s: ddcutil {{.Command}}"

[DDCTimeoutWithOutput]
other = "команда ddcutil превысила таймаут {{.Seconds}}s: ddcutil {{.Command}}\n{{.Raw}}"

[DDCCommandFailedNoOutput]
other = "ddcutil {{.Command}}: {{.Error}}"

[DDCCommandFailedWithOutput]
other = "ddcutil {{.Command}}: {{.Error}}\n{{.Raw}}"

[AutoDetectFailed]
other = "автоопределение bus не удалось"

[AutoDetectNoDisplays]
other = "автоопределение bus не нашло дисплеи в выводе ddcutil detect --brief"

[AutoDetectNoResponsive]
other = "автоопределение bus проверило {{.Buses}}, но ни один bus не ответил на getvcp {{.VCP}}"

[NoMonitorUnknownBus]
other = "монитор не найден на указанном bus. Проверь доступные дисплеи: dellkvm detect. Затем исправь bus в config.toml или оставь bus = 0 для автоопределения"

[NoMonitorBus]
other = "монитор не найден на bus {{.Bus}}. Проверь доступные дисплеи: dellkvm detect. Затем исправь bus в config.toml или оставь bus = 0 для автоопределения"

[ListTitle]
other = "Входы"

[Refreshing]
other = "Обновляю текущий вход..."

[Switching]
other = "Переключаю на {{.Name}}..."

[ErrorPrefix]
other = "Ошибка: {{.Error}}"

[CurrentRefreshed]
other = "Текущий вход обновлен через bus {{.Bus}}."

[CurrentCodeNotFound]
other = "Текущий код отсутствует в списке входов."

[TuiLastSelected]
other = "Последний выбранный вход: {{.Name}} (bus {{.Bus}}, не подтверждено)"

[TuiCurrentUnknown]
other = "Текущий вход: неизвестно"

[TuiCurrentInput]
other = "Текущий вход: {{.Name}} ({{.Code}}, bus {{.Bus}})"

[TuiCurrentCode]
other = "Текущий вход: {{.Code}} (bus {{.Bus}})"

[Keys]
other = "Enter: переключить  r: обновить  q/Ctrl+C: выйти"
`,
	"fr": `[UsageDetect]
other = "Utilisation : dellkvm detect"

[UsageCurrent]
other = "Utilisation : dellkvm current"

[UsageSwitch]
other = "Utilisation : dellkvm switch <id>"

[UsageRoot]
other = "Utilisation : dellkvm [tray|tui|init|detect|current|switch <id>|status --json|help]"

[Help]
other = "Utilisation : dellkvm [commande]\n\nCommandes:\n  init         créer la configuration si absente\n  tray         ouvrir la zone de notification\n  tui          ouvrir l’interface terminal\n  status --json  afficher le statut JSON\n  detect       liste les moniteurs disponibles\n  current      affiche l'entrée actuelle\n  switch <id>  bascule vers une entrée configurée\n  help         affiche cette aide\n\nLance sans commande pour ouvrir la zone de notification. dellkvm tui ouvre la TUI.\nLa configuration est facultative ; lance dellkvm init pour personnaliser les entrées. Utilise dellkvm detect pour lister les moniteurs, ou garde bus = 0 pour l'auto-détection."

[ConfigReadFailed]
other = "impossible de lire {{.Path}}"

[ConfigNotFound]
other = "configuration introuvable. Chemins vérifiés : {{.Paths}}\nCopie config.toml.default vers config.toml ou ~/.config/dellkvm/config.toml et modifie-le.\nExécute dellkvm init pour créer une configuration initiale."

[InvalidBus]
other = "dans config.toml, bus doit valoir 0 pour l'auto-détection ou un nombre positif"

[MissingInputs]
other = "config.toml doit contenir au moins une entrée"

[InputNotFound]
other = "l'entrée avec l'id {{.ID}} est introuvable"

[SwitchSentNoVerify]
other = "Commande de bascule envoyée vers {{.Name}} ({{.ID}}, {{.Code}}) via le bus {{.Bus}}, mais le moniteur ne répond pas encore à la vérification getvcp {{.VCP}} : {{.Error}}"

[SwitchVerifiedCode]
other = "Basculé vers {{.Name}} ({{.ID}}, {{.Code}}) via le bus {{.Bus}}. Vérification : code actuel {{.CurrentCode}}"

[SwitchVerifiedInput]
other = "Basculé vers {{.Name}} ({{.ID}}, {{.Code}}) via le bus {{.Bus}}. Vérification : {{.CurrentName}} ({{.CurrentCode}})"

[SwitchVerificationMismatch]
other = "Commande de bascule envoyée vers {{.Name}} ({{.ID}}, {{.Requested}}) via le bus {{.Bus}}, mais la vérification indique {{.Observed}}."

[AutoDetectedPrefix]
other = "Bus {{.Bus}} auto-détecté. {{.Message}}"

[VCPCodeNotFound]
other = "sl=0x.. est introuvable dans la sortie getvcp :\n{{.Raw}}"

[VCPCodeParseFailed]
other = "impossible d'analyser sl=0x{{.Value}}"

[CurrentCode]
other = "Entrée actuelle : {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Entrée actuelle : {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[MissingDDC]
other = "ddcutil est introuvable. Installe-le avec : sudo dnf install ddcutil i2c-tools"

[I2CAccess]
other = "Ajoute ton utilisateur au groupe i2c : sudo usermod -aG i2c $USER, puis reconnecte-toi"

[NoI2CDevices]
other = "Aucun périphérique I2C trouvé. Charge le module : sudo modprobe i2c-dev. Si cela ne suffit pas, vérifie DDC/CI dans l'OSD du moniteur"

[DDCRetry]
other = "le moniteur a été trouvé, mais la requête DDC/VCP n'a pas obtenu de réponse stable. Cela ne ressemble pas à un problème sudo. Vérifie DDC/CI dans l'OSD du moniteur, le câble, le dock ou le port ; la bascule peut quand même fonctionner si la lecture de l'entrée actuelle est instable"

[DDCTimeout]
other = "la commande ddcutil a dépassé le délai de {{.Seconds}}s : ddcutil {{.Command}}"

[DDCTimeoutWithOutput]
other = "la commande ddcutil a dépassé le délai de {{.Seconds}}s : ddcutil {{.Command}}\n{{.Raw}}"

[DDCCommandFailedNoOutput]
other = "ddcutil {{.Command}} : {{.Error}}"

[DDCCommandFailedWithOutput]
other = "ddcutil {{.Command}} : {{.Error}}\n{{.Raw}}"

[AutoDetectFailed]
other = "échec de l'auto-détection du bus"

[AutoDetectNoDisplays]
other = "l'auto-détection du bus n'a trouvé aucun écran dans la sortie de ddcutil detect --brief"

[AutoDetectNoResponsive]
other = "l'auto-détection du bus a vérifié {{.Buses}}, mais aucun bus n'a répondu à getvcp {{.VCP}}"

[NoMonitorUnknownBus]
other = "moniteur introuvable sur le bus sélectionné. Vérifie les écrans disponibles : dellkvm detect. Corrige ensuite bus dans config.toml, ou garde bus = 0 pour l'auto-détection"

[NoMonitorBus]
other = "moniteur introuvable sur le bus {{.Bus}}. Vérifie les écrans disponibles : dellkvm detect. Corrige ensuite bus dans config.toml, ou garde bus = 0 pour l'auto-détection"

[ListTitle]
other = "Entrées"

[Refreshing]
other = "Actualisation de l'entrée actuelle..."

[Switching]
other = "Bascule vers {{.Name}}..."

[ErrorPrefix]
other = "Erreur : {{.Error}}"

[CurrentRefreshed]
other = "Entrée actuelle actualisée via le bus {{.Bus}}."

[CurrentCodeNotFound]
other = "Le code actuel ne figure pas dans la liste des entrées."

[TuiLastSelected]
other = "Dernière entrée choisie : {{.Name}} (bus {{.Bus}}, non vérifiée)"

[TuiCurrentUnknown]
other = "Entrée actuelle : inconnue"

[TuiCurrentInput]
other = "Entrée actuelle : {{.Name}} ({{.Code}}, bus {{.Bus}})"

[TuiCurrentCode]
other = "Entrée actuelle : {{.Code}} (bus {{.Bus}})"

[Keys]
other = "Enter : basculer  r : actualiser  q/Ctrl+C : quitter"
`,
	"de": `[UsageDetect]
other = "Verwendung: dellkvm detect"

[UsageCurrent]
other = "Verwendung: dellkvm current"

[UsageSwitch]
other = "Verwendung: dellkvm switch <id>"

[UsageRoot]
other = "Verwendung: dellkvm [tray|tui|init|detect|current|switch <id>|status --json|help]"

[Help]
other = "Verwendung: dellkvm [Befehl]\n\nBefehle:\n  init         fehlende Startkonfiguration erstellen\n  tray         Infobereich öffnen\n  tui          Terminal-Oberfläche öffnen\n  status --json  Status als JSON ausgeben\n  detect       listet verfügbare Monitore auf\n  current      zeigt den aktuellen Eingang\n  switch <id>  schaltet auf einen konfigurierten Eingang\n  help         zeigt diese Hilfe\n\nOhne Befehl wird das Infobereich-Symbol geöffnet. dellkvm tui öffnet die TUI.\nDie Konfiguration ist optional; mit dellkvm init kannst du Eingänge anpassen. dellkvm detect listet Monitore auf; bus = 0 aktiviert die Auto-Erkennung."

[ConfigReadFailed]
other = "{{.Path}} konnte nicht gelesen werden"

[ConfigNotFound]
other = "Konfiguration nicht gefunden. Geprüfte Pfade: {{.Paths}}\nKopiere config.toml.default nach config.toml oder ~/.config/dellkvm/config.toml und bearbeite sie.\nFühre dellkvm init aus, um eine Startkonfiguration zu erstellen."

[InvalidBus]
other = "in config.toml muss bus 0 für Auto-Erkennung oder eine positive Nummer sein"

[MissingInputs]
other = "config.toml muss mindestens einen input enthalten"

[InputNotFound]
other = "input mit id {{.ID}} wurde nicht gefunden"

[SwitchSentNoVerify]
other = "Umschaltbefehl an {{.Name}} ({{.ID}}, {{.Code}}) über bus {{.Bus}} gesendet, aber der Monitor antwortet noch nicht auf die getvcp-{{.VCP}}-Prüfung: {{.Error}}"

[SwitchVerifiedCode]
other = "Auf {{.Name}} ({{.ID}}, {{.Code}}) über bus {{.Bus}} umgeschaltet. Prüfung: aktueller Code {{.CurrentCode}}"

[SwitchVerifiedInput]
other = "Auf {{.Name}} ({{.ID}}, {{.Code}}) über bus {{.Bus}} umgeschaltet. Prüfung: {{.CurrentName}} ({{.CurrentCode}})"

[SwitchVerificationMismatch]
other = "Umschaltbefehl an {{.Name}} ({{.ID}}, {{.Requested}}) über bus {{.Bus}} gesendet, aber die Prüfung meldete {{.Observed}}."

[AutoDetectedPrefix]
other = "Bus {{.Bus}} automatisch erkannt. {{.Message}}"

[VCPCodeNotFound]
other = "sl=0x.. wurde in der getvcp-Ausgabe nicht gefunden:\n{{.Raw}}"

[VCPCodeParseFailed]
other = "sl=0x{{.Value}} konnte nicht geparst werden"

[CurrentCode]
other = "Aktueller Eingang: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Aktueller Eingang: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[MissingDDC]
other = "ddcutil wurde nicht gefunden. Installiere es mit: sudo dnf install ddcutil i2c-tools"

[I2CAccess]
other = "Füge deinen Benutzer zur Gruppe i2c hinzu: sudo usermod -aG i2c $USER, melde dich danach neu an"

[NoI2CDevices]
other = "Keine I2C-Geräte gefunden. Lade das Modul: sudo modprobe i2c-dev. Wenn das nicht hilft, prüfe DDC/CI im OSD des Monitors"

[DDCRetry]
other = "Monitor gefunden, aber die DDC/VCP-Anfrage bekam keine stabile Antwort. Das sieht nicht nach einem sudo-Problem aus. Prüfe DDC/CI im OSD des Monitors, Kabel, Dock oder Port; Umschalten kann trotzdem funktionieren, auch wenn das Lesen des aktuellen Eingangs instabil ist"

[DDCTimeout]
other = "ddcutil-Befehl hat nach {{.Seconds}}s das Zeitlimit erreicht: ddcutil {{.Command}}"

[DDCTimeoutWithOutput]
other = "ddcutil-Befehl hat nach {{.Seconds}}s das Zeitlimit erreicht: ddcutil {{.Command}}\n{{.Raw}}"

[DDCCommandFailedNoOutput]
other = "ddcutil {{.Command}}: {{.Error}}"

[DDCCommandFailedWithOutput]
other = "ddcutil {{.Command}}: {{.Error}}\n{{.Raw}}"

[AutoDetectFailed]
other = "bus-Auto-Erkennung fehlgeschlagen"

[AutoDetectNoDisplays]
other = "bus-Auto-Erkennung fand keine Displays in der Ausgabe von ddcutil detect --brief"

[AutoDetectNoResponsive]
other = "bus-Auto-Erkennung prüfte {{.Buses}}, aber kein bus antwortete auf getvcp {{.VCP}}"

[NoMonitorUnknownBus]
other = "Monitor auf dem ausgewählten bus nicht gefunden. Prüfe verfügbare Displays: dellkvm detect. Korrigiere danach bus in config.toml, oder behalte bus = 0 für Auto-Erkennung"

[NoMonitorBus]
other = "Monitor auf bus {{.Bus}} nicht gefunden. Prüfe verfügbare Displays: dellkvm detect. Korrigiere danach bus in config.toml, oder behalte bus = 0 für Auto-Erkennung"

[ListTitle]
other = "Eingänge"

[Refreshing]
other = "Aktueller Eingang wird aktualisiert..."

[Switching]
other = "Schalte auf {{.Name}} um..."

[ErrorPrefix]
other = "Fehler: {{.Error}}"

[CurrentRefreshed]
other = "Aktueller Eingang über bus {{.Bus}} aktualisiert."

[CurrentCodeNotFound]
other = "Der aktuelle Code steht nicht in der Eingangsliste."

[TuiLastSelected]
other = "Zuletzt gewählt: {{.Name}} (Bus {{.Bus}}, nicht bestätigt)"

[TuiCurrentUnknown]
other = "Aktueller Eingang: unbekannt"

[TuiCurrentInput]
other = "Aktueller Eingang: {{.Name}} ({{.Code}}, bus {{.Bus}})"

[TuiCurrentCode]
other = "Aktueller Eingang: {{.Code}} (bus {{.Bus}})"

[Keys]
other = "Enter: umschalten  r: aktualisieren  q/Ctrl+C: beenden"
`,
	"zh": `[UsageDetect]
other = "用法：dellkvm detect"

[UsageCurrent]
other = "用法：dellkvm current"

[UsageSwitch]
other = "用法：dellkvm switch <id>"

[UsageRoot]
other = "用法：dellkvm [tray|tui|init|detect|current|switch <id>|status --json|help]"

[Help]
other = "用法：dellkvm [command]\n\n命令：\n  init         创建缺失的初始配置\n  tray         打开系统托盘\n  tui          打开终端界面\n  status --json  输出 JSON 状态\n  detect       列出可用显示器\n  current      显示当前输入\n  switch <id>  切换到已配置的输入\n  help         显示此帮助\n\n不带命令运行会打开系统托盘。使用 dellkvm tui 打开终端界面。\n配置可选；运行 dellkvm init 可自定义输入。使用 dellkvm detect 查看显示器，或保留 bus = 0 自动检测。"

[ConfigReadFailed]
other = "无法读取 {{.Path}}"

[ConfigNotFound]
other = "未找到配置。已检查路径：{{.Paths}}\n将 config.toml.default 复制到 config.toml 或 ~/.config/dellkvm/config.toml 并编辑它。\n运行 dellkvm init 创建初始配置。"

[InvalidBus]
other = "config.toml 中 bus 必须为 0（自动检测）或正数"

[MissingInputs]
other = "config.toml 必须至少包含一个 input"

[InputNotFound]
other = "未找到 id 为 {{.ID}} 的 input"

[SwitchSentNoVerify]
other = "已向 {{.Name}} ({{.ID}}, {{.Code}}) 通过 bus {{.Bus}} 发送切换命令，但显示器暂时未响应 getvcp {{.VCP}} 验证：{{.Error}}"

[SwitchVerifiedCode]
other = "已切换到 {{.Name}} ({{.ID}}, {{.Code}})，通过 bus {{.Bus}}。验证：当前代码 {{.CurrentCode}}"

[SwitchVerifiedInput]
other = "已切换到 {{.Name}} ({{.ID}}, {{.Code}})，通过 bus {{.Bus}}。验证：{{.CurrentName}} ({{.CurrentCode}})"

[SwitchVerificationMismatch]
other = "已向 {{.Name}} ({{.ID}}, {{.Requested}}) 通过 bus {{.Bus}} 发送切换命令，但验证结果为 {{.Observed}}。"

[AutoDetectedPrefix]
other = "自动检测到 bus {{.Bus}}。{{.Message}}"

[VCPCodeNotFound]
other = "getvcp 输出中未找到 sl=0x..：\n{{.Raw}}"

[VCPCodeParseFailed]
other = "无法解析 sl=0x{{.Value}}"

[CurrentCode]
other = "当前输入：{{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "当前输入：{{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[MissingDDC]
other = "未找到 ddcutil。请安装：sudo dnf install ddcutil i2c-tools"

[I2CAccess]
other = "将当前用户加入 i2c 组：sudo usermod -aG i2c $USER，然后重新登录"

[NoI2CDevices]
other = "未找到 I2C 设备。加载模块：sudo modprobe i2c-dev。如果仍然无效，请检查显示器 OSD 中的 DDC/CI"

[DDCRetry]
other = "已找到显示器，但 DDC/VCP 请求没有获得稳定响应。这看起来不像 sudo 问题。请检查显示器 OSD 中的 DDC/CI、线缆、扩展坞或端口；即使读取当前输入不稳定，切换仍可能可用"

[DDCTimeout]
other = "ddcutil 命令在 {{.Seconds}}s 后超时：ddcutil {{.Command}}"

[DDCTimeoutWithOutput]
other = "ddcutil 命令在 {{.Seconds}}s 后超时：ddcutil {{.Command}}\n{{.Raw}}"

[DDCCommandFailedNoOutput]
other = "ddcutil {{.Command}}：{{.Error}}"

[DDCCommandFailedWithOutput]
other = "ddcutil {{.Command}}：{{.Error}}\n{{.Raw}}"

[AutoDetectFailed]
other = "bus 自动检测失败"

[AutoDetectNoDisplays]
other = "bus 自动检测在 ddcutil detect --brief 输出中未找到显示器"

[AutoDetectNoResponsive]
other = "bus 自动检测检查了 {{.Buses}}，但没有 bus 响应 getvcp {{.VCP}}"

[NoMonitorUnknownBus]
other = "在所选 bus 上未找到显示器。查看可用显示器：dellkvm detect。然后修正 config.toml 中的 bus，或保留 bus = 0 进行自动检测"

[NoMonitorBus]
other = "在 bus {{.Bus}} 上未找到显示器。查看可用显示器：dellkvm detect。然后修正 config.toml 中的 bus，或保留 bus = 0 进行自动检测"

[ListTitle]
other = "输入"

[Refreshing]
other = "正在刷新当前输入..."

[Switching]
other = "正在切换到 {{.Name}}..."

[ErrorPrefix]
other = "错误：{{.Error}}"

[CurrentRefreshed]
other = "已通过 bus {{.Bus}} 刷新当前输入。"

[CurrentCodeNotFound]
other = "当前代码不在输入列表中。"

[TuiLastSelected]
other = "上次选择的输入：{{.Name}} (bus {{.Bus}}，未验证)"

[TuiCurrentUnknown]
other = "当前输入：未知"

[TuiCurrentInput]
other = "当前输入：{{.Name}} ({{.Code}}, bus {{.Bus}})"

[TuiCurrentCode]
other = "当前输入：{{.Code}} (bus {{.Bus}})"

[Keys]
other = "Enter：切换  r：刷新  q/Ctrl+C：退出"
`,
}
