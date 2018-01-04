(function () {
    var tomlEditor = null

    function LoadConfig() {
        $.ajax({
            type: 'POST',
            url: contextPath + "/loadConfig",
            success: function (content, textStatus, request) {
                tomlEditor.setValue(content.Conf)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function toogleConfigDiv(show) {
        $('.hiddable').hide()
        if (show) {
            $('#configDiv').show()
        } else {
            $('#tableArea').show()
        }
    }

    $('#CloseConfig').click(function () {
        toogleConfigDiv()
    })

    $('#SaveConfig').click(function () {
        $.ajax({
            type: 'POST',
            url: contextPath + "/saveConfig",
            data: {config: tomlEditor.getValue()},
            success: function (content, textStatus, request) {
                if (content.OK === "OK") {
                    toogleConfigDiv()
                } else {
                    alert(content)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    })

    $('#ConfigBtn').click(function () {
        toogleConfigDiv(true)

        if (tomlEditor == null) {
            tomlEditor = CodeMirror.fromTextArea(document.getElementById("tomlEditor"), {
                mode: 'text/x-toml', lineNumbers: true
            })
        }

        LoadConfig()
    })
})()