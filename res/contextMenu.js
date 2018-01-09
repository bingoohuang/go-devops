(function () {
    $('#closeFileContent').click(function () {
        $('#fileContent').hide()
        $('#tableArea').show()

        $.stopTailFLog()
    })

    $('#clearFileContent').click(function () {
        $('#preContent .preWrap').html('')
    })

    $.createTailTabs = function (content) {
        var tailTabsHtml = ''
        for (var i = 0; i < content.length; ++i) {
            tailTabsHtml += '<button class="tablinks">' + content[i].MachineName + '</button>'
        }
        $('#preContent .tabs').html(tailTabsHtml)
    }


    $.createTailContents = function (content) {
        var datasHtml = ''
        for (var i = 0; i < content.length; ++i) {
            datasHtml += '<div id="machine-' + content[i].MachineName
                + '" class="tabcontent"><pre class="preWrap">' + content[i].TailContent + '</pre></div>'
        }
        $('#preContent .datas').html(datasHtml)

        $.scrollToBottom()
    }

    $.bindTabClicks = function () {
        $('button.tablinks').click(function () {
            $('button.tablinks').removeClass('active')
            $(this).addClass('active')
            $('div.tabcontent').removeClass('active').hide()
            $('#machine-' + $(this).text()).addClass('active').show()
        }).first().click()
    }

    function TailLogFile(loggerName, logPath, lines) {
        $('#refresh').unbind('click')
        $('#locateLogSpan').show()

        $.ajax({
            type: 'POST',
            url: contextPath + "/tailLogFile/" + loggerName + "/" + lines,
            success: function (content, textStatus, request) {
                $.createTailTabs(content)
                $.createTailContents(content)
                $.bindTabClicks()

                $('#tableArea').hide()
                $('#fileContent').show()
                $('#refresh').click(function () {
                    TailLogFile(loggerName, logPath, lines)
                })

                $.bindLocateLogClick(loggerName)
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
                var $row = $(this).parent()
                var loggerName = $row.find('td.LoggerName').text()
                var logPath = $row.find('td.LogPath').text()

                if (key === "TailLogFile") {
                    $('#refresh').unbind('click').show().find('span').text('Refresh')
                    var lines = $.contextMenu.getInputValues(options).Lines
                    TailLogFile(loggerName, logPath, lines)
                } else if (key === 'TailFLog') {
                    $.bindTailFLogEvent(loggerName, logPath)
                }
            },
            items: {
                Lines: {name: "Tail Last Lines:", type: 'text', value: "300"},
                TailLogFile: {name: "Tail Log", icon: "tail"},
                TailFLog: {name: "Tail -F Log", icon: "tail"},
            }
        })
    }
})()