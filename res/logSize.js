(function () {
    function TruncateLogFile($cell, logMachine, loggerName) {
        $.ajax({
            type: 'POST',
            url: contextPath + "/truncateLogFile/" + loggerName + "/" + logMachine,
            success: function (content, textStatus, request) {
                $cell.addClass('changed').text(content.FileSize)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $.createLogFileSizeContextMenu = function () {
        $.contextMenu({
            selector: '.LogFileSize',
            callback: function (key, options) {
                var $cell = $(this)
                var $row = $cell.parent()
                var logMachine = $row.find('td.LogMachine').text()
                var loggerName = $row.find('td.LoggerName').text()
                if (key === "TruncateLogFile") {
                    TruncateLogFile($cell, logMachine, loggerName)
                }
            },
            items: {
                TruncateLogFile: {name: "Truncate Log File", icon: "truncate"}
            }
        })
    }
})()