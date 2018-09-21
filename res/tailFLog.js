(function () {
    var ttlTailTimeout = null
    var NewLogSeq = null

    function TailFLog(loggerName, traceMobile, tailSeq) {
        $('#locateLogSpan').hide()
        $.ajax({
            type: 'POST',
            url: contextPath + "/tailFLog/" + loggerName + "/" + traceMobile + "/" + tailSeq,
            success: function (content, textStatus, request) {
                if (tailSeq == "init") {
                    $.createTailTabs(content.Results)
                    $.createTailContents(content.Results)
                    $.bindTabClicks()

                    $('#tableArea').hide()
                    $('#fileContent').show()
                } else {
                    appendTailContents(content.Results)
                }

                if (ttlTailTimeout != null) {
                    NewLogSeq = content.NewLogSeq
                    ttlTailTimeout = setTimeout(function () {
                        TailFLog(loggerName, traceMobile, NewLogSeq)
                    }, 1000)
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    function appendTailContents(content) {
        var maxSize = 1 * 1024 * 1024; // 1M
        for (var i = 0; i < content.length; ++i) {
            if (content[i].TailContent == "") continue

            var machinePreWrap = $('#machine-' + content[i].MachineName + " .preWrap")
            machinePreWrap.append(content[i].TailContent.escapeHtml())
            var textLength = machinePreWrap.text().length

            if (ttlTailTimeout != null && textLength > maxSize) {
                machinePreWrap.text(machinePreWrap.text().substring(textLength - maxSize))
            }

            if (ttlTailTimeout != null) {
                window.scrollTo(0, document.body.scrollHeight)
            }
        }
    }

    $.stopTailFLog = function () {
        var temp = ttlTailTimeout
        ttlTailTimeout = null
        clearTimeout(temp)
    }

    $.bindTailFLogEvent = function (loggerName, traceMobile) {
        $('#refresh').unbind('click').click(function () {
            var span = $(this).find('span');
            if (span.text() === "Stop") {
                $.stopTailFLog()

                span.text('Resume')
            } else if (span.text() === "Resume") {
                ttlTailTimeout = {}
                TailFLog(loggerName, traceMobile, NewLogSeq)
                span.text('Stop')
            }
        }).find('span').text('Stop')

        ttlTailTimeout = {}
        TailFLog(loggerName, traceMobile, "init")
    }
})()