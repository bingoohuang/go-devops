(function () {
    $('#closeFileContent').click(function () {
        $('#fileContent').hide()
        $('#tableArea').show()
    })

    function TailLogFile(loggerName, logPath) {
        $('#refresh').unbind('click')
        $.ajax({
            type: 'POST',
            url: contextPath + "/tailLogFile/" + loggerName,
            success: function (content, textStatus, request) {
                var tailContent = '<pre class="preWrap">'

                for (var i = 0; i < content.length; ++i) {
                    if (i > 0) tailContent += '\n\n'
                    tailContent += content[i].MachineName + ' ' + logPath + ':\n'
                        + content[i].TailContent
                }
                tailContent += '</pre>'
                $('#preContent').html(tailContent)

                $('#tableArea').hide()
                $('#fileContent').show()
                $('#refresh').click(function () {
                    TailLogFile(loggerName, logPath)
                })
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
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

    $.createLogFileTailContextMenu = function () {
        $.contextMenu({
            selector: '.LogPath',
            callback: function (key, options) {
                if (key === "TailLogFile") {
                    var $row = $(this).parent();
                    var loggerName = $row.find('td.LoggerName').text();
                    var logPath = $row.find('td.LogPath').text();
                    TailLogFile(loggerName, logPath)
                }
            },
            items: {
                "TailLogFile": {name: "Tail Log File", icon: "tail"}
            }
        })
    }

    $.createLogFileSizeContextMenu = function () {
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
                "TruncateLogFile": {name: "Truncate Log File", icon: "truncate"}
            }
        })
    }

})()