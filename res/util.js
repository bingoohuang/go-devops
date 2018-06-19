(function () {
    $.scrollToBottom = function () {
        setTimeout(function () {
            $('html, body').animate({scrollTop: $(document).height()}, 1000)
        }, 300)
    }

    function escapeHtml() {
        return this.replace(/[&<>"'\/]/g, function (s) {
            var entityMap = {
                "&": "&amp;",
                "<": "&lt;",
                ">": "&gt;",
                '"': '&quot;',
                "'": '&#39;',
                "/": '&#x2F;'
            };

            return entityMap[s]
        })
    }

    if (typeof(String.prototype.escapeHtml) !== 'function') {
        String.prototype.escapeHtml = escapeHtml
    }
})()

