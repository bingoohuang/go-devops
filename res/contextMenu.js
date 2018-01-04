(function () {
    $('#closeFileContent').click(function () {
        $('#fileContent').hide()
        $('#tableArea').show()

        clearTimeout(ttlTailTimeout)
        ttlTailTimeout = null
    })

    function createTailTabs(content) {
        var tailTabsHtml = ''
        for (var i = 0; i < content.length; ++i) {
            tailTabsHtml += '<button class="tablinks">' + content[i].MachineName + '</button>'
        }
        $('#preContent .tabs').html(tailTabsHtml)
    }

    var scrollToBottom = function () {
        $('html, body').scrollTop($(document).height())
    }

    function appendTailContents(content) {
        for (var i = 0; i < content.length; ++i) {
            if (content[i].TailContent == "") continue

            $('#machine-' + content[i].MachineName + " .preWrap").append(content[i].TailContent)
            scrollToBottom()
        }
    }

    function createTailContents(content) {
        var datasHtml = ''
        for (var i = 0; i < content.length; ++i) {
            datasHtml += '<div id="machine-' + content[i].MachineName
                + '" class="tabcontent"><pre class="preWrap">' + content[i].TailContent + '</pre></div>'
        }
        $('#preContent .datas').html(datasHtml)
    }

    function bindTabClicks() {
        $('button.tablinks').click(function () {
            $('button.tablinks').removeClass('active')
            $(this).addClass('active')
            $('div.tabcontent').removeClass('active').hide()
            $('#machine-' + $(this).text()).addClass('active').show()
        }).first().click()
    }

    var ttlTailTimeout = null
    var NewLogSeq = null

    function TailFLog(loggerName, logPath, tailSeq) {
        $('#locateLogSpan').hide()
        $.ajax({
            type: 'POST',
            url: contextPath + "/tailFLog/" + loggerName + "/" + tailSeq,
            success: function (content, textStatus, request) {
                if (tailSeq == "init") {
                    createTailTabs(content.Results)
                    createTailContents(content.Results)
                    bindTabClicks()

                    $('#tableArea').hide()
                    $('#fileContent').show()
                } else {
                    appendTailContents(content.Results)
                }

                NewLogSeq = content.NewLogSeq
                ttlTailTimeout = setTimeout(function () {
                    TailFLog(loggerName, logPath, NewLogSeq)
                }, 1000)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function replaceLocateContents(content) {
        for (var i = 0; i < content.length; ++i) {
            $('#machine-' + content[i].MachineName + " .preWrap").html(content[i].Stdout)
            scrollToBottom()
        }
    }

    function locateLogClick(loggerName) {
        $('#locateLog').unbind('click').click(function () {
            var fromTimestamp = $('#fromTimestamp').val()
            var toTimestamp = $('#toTimestamp').val()
            $.ajax({
                type: 'POST',
                url: contextPath + "/locateLog/" + loggerName + "/" + fromTimestamp + "/" + toTimestamp,
                success: function (content, textStatus, request) {
                    replaceLocateContents(content)
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    }

    function TailLogFile(loggerName, logPath, lines) {
        $('#refresh').unbind('click')
        $('#locateLogSpan').show()

        $.ajax({
            type: 'POST',
            url: contextPath + "/tailLogFile/" + loggerName + "/" + lines,
            success: function (content, textStatus, request) {
                createTailTabs(content)
                createTailContents(content)
                bindTabClicks()

                $('#tableArea').hide()
                $('#fileContent').show()
                $('#refresh').click(function () {
                    TailLogFile(loggerName, logPath, lines)
                })

                locateLogClick(loggerName);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

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
        var lines = 10
        $.contextMenu({
            selector: '.LogPath',
            callback: function (key, options, rootMenu, originalEvent) {
                if (key === "TailLogFile") {
                    var $row = $(this).parent();
                    var loggerName = $row.find('td.LoggerName').text();
                    var logPath = $row.find('td.LogPath').text();
                    $('#refresh').unbind('click').show().find('span').text('Refresh')
                    TailLogFile(loggerName, logPath, lines)
                } else if (key === 'TailFLog') {
                    var $row = $(this).parent();
                    var loggerName = $row.find('td.LoggerName').text();
                    var logPath = $row.find('td.LogPath').text();
                    $('#refresh').unbind('click').click(function () {
                        var span = $(this).find('span');
                        if (span.text() === "Stop") {
                            clearTimeout(ttlTailTimeout)
                            ttlTailTimeout = null
                            span.text('Resume')
                        } else if (span.text() === "Resume") {
                            TailFLog(loggerName, logPath, NewLogSeq)
                            span.text('Stop')
                        }
                    }).find('span').text('Stop')
                    TailFLog(loggerName, logPath, "init")
                }
            },
            items: {
                // <input type="text">
                Lines: {
                    name: "How Many Lines to Tail",
                    type: 'text',
                    value: "10",
                    events: {
                        keyup: function (e) {
                            lines = $(this).val()
                        }
                    }
                },
                TailLogFile: {name: "Tail Log", icon: "tail"},
                TailFLog: {name: "Tail -F Log", icon: "tail"},
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
                TruncateLogFile: {name: "Truncate Log File", icon: "truncate"}
            }
        })
    }

    $.createProcessInfoContextMenu = function () {
        $.contextMenu({
            selector: '.ProcessInfo',
            callback: function (key, options) {
                if (key === "RestartProcess") {
                    var $cell = $(this);
                    var $row = $cell.parent();
                    var logMachine = $row.find('td.LogMachine').text();
                    var loggerName = $row.find('td.LoggerName').text();
                    RestartProcess($cell, logMachine, loggerName)
                }
            },
            items: {
                RestartProcess: {name: "Restart process", icon: "restart"}
            }
        })
    }

})()