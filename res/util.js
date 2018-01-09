var contextPath = window.location.pathname
if (contextPath.lastIndexOf("/", contextPath.length - 1) !== -1) {
    contextPath = contextPath.substring(0, contextPath.length - 1)
}


(function () {
    $.scrollToBottom = function () {
        setTimeout(function () {
            $('html, body').animate({scrollTop: $(document).height()}, 1000)
        }, 300)
    }
})()

