(function () {
    var intervalFn = null
    var countSeconds = 0

    $.startWaiting = function () {
        $('#screenShot').show()
        intervalFn = setInterval(function () {
            $('#screenShotTip').html('加载中，请稍后。(' + ++countSeconds + 's)...')
        }, 1000)
    }

    $.endWaiting = function () {
        $('#screenShot').hide()

        if (intervalFn != null) clearInterval(intervalFn)
        intervalFn = null
        countSeconds = 0
    }
})()