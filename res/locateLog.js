(function () {
    $.replaceLocateContents = function (content) {
        var hasContentActivated = false
        for (var i = 0; i < content.length; ++i) {
            var out = content[i].Stdout
            var machineName = content[i].MachineName;
            $('#machine-' + machineName + " .preWrap").html(out)

            if (out != "") {
                $.scrollToBottom()

                if (!hasContentActivated) {
                    hasContentActivated = true
                    setTimeout(function () {
                        $('.tablink-' + machineName).click()
                    }, 500)
                }
            }
        }
    }

    $.bindLocateLogClick = function (loggerName) {
        $('#locateLog').unbind('click').click(function () {
            var fromTimestamp = $('#fromTimestamp').val()
            var toTimestamp = $('#toTimestamp').val()
            $.ajax({
                type: 'POST',
                url: contextPath + "/locateLog/" + loggerName + "/" + fromTimestamp + "/" + toTimestamp,
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