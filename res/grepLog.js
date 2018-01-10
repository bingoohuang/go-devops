(function () {
    $.bindGrepLogClick = function(loggerName) {
        $('#grepLog').unbind('click').click(function () {
            var grepText = $('#grepText').val()
            $.ajax({
                type: 'POST',
                url: contextPath + "/grepLog/" + loggerName + "/" + grepText,
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
