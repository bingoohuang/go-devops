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

    if (typeof (String.prototype.escapeHtml) !== 'function') {
        String.prototype.escapeHtml = escapeHtml
    }

    $('#TestQyMsg').click(function () {
        $.confirm({
            title: '测试黑猫消息推送',
            content: '' +
                '<div>' +
                '<textarea class="input" rows="10" cols="50" >写点啥吧</textarea>' +
                '</div>',
            buttons: {
                formSubmit: {
                    text: 'Submit',
                    btnClass: 'btn-blue',
                    action: function () {
                        var input = this.$content.find('.input').val();
                        $.ajax({
                            type: 'POST',
                            url: contextPath + "/testQywxMsg",
                            data: {msg: input},
                            success: function (content, textStatus, request) {
                            },
                            error: function (jqXHR, textStatus, errorThrown) {
                                alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
                            }
                        })
                    }
                },
                cancel: function () {
                    //close
                }
            }
        });
    })
})()

