(function () {
    $.ajax({
        type: 'POST',
        url: contextPath + "/machines",
        success: function (content, textStatus, request) {
            var machines = $('#machines')
            var machinesHtml = '<table>' +
                '<tr><td>Name</td><td>State</td><td>Hostname</td><td>IP</td>' +
                '</td><td>OS</td><td>Cores</td><td>Total Memory</td><td>Available Memory</td>' +
                '<td>Memory Used</td><td>Total Disk</td><td>Free Disk</td><td>Disk Used</td></tr>'
            if (content && content.length) {
                for (var j = 0; j < content.length; j++) {
                    var machineResult = content[j]
                    var hardwareInfo = machineResult.MachineInfo

                    machinesHtml += '<tr><td>' + machineResult.Name + '</td>'
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
                        + '</tr>'
                }
            }

            machinesHtml += '</table>'
            machines.html(machinesHtml)
        },
        error: function (jqXHR, textStatus, errorThrown) {
            alert(jqXHR.responseText + "\nStatus: " + textStatus + "\nError: " + errorThrown)
        }
    })
})()