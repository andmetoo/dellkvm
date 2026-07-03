//go:build linux

package main

// localeData is generated from cmd/dellkvm/locales/active.*.toml.
var localeData = map[string]string{
	"en": `[UsageDetect]
other = "Usage: dellkvm detect"

[UsageCurrent]
other = "Usage: dellkvm current"

[UsageSwitch]
other = "Usage: dellkvm switch <id>"

[UsageLearn]
other = "Usage: dellkvm learn"

[UsageRoot]
other = "Usage: dellkvm [detect|current|switch <id>|learn]"

[ConfigReadFailed]
other = "failed to read {{.Path}}"

[ConfigNotFound]
other = "config not found. Checked paths: {{.Paths}}\nRun: dellkvm learn"

[InvalidBus]
other = "config.toml bus must be 0 for auto-detect or a positive number"

[MissingInputs]
other = "config.toml must contain at least one input"

[InputNotFound]
other = "input with id {{.ID}} was not found in config.toml"

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

[LearnPrompt]
other = "Switch the monitor to {{.Name}} through OSD and press Enter..."

[LearnCodeReadFailed]
other = "failed to read code for {{.ID}}"

[ConfigSaved]
other = "Config saved: {{.Path}}"

[CurrentCode]
other = "Current input: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Current input: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[BusPrompt]
other = "Bus: "

[InvalidBusPrompt]
other = "Enter a positive bus number, for example 6."

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
other = "monitor not found on the selected bus. Check available displays: dellkvm detect. Then fix bus in config.toml or run: dellkvm learn"

[NoMonitorBus]
other = "monitor not found on bus {{.Bus}}. Check available displays: dellkvm detect. Then fix bus in config.toml or run: dellkvm learn"

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
other = "Current code was not found in config.toml."

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

[UsageLearn]
other = "Использование: dellkvm learn"

[UsageRoot]
other = "Использование: dellkvm [detect|current|switch <id>|learn]"

[ConfigReadFailed]
other = "не удалось прочитать {{.Path}}"

[ConfigNotFound]
other = "конфиг не найден. Проверены пути: {{.Paths}}\nЗапусти: dellkvm learn"

[InvalidBus]
other = "в config.toml bus должен быть 0 для auto или положительным номером"

[MissingInputs]
other = "в config.toml должен быть хотя бы один input"

[InputNotFound]
other = "input с id {{.ID}} не найден в config.toml"

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

[LearnPrompt]
other = "Переключи монитор на {{.Name}} через OSD и нажми Enter..."

[LearnCodeReadFailed]
other = "не удалось прочитать код для {{.ID}}"

[ConfigSaved]
other = "Конфиг сохранен: {{.Path}}"

[CurrentCode]
other = "Текущий вход: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Текущий вход: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[BusPrompt]
other = "Bus: "

[InvalidBusPrompt]
other = "Введи положительный номер bus, например 6."

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
other = "монитор не найден на указанном bus. Проверь доступные дисплеи: dellkvm detect. Затем исправь bus в config.toml или запусти: dellkvm learn"

[NoMonitorBus]
other = "монитор не найден на bus {{.Bus}}. Проверь доступные дисплеи: dellkvm detect. Затем исправь bus в config.toml или запусти: dellkvm learn"

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
other = "Текущий код не найден в config.toml."

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

[UsageLearn]
other = "Utilisation : dellkvm learn"

[UsageRoot]
other = "Utilisation : dellkvm [detect|current|switch <id>|learn]"

[ConfigReadFailed]
other = "impossible de lire {{.Path}}"

[ConfigNotFound]
other = "configuration introuvable. Chemins vérifiés : {{.Paths}}\nExécute : dellkvm learn"

[InvalidBus]
other = "dans config.toml, bus doit valoir 0 pour l'auto-détection ou un nombre positif"

[MissingInputs]
other = "config.toml doit contenir au moins une entrée"

[InputNotFound]
other = "l'entrée avec l'id {{.ID}} est introuvable dans config.toml"

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

[LearnPrompt]
other = "Bascule le moniteur vers {{.Name}} via l'OSD puis appuie sur Entrée..."

[LearnCodeReadFailed]
other = "impossible de lire le code pour {{.ID}}"

[ConfigSaved]
other = "Configuration enregistrée : {{.Path}}"

[CurrentCode]
other = "Entrée actuelle : {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Entrée actuelle : {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[BusPrompt]
other = "Bus : "

[InvalidBusPrompt]
other = "Saisis un numéro de bus positif, par exemple 6."

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
other = "moniteur introuvable sur le bus sélectionné. Vérifie les écrans disponibles : dellkvm detect. Corrige ensuite bus dans config.toml ou exécute : dellkvm learn"

[NoMonitorBus]
other = "moniteur introuvable sur le bus {{.Bus}}. Vérifie les écrans disponibles : dellkvm detect. Corrige ensuite bus dans config.toml ou exécute : dellkvm learn"

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
other = "Le code actuel est introuvable dans config.toml."

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

[UsageLearn]
other = "Verwendung: dellkvm learn"

[UsageRoot]
other = "Verwendung: dellkvm [detect|current|switch <id>|learn]"

[ConfigReadFailed]
other = "{{.Path}} konnte nicht gelesen werden"

[ConfigNotFound]
other = "Konfiguration nicht gefunden. Geprüfte Pfade: {{.Paths}}\nAusführen: dellkvm learn"

[InvalidBus]
other = "in config.toml muss bus 0 für Auto-Erkennung oder eine positive Nummer sein"

[MissingInputs]
other = "config.toml muss mindestens einen input enthalten"

[InputNotFound]
other = "input mit id {{.ID}} wurde in config.toml nicht gefunden"

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

[LearnPrompt]
other = "Schalte den Monitor per OSD auf {{.Name}} um und drücke Enter..."

[LearnCodeReadFailed]
other = "Code für {{.ID}} konnte nicht gelesen werden"

[ConfigSaved]
other = "Konfiguration gespeichert: {{.Path}}"

[CurrentCode]
other = "Aktueller Eingang: {{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "Aktueller Eingang: {{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[BusPrompt]
other = "Bus: "

[InvalidBusPrompt]
other = "Gib eine positive bus-Nummer ein, zum Beispiel 6."

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
other = "Monitor auf dem ausgewählten bus nicht gefunden. Prüfe verfügbare Displays: dellkvm detect. Korrigiere danach bus in config.toml oder führe aus: dellkvm learn"

[NoMonitorBus]
other = "Monitor auf bus {{.Bus}} nicht gefunden. Prüfe verfügbare Displays: dellkvm detect. Korrigiere danach bus in config.toml oder führe aus: dellkvm learn"

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
other = "Aktueller Code wurde in config.toml nicht gefunden."

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

[UsageLearn]
other = "用法：dellkvm learn"

[UsageRoot]
other = "用法：dellkvm [detect|current|switch <id>|learn]"

[ConfigReadFailed]
other = "无法读取 {{.Path}}"

[ConfigNotFound]
other = "未找到配置。已检查路径：{{.Paths}}\n运行：dellkvm learn"

[InvalidBus]
other = "config.toml 中 bus 必须为 0（自动检测）或正数"

[MissingInputs]
other = "config.toml 必须至少包含一个 input"

[InputNotFound]
other = "config.toml 中未找到 id 为 {{.ID}} 的 input"

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

[LearnPrompt]
other = "请通过 OSD 将显示器切换到 {{.Name}}，然后按 Enter..."

[LearnCodeReadFailed]
other = "无法读取 {{.ID}} 的代码"

[ConfigSaved]
other = "配置已保存：{{.Path}}"

[CurrentCode]
other = "当前输入：{{.Code}} (bus {{.Bus}})"

[CurrentInput]
other = "当前输入：{{.Name}} ({{.ID}}, {{.Code}}, bus {{.Bus}})"

[BusPrompt]
other = "Bus: "

[InvalidBusPrompt]
other = "请输入正数 bus 编号，例如 6。"

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
other = "在所选 bus 上未找到显示器。查看可用显示器：dellkvm detect。然后修正 config.toml 中的 bus 或运行：dellkvm learn"

[NoMonitorBus]
other = "在 bus {{.Bus}} 上未找到显示器。查看可用显示器：dellkvm detect。然后修正 config.toml 中的 bus 或运行：dellkvm learn"

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
other = "当前代码未在 config.toml 中找到。"

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
