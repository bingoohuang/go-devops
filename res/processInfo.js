(function () {
    function RestartProcess($cell, logMachine, loggerName) {
        $.ajax({
            type: 'POST',
            url: contextPath + "/restartProcess/" + loggerName + "/" + logMachine,
            success: function (content, textStatus, request) {
                if (content.Error !== "") {
                    alert(content.Error)
                    return
                }

                $cell.addClass('changed').text(content.ProcessInfo)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $.createProcessInfoContextMenu = function () {
        $.contextMenu({
            selector: '.ProcessInfo',
            callback: function (key, options) {
                var $cell = $(this)
                var $row = $cell.parent()
                var logMachine = $row.find('td.LogMachine').text()
                var loggerName = $row.find('td.LoggerName').text()
                if (key === "RestartProcess") {
                    RestartProcess($cell, logMachine, loggerName)
                }
            },
            items: {
                RestartProcess: {name: "Restart process", icon: "restart"}
            }
        })
    }
})()