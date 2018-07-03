(function () {
    var replaceLocateContents = function (content) {
        var tablLink = null
        for (var i = 0; i < content.length; ++i) {
            if (!content[i].Ready) {
                continue
            }

            var out = content[i].Stdout.escapeHtml()
            var machineName = content[i].MachineName
            $('#machine-' + machineName + " .preWrap").html(out)

            if (tablLink == null && out !== "") {
                tablLink = $('.tablink-' + machineName)
            }
        }

        if (tablLink != null) {
            setTimeout(function () {
                tablLink.click()
            }, 500)
        }
    }

    var locateLogResult = function (ShellResultCommandResult) {
        var data = {}
        for (var i = 0; i < ShellResultCommandResult.length; ++i) {
            var result = ShellResultCommandResult[i]
            if (!result.Ready) {
                data[result.MachineName] = result.ShellKey
            }
        }

        replaceLocateContents(ShellResultCommandResult)

        if (Object.keys(data).length === 0) {
            $.endWaiting()
            return
        }

        $.ajax({
            type: 'GET',
            url: contextPath + "/locateLogResult/",
            data: data,
            success: function (content, textStatus, request) {
                if (content.Err != "") {
                    alert(content.Err)
                    return
                }

                setTimeout(function () {
                    locateLogResult(content.Results)
                }, 10000)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
            }
        })
    }

    $.bindLocateLogClick = function (loggerName) {
        $('#locateLog').unbind('click').click(function () {
            var logKey = $('#logKey').val()
            var preLines = $('#preLines').val()
            var lines = $('#lines').val()

            $.startWaiting()
            $.ajax({
                type: 'POST',
                url: contextPath + "/locateLog/" + loggerName + "/" + logKey + "/" + preLines + "/" + lines,
                success: function (content, textStatus, request) {
                    setTimeout(function () {
                        locateLogResult(content)
                    }, 10000)
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    }
})()