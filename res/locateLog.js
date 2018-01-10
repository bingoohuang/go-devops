(function () {
    $.replaceLocateContents = function(content) {
        for (var i = 0; i < content.length; ++i) {
            $('#machine-' + content[i].MachineName + " .preWrap").html(content[i].Stdout)
            $.scrollToBottom()
        }
    }

    $.bindLocateLogClick = function(loggerName) {
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
})()