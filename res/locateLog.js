(function () {
    $.replaceLocateContents = function (content) {
        var tablLink = null
        for (var i = 0; i < content.length; ++i) {
            var out = content[i].Stdout.escapeHtml()
            var machineName = content[i].MachineName;
            $('#machine-' + machineName + " .preWrap").html(out)

            if (tablLink == null && out !== "") {
                tablLink = $('.tablink-' + machineName)
            }
        }

        if (tablLink != null) {
            setTimeout(function () {tablLink.click()}, 500)
        }
    }

    $.bindLocateLogClick = function (loggerName) {
        $('#locateLog').unbind('click').click(function () {
            var logKey = $('#logKey').val()
            var preLines = $('#preLines').val()
            var lines = $('#lines').val()
            $.ajax({
                type: 'POST',
                url: contextPath + "/locateLog/" + loggerName + "/" + logKey + "/" + preLines + "/" + lines,
                success: function (content, textStatus, request) {
                    $.replaceLocateContents(content)
                },
                error: function (jqXHR, textStatus, errorThrown) {
                    alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                }
            })
        })
    }
})()