(function () {
    var lastHighlightMachineCells = []

    function highlightLogMachineCells(machineName) {
        $.each(lastHighlightMachineCells, function (index, value) {
            value.removeClass('highlight')
        })

        lastHighlightMachineCells = []
        $("#logs td.LogMachine").each(function () {
            if ($(this).text() === machineName) {
                lastHighlightMachineCells.push($(this))
            }
            return true
        })

        $.each(lastHighlightMachineCells, function (index, value) {
            value.addClass('highlight')
        })
    }

    var lastHighlightMachineRow = null

    function highlightMachineRow(machineName) {
        if (lastHighlightMachineRow != null) {
            lastHighlightMachineRow.removeClass('highlight')
        }

        $('#machines td.MachineName').each(function () {
            if (machineName !== $(this).text()) return true

            lastHighlightMachineRow = $(this).parent()
            lastHighlightMachineRow.addClass('highlight')
            return false
        })
    }

    $("#logs").on("click", "td.LogMachine", function () {
        var machineName = $(this).text();
        highlightLogMachineCells(machineName)
        highlightMachineRow(machineName)
    })

})()