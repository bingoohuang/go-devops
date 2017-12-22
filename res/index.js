(function () {
    $.ajax({
        type: 'POST',
        url: contextPath + "/machines",
        success: function (content, textStatus, request) {
            var machinesHtml = createMachinesTable(content)
            $('#machines').html(machinesHtml)
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })

    $.ajax({
        type: 'POST',
        url: contextPath + "/logs",
        success: function (content, textStatus, request) {
            var logsHtml = createLogsTable(content)
            $('#logs').html(logsHtml)
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })


    var lastHighlightMachineCells = []

    function hightlightLogMachineCells(machineName) {
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
        hightlightLogMachineCells(machineName)
        highlightMachineRow(machineName)
    })

    function createMachinesTable(content) {
        var machinesHtml = '<table>' +
            '<tr>' +
            '<td>Name</td>' +
            '<td>State</td>' +
            '<td>Hostname</td>' +
            '<td>IP</td>' +
            '<td>OS</td>' +
            '<td>Cores</td>' +
            '<td>Total Memory</td>' +
            '<td>Available Memory</td>' +
            '<td>Memory Used</td>' +
            '<td>Total Disk</td>' +
            '<td>Free Disk</td>' +
            '<td>Disk Used</td>' +
            '<td>Cost</td>' +
            '</tr>'
        if (content && content.length) {
            for (var j = 0; j < content.length; j++) {
                var machineResult = content[j]
                var hardwareInfo = machineResult.MachineInfo

                machinesHtml += '<tr><td class="MachineName">' + machineResult.Name + '</td>'
                    + '<td>' + (machineResult.Error || 'OK') + '</td>'
                    + '<td>' + hardwareInfo.Hostname + '</td>'
                    + '<td>' + (hardwareInfo.Ips && hardwareInfo.Ips.length > 0 ? hardwareInfo.Ips.join(', ') : '') + '</td>'
                    + '<td>' + hardwareInfo.OS + '</td>'
                    + '<td>' + hardwareInfo.Cores + '</td>'
                    + '<td>' + hardwareInfo.HumanizedTotalMemory + '</td>'
                    + '<td>' + hardwareInfo.HumanizedAvailableMemory + '</td>'
                    + '<td>' + new Number(hardwareInfo.MemoryUsedPercent).toFixed(2) + '%</td>'
                    + '<td>' + hardwareInfo.HumanizedTotalDisk + '</td>'
                    + '<td>' + hardwareInfo.HumanizedFreeDisk + '</td>'
                    + '<td>' + new Number(hardwareInfo.DiskUsedPercent).toFixed(2) + '%</td>'
                    + '<td>' + machineResult.CostTime + '</td>'
                    + '</tr>'
            }
        }

        machinesHtml += '</table>'
        return machinesHtml;
    }

    function createLogsTable(content) {
        var logsHtml = '<table>' +
            '<tr>' +
            '<td>Logger Name</td>' +
            '<td>Log Path</td>' +
            '<td>State</td>' +
            '<td>Machine</td>' +
            '<td>Size</td>' +
            '<td>Last Modified</td>' +
            '<td>ProcessInfo</td>' +
            '<td>Cost</td>' +
            '</tr>'

        if (content && content.length) {
            for (var j = 0; j < content.length; j++) {
                var log = content[j]
                var logMachines = log.Logs

                for (var i = 0; i < logMachines.length; ++i) {
                    var logMachine = log.Logs[i]

                    logsHtml += '<tr>'
                    if (i == 0) {
                        logsHtml += '<td rowspan="' + logMachines.length + '">' + log.Logger + '</td>'
                            + '<td class="LogPath" rowspan="' + logMachines.length + '">' + log.LogPath + '</td>'
                    } else {
                        logsHtml += '<td class="hidden">' + log.Logger + '</td>'
                            + '<td class="LogPath hidden">' + log.LogPath + '</td>'
                    }

                    logsHtml += '<td>' + (logMachine.Error || 'OK') + '</td>'
                        + '<td class="LogMachine">' + logMachine.MachineName + '</td>'
                        + '<td class="LogFileSize">' + logMachine.FileSize + '</td>'
                        + '<td>' + logMachine.LastModified + '</td>'
                        + '<td class="ProcessInfo">' + logMachine.ProcessInfo + '</td>'
                        + '<td>' + logMachine.CostTime + '</td>'
                        + '</tr>'
                }
            }
        }

        logsHtml += '</table>'
        return logsHtml;
    }

    $.contextMenu({
        selector: '.LogFileSize',
        callback: function (key, options) {
            if (key === "TruncateLogFile") {
                alert($(this).parent().find('td.LogPath').text())
            }
        },
        items: {
            "TruncateLogFile": {name: "TruncateLogFile", icon: "cut"}
        }
    })
})()