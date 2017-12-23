(function () {
    $.ajax({
        type: 'POST',
        url: contextPath + "/logs",
        success: function (content, textStatus, request) {
            var logsHtml = createLogsTable(content)
            $('#logs').html(logsHtml)
            createContextMenu()
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })

    function createLogsTable(content) {
        var logsHtml = '<table>' +
            '<tr>' +
            '<td>Logger Name</td>' +
            '<td>Log Path</td>' +
            '<td>State</td>' +
            '<td>Machine</td>' +
            '<td>Size</td>' +
            '<td>Last Modified</td>' +
            '<td>ProcessInfo</td>' +
            '<td>Cost</td>' +
            '</tr>'

        if (content && content.length) {
            for (var j = 0; j < content.length; j++) {
                var log = content[j]
                var logMachines = log.Logs

                for (var i = 0; i < logMachines.length; ++i) {
                    var logMachine = log.Logs[i]

                    logsHtml += '<tr>'
                    if (i == 0) {
                        logsHtml += '<td rowspan="' + logMachines.length + '" class="LoggerName">' + log.Logger + '</td>'
                            + '<td class="LogPath" rowspan="' + logMachines.length + '">' + log.LogPath + '</td>'
                    } else {
                        logsHtml += '<td class="hidden LoggerName">' + log.Logger + '</td>'
                            + '<td class="LogPath hidden">' + log.LogPath + '</td>'
                    }

                    logsHtml += '<td>' + (logMachine.Error || 'OK') + '</td>'
                        + '<td class="LogMachine">' + logMachine.MachineName + '</td>'
                        + '<td class="LogFileSize">' + logMachine.FileSize + '</td>'
                        + '<td>' + logMachine.LastModified + '</td>'
                        + '<td class="ProcessInfo">' + logMachine.ProcessInfo + '</td>'
                        + '<td>' + logMachine.CostTime + '</td>'
                        + '</tr>'
                }
            }
        }

        logsHtml += '</table>'
        return logsHtml;
    }

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

    function createContextMenu() {
        $.contextMenu({
            selector: '.LogFileSize',
            callback: function (key, options) {
                if (key === "TruncateLogFile") {
                    var $cell = $(this);
                    var $row = $cell.parent();
                    var logMachine = $row.find('td.LogMachine').text();
                    var loggerName = $row.find('td.LoggerName').text();
                    TruncateLogFile($cell, logMachine, loggerName)
                }
            },
            items: {
                "TruncateLogFile": {name: "TruncateLogFile", icon: "cut"}
            }
        })
    }

})()