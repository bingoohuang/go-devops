(function () {
    $.scrollToBottom = function () {
        setTimeout(function () {
            $('html, body').animate({scrollTop: $(document).height()}, 1000)
        }, 300)
    }
})()

