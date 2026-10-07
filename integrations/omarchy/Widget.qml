import QtQuick
import QtQuick.Controls as Controls
import QtQuick.Layouts
import Quickshell
import Quickshell.Io
import qs.Ui

BarWidget {
    id: root
    moduleName: "dellkvm.inputs"
    property var inputs: []
    property string currentCode: ""
    property string message: "Open to read monitor status"
    property string stdoutText: ""
    property string stderrText: ""
    property bool readingStatus: false
    readonly property bool busy: operation.running
    readonly property bool opened: popup.visible

    implicitWidth: button.implicitWidth
    implicitHeight: button.implicitHeight

    function open() {
        popup.visible = true
        refresh()
    }
    function close() { popup.visible = false }
    function execute(args, isStatus) {
        if (busy) return
        readingStatus = isStatus
        stdoutText = ""
        stderrText = ""
        currentCode = ""
        message = isStatus ? "Reading monitor…" : "Switching input…"
        operation.command = [String(setting("executable", "dellkvm"))].concat(args)
        operation.running = true
    }
    function refresh() { execute(["status", "--json"], true) }

    WidgetButton {
        id: button
        anchors.fill: parent
        bar: root.bar
        text: "KVM"
        tooltipText: root.message
        onPressed: mouseButton => {
            if (mouseButton === Qt.LeftButton) {
                if (root.opened) root.close()
                else root.open()
            }
        }
    }

    PopupWindow {
        id: popup
        anchor.item: button
        anchor.edges: root.vertical ? Edges.Right : Edges.Bottom
        anchor.gravity: root.vertical ? Edges.Right : Edges.Bottom
        implicitWidth: 340
        implicitHeight: Math.min(content.implicitHeight + 28, 520)
        color: root.bar ? root.bar.barBackground : "#20252b"

        Controls.ScrollView {
            anchors.fill: parent
            anchors.margins: 14
            contentWidth: availableWidth
            ColumnLayout {
                id: content
                width: parent.width
                spacing: 8
                Controls.Label {
                    text: "Monitor inputs"
                    font.bold: true
                    color: root.bar ? root.bar.barForeground : "white"
                }
                Controls.Label {
                    Layout.fillWidth: true
                    text: root.message
                    textFormat: Text.PlainText
                    wrapMode: Text.Wrap
                    color: root.bar ? root.bar.barForeground : "white"
                }
                Repeater {
                    model: root.inputs
                    Controls.Button {
                        required property var modelData
                        Layout.fillWidth: true
                        text: (root.currentCode === modelData.code ? "● " : "") + modelData.name
                        enabled: !root.busy
                        onClicked: root.execute(["switch", modelData.id], false)
                    }
                }
                RowLayout {
                    Controls.Button { text: "Refresh"; enabled: !root.busy; onClicked: root.refresh() }
                    Controls.Button { text: "Close"; onClicked: root.close() }
                }
            }
        }
    }

    Process {
        id: operation
        stdout: StdioCollector { onStreamFinished: root.stdoutText = text }
        stderr: StdioCollector { onStreamFinished: root.stderrText = text }
        onExited: (exitCode, exitStatus) => {
            if (exitCode !== 0 || exitStatus !== 0) {
                root.message = root.stderrText.trim() || "Could not run dellkvm. Check the executable path."
                return
            }
            if (!root.readingStatus) {
                // Keep the switch result: another read may fail because this
                // host is intentionally disconnected by the input change.
                root.message = root.stdoutText.trim()
                return
            }
            try {
                const data = JSON.parse(root.stdoutText)
                root.inputs = data.inputs || []
                root.currentCode = data.code || ""
                root.message = data.error || ("Monitor " + data.bus + " · " + data.code)
            } catch (error) {
                root.message = "Invalid response from dellkvm: " + error
            }
        }
    }
}
